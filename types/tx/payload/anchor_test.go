package payload_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// number 1, balance 2, 32-byte hash of 0x11, empty URI, deposit 1, timestamps 10/20/30/40.
const anchoredAccountHex = "01000000" +
	"0200000000000000" +
	"01" +
	"20" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"00" +
	"00" +
	"0100000000000000" +
	"0a000000" +
	"14000000" +
	"1e000000" +
	"28000000"

const plainAccountHex = "010000000200000000000000"

// Set: BLS address, hash of 0x11, URI "abc", type 0x04, deposit 300 as a varint.
const frozenSetHex = "020102030405060708090a0b0c0d0e0f1011121314" +
	"00" +
	"20" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"03" +
	"616263" +
	"04" +
	"ac02"

const frozenDeleteHex = "020102030405060708090a0b0c0d0e0f1011121314" + "01"

func blsAddr() crypto.Address {
	return crypto.NewAddress(crypto.AddressTypeBLSAccount, []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
	})
}

func encodeAnchor(tb testing.TB, pld *payload.AnchorPayload) []byte {
	tb.Helper()

	buf := bytes.NewBuffer(nil)
	require.NoError(tb, pld.Encode(buf))
	require.Len(tb, buf.Bytes(), pld.SerializeSize())

	return buf.Bytes()
}

func decodeAnchor(t *testing.T, data []byte) *payload.AnchorPayload {
	t.Helper()

	pld := &payload.AnchorPayload{}
	require.NoError(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(data)))

	return pld
}

func setPayload(
	addr crypto.Address, hashLen int, uri string, anchorType uint8, deposit amount.Amount,
) *payload.AnchorPayload {
	root := bytes.Repeat([]byte{0x11}, hashLen)

	return &payload.AnchorPayload{
		From:        addr,
		Action:      payload.AnchorActionSet,
		RootHash:    root,
		ManifestURI: uri,
		AnchorType:  anchorType,
		Deposit:     deposit,
	}
}

func TestAnchorTypeString(t *testing.T) {
	require.Equal(t, "anchor", payload.TypeAnchor.String())
	require.Equal(t, "0", payload.Type(0).String())
	require.Equal(t, "8", payload.Type(8).String())
	require.Equal(t, "transfer", payload.TypeTransfer.String())
	require.Equal(t, "bond", payload.TypeBond.String())
	require.Equal(t, "sortition", payload.TypeSortition.String())
	require.Equal(t, "unbond", payload.TypeUnbond.String())
	require.Equal(t, "withdraw", payload.TypeWithdraw.String())
	require.Equal(t, "batch-transfer", payload.TypeBatchTransfer.String())

	set := setPayload(blsAddr(), 32, "", 0, 0)
	del := &payload.AnchorPayload{From: blsAddr(), Action: payload.AnchorActionDelete}
	require.Equal(t, payload.TypeAnchor, set.Type())
	require.Equal(t, payload.TypeAnchor, del.Type())
	require.Equal(t, set.From, set.Signer())
	require.Equal(t, del.From, del.Signer())
	require.Contains(t, set.LogString(), "Anchor")
	require.Contains(t, del.LogString(), "Anchor")
	require.NotContains(t, del.LogString(), "11")
}

func TestAnchorRoundTrip(t *testing.T) {
	addr := blsAddr()
	cases := []struct {
		hashLen    int
		uri        string
		anchorType uint8
		deposit    amount.Amount
	}{
		{32, "", 0x00, 0},
		{32, "", 0x04, 1},
		{64, strings.Repeat("a", 128), 0xFF, amount.Amount(amount.MaxNanoPAC)},
		{33, "abc", 0x04, 127},
		{63, "abc", 0x00, 128},
		{48, "é", 0xFF, 300},
	}

	for _, tc := range cases {
		original := setPayload(addr, tc.hashLen, tc.uri, tc.anchorType, tc.deposit)
		require.NoError(t, original.BasicCheck())
		require.Equal(t, tc.deposit, original.Value())

		data := encodeAnchor(t, original)
		decoded := decodeAnchor(t, data)
		require.Equal(t, original.From, decoded.From)
		require.Equal(t, original.Action, decoded.Action)
		require.Equal(t, original.RootHash, decoded.RootHash)
		require.Equal(t, original.ManifestURI, decoded.ManifestURI)
		require.Equal(t, original.AnchorType, decoded.AnchorType)
		require.Equal(t, original.Deposit, decoded.Deposit)
		require.Equal(t, data, encodeAnchor(t, decoded))
	}

	empty := setPayload(addr, 32, "", 0, 0)
	require.Len(t, encodeAnchor(t, empty), 58)

	full := setPayload(addr, 64, strings.Repeat("b", 128), 0xFF, amount.Amount(amount.MaxNanoPAC))
	require.Len(t, encodeAnchor(t, full), 225)

	small := setPayload(addr, 32, "", 0, 127)
	wide := setPayload(addr, 32, "", 0, 128)
	require.Len(t, encodeAnchor(t, wide), len(encodeAnchor(t, small))+1)
	require.Equal(t, []byte{0x80, 0x01}, encodeAnchor(t, wide)[len(encodeAnchor(t, wide))-2:])
}

func TestAnchorFrozenWire(t *testing.T) {
	raw, err := hex.DecodeString(frozenSetHex)
	require.NoError(t, err)

	decoded := decodeAnchor(t, raw)
	require.Equal(t, blsAddr(), decoded.From)
	require.Equal(t, payload.AnchorActionSet, decoded.Action)
	require.Equal(t, bytes.Repeat([]byte{0x11}, 32), decoded.RootHash)
	require.Equal(t, "abc", decoded.ManifestURI)
	require.Equal(t, uint8(0x04), decoded.AnchorType)
	require.Equal(t, amount.Amount(300), decoded.Deposit)
	require.Equal(t, frozenSetHex, hex.EncodeToString(encodeAnchor(t, decoded)))
	require.Equal(t, []byte{0xac, 0x02}, raw[len(raw)-2:])

	delRaw, err := hex.DecodeString(frozenDeleteHex)
	require.NoError(t, err)
	deleted := decodeAnchor(t, delRaw)
	require.Equal(t, payload.AnchorActionDelete, deleted.Action)
	require.Equal(t, amount.Amount(0), deleted.Value())
	require.Len(t, delRaw, 22)
	require.Equal(t, frozenDeleteHex, hex.EncodeToString(encodeAnchor(t, deleted)))

	setRaw := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 0x04, 300))
	require.Equal(t, raw[:21], setRaw[:21])
	require.Equal(t, raw[:21], delRaw[:21])
	require.Equal(t, payload.AnchorActionSet, setRaw[21])
	require.Equal(t, payload.AnchorActionDelete, delRaw[21])
}

func TestAnchorDeleteDropsSetFields(t *testing.T) {
	pld := &payload.AnchorPayload{
		From:        blsAddr(),
		Action:      payload.AnchorActionDelete,
		RootHash:    bytes.Repeat([]byte{0x11}, 32),
		ManifestURI: "still-here",
		AnchorType:  0x04,
		Deposit:     300,
	}
	require.Equal(t, amount.Amount(0), pld.Value())
	require.NoError(t, pld.BasicCheck())
	require.Len(t, encodeAnchor(t, pld), 22)

	decoded := decodeAnchor(t, encodeAnchor(t, pld))
	require.Empty(t, decoded.RootHash)
	require.Empty(t, decoded.ManifestURI)
	require.Equal(t, amount.Amount(0), decoded.Deposit)
	require.Equal(t, pld.From, decoded.Signer())

	junk := &payload.AnchorPayload{
		From:        blsAddr(),
		Action:      payload.AnchorActionDelete,
		RootHash:    bytes.Repeat([]byte{0x11}, 31),
		ManifestURI: string([]byte{0xC3}),
		Deposit:     -1,
	}
	require.NoError(t, junk.BasicCheck())
	require.Equal(t, amount.Amount(0), junk.Value())
	require.Len(t, encodeAnchor(t, junk), 22)
}

func TestAnchorSignerNotTreasury(t *testing.T) {
	pld := &payload.AnchorPayload{From: crypto.TreasuryAddress, Action: payload.AnchorActionSet}
	data := encodeAnchor(t, &payload.AnchorPayload{
		From:   blsAddr(),
		Action: payload.AnchorActionDelete,
	})
	require.NoError(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(data)))
	require.Equal(t, blsAddr(), pld.Signer())
	require.False(t, pld.Signer().IsTreasuryAddress())
}

func TestAnchorBasicCheckBounds(t *testing.T) {
	addr := blsAddr()
	require.NoError(t, setPayload(addr, 32, "", 0, 0).BasicCheck())
	require.NoError(t, setPayload(addr, 32, "", 0, 1).BasicCheck())
	require.NoError(t, setPayload(addr, 32, "", 0, amount.Amount(amount.MaxNanoPAC)).BasicCheck())
	require.NoError(t, setPayload(addr, 32, strings.Repeat("a", 128), 0, 1).BasicCheck())
	require.NoError(t, setPayload(addr, 33, "", 0, 1).BasicCheck())
	require.NoError(t, setPayload(addr, 63, "", 0, 1).BasicCheck())
	require.NoError(t, setPayload(addr, 64, "", 0, 1).BasicCheck())
	require.NoError(t, (&payload.AnchorPayload{From: addr, Action: payload.AnchorActionDelete}).BasicCheck())

	ts := testsuite.NewTestSuite(t)
	blsPub, _ := ts.RandBLSKeyPair()
	edPub, _ := ts.RandEd25519KeyPair()
	secpPub, _ := ts.RandSecp256k1KeyPair()
	require.NoError(t, setPayload(blsPub.AccountAddress(), 32, "", 0, 1).BasicCheck())
	require.NoError(t, setPayload(edPub.AccountAddress(), 32, "", 0, 1).BasicCheck())
	require.NoError(t, setPayload(secpPub.AccountAddress(), 32, "", 0, 1).BasicCheck())

	require.Error(t, setPayload(addr, 31, "", 0, 1).BasicCheck())
	require.Error(t, setPayload(addr, 65, "", 0, 1).BasicCheck())
	require.Error(t, setPayload(addr, 0, "", 0, 1).BasicCheck())
	require.Error(t, setPayload(addr, 32, strings.Repeat("a", 129), 0, 1).BasicCheck())
	require.Error(t, setPayload(addr, 32, "", 0, amount.Amount(amount.MaxNanoPAC)+1).BasicCheck())
	require.Error(t, setPayload(addr, 32, "", 0, -1).BasicCheck())
	require.Error(t, setPayload(crypto.TreasuryAddress, 32, "", 0, 1).BasicCheck())
	require.Error(t, setPayload(ts.RandValAddress(), 32, "", 0, 1).BasicCheck())
	require.Error(t, (&payload.AnchorPayload{From: addr, Action: 2}).BasicCheck())
	require.Error(t, (&payload.AnchorPayload{From: addr, Action: 0xFF}).BasicCheck())
	require.Error(t, setPayload(addr, 32, string([]byte{0xC3}), 0, 1).BasicCheck())
}

func TestAnchorDecodeRejects(t *testing.T) {
	valid := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 0x04, 300))
	for length := 0; length < len(valid); length++ {
		pld := &payload.AnchorPayload{}
		err := pld.Decode(payload.DecodeContext{}, bytes.NewReader(valid[:length]))
		require.Error(t, err)
		require.Error(t, pld.BasicCheck())
	}

	filled := setPayload(blsAddr(), 32, "abc", 0x04, 1)
	require.Error(t, filled.Decode(payload.DecodeContext{}, bytes.NewReader(nil)))
	require.Error(t, filled.BasicCheck())

	addr := blsAddr().Bytes()
	actionSet := append(append([]byte{}, addr...), payload.AnchorActionSet)
	for _, action := range []byte{2, 0xFF} {
		body := append(append([]byte{}, addr...), action)
		pld := &payload.AnchorPayload{}
		require.Error(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(body)))
		require.Error(t, pld.BasicCheck())
	}

	for _, hashLen := range []byte{0, 31, 65, 255} {
		body := append(append([]byte{}, actionSet...), hashLen)
		body = append(body, bytes.Repeat([]byte{0x11}, int(hashLen))...)
		pld := &payload.AnchorPayload{}
		require.Error(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(body)))
	}

	hash32 := append(append([]byte{}, actionSet...), 32)
	hash32 = append(hash32, bytes.Repeat([]byte{0x11}, 32)...)
	pld := &payload.AnchorPayload{}
	require.Error(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(append(hash32, 255))))
	require.Error(t, pld.Decode(payload.DecodeContext{}, bytes.NewReader(append(hash32, 129, 0))))

	badURI := []struct {
		name string
		uri  []byte
	}{
		{"truncated", []byte{0xC3}},
		{"overlong", []byte{0xC0, 0x80}},
		{"surrogate", []byte{0xED, 0xA0, 0x80}},
	}
	for _, tc := range badURI {
		body := append(append([]byte{}, hash32...), byte(len(tc.uri)))
		body = append(body, tc.uri...)
		got := &payload.AnchorPayload{}
		require.Error(t, got.Decode(payload.DecodeContext{}, bytes.NewReader(body)), tc.name)
		require.Error(t, got.BasicCheck())
	}

	shifted := append([]byte{0x00}, valid...)
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(shifted)))
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader([]byte{0x05})))

	nonMinimal := append([]byte{}, valid[:len(valid)-2]...)
	nonMinimal = append(nonMinimal, 0x81, 0x00)
	err := (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(nonMinimal))
	require.ErrorIs(t, err, encoding.ErrNonCanonical)

	overflow := append([]byte{}, hash32...)
	overflow = append(overflow, 0x00, 0x00)
	overflow = append(overflow, bytes.Repeat([]byte{0x80}, 10)...)
	err = (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(overflow))
	require.ErrorIs(t, err, encoding.ErrOverflow)

	overMax := append([]byte{}, valid[:len(valid)-2]...)
	overMax = append(overMax, varInt(uint64(amount.Amount(amount.MaxNanoPAC))+1)...)
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(overMax)))
}

func TestAnchorParserAttacks(t *testing.T) {
	setData := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 4, 300))
	delData := encodeAnchor(t, &payload.AnchorPayload{From: blsAddr(), Action: payload.AnchorActionDelete})
	stuck := append(append([]byte{}, delData...), setData...)
	reader := bytes.NewReader(stuck)
	decoded := &payload.AnchorPayload{}
	require.NoError(t, decoded.Decode(payload.DecodeContext{}, reader))
	require.Empty(t, decoded.RootHash)
	require.Equal(t, len(stuck)-len(delData), reader.Len())

	withDelegation := decodeAnchor(t, setData)
	again := &payload.AnchorPayload{}
	require.NoError(t, again.Decode(payload.DecodeContext{WithDelegation: true}, bytes.NewReader(setData)))
	require.Equal(t, withDelegation.RootHash, again.RootHash)
	require.Equal(t, withDelegation.ManifestURI, again.ManifestURI)
	require.Equal(t, withDelegation.Deposit, again.Deposit)

	hostile := []string{"\n", "\x00", "javascript:alert(1)", "é", "e\u0301"}
	for _, uri := range hostile {
		original := setPayload(blsAddr(), 32, uri, 0, 1)
		require.NoError(t, original.BasicCheck())
		got := decodeAnchor(t, encodeAnchor(t, original))
		require.Equal(t, uri, got.ManifestURI)
		require.LessOrEqual(t, len(got.ManifestURI), 128)
		require.True(t, utf8.ValidString(got.ManifestURI))
	}
	require.NotEqual(t, "é", "e\u0301")
	require.NoError(t, setPayload(blsAddr(), 32, strings.Repeat("é", 64), 0, 1).BasicCheck())
	require.Error(t, setPayload(blsAddr(), 32, strings.Repeat("a", 129), 0, 1).BasicCheck())

	original := setPayload(blsAddr(), 32, "abc", 4, 300)
	data := encodeAnchor(t, original)
	input := append([]byte{}, data...)
	first := decodeAnchor(t, input)
	second := decodeAnchor(t, input)
	input[len(input)-1] ^= 0xff
	require.Equal(t, bytes.Repeat([]byte{0x11}, 32), first.RootHash)
	first.RootHash[0] ^= 0xff
	require.NotEqual(t, first.RootHash[0], second.RootHash[0])

	encoded := encodeAnchor(t, original)
	extended := append([]byte{}, encoded...)
	extended = append(extended, 0xff)
	require.Equal(t, encoded, encodeAnchor(t, original))
	require.NotEqual(t, extended, encodeAnchor(t, original))

	longHash := setPayload(blsAddr(), 32, "", 0, 1)
	longHash.RootHash = bytes.Repeat([]byte{0x11}, 300)
	buf := bytes.NewBuffer(nil)
	require.Error(t, longHash.Encode(buf))
	require.Len(t, buf.Bytes(), 22)

	longURI := setPayload(blsAddr(), 32, strings.Repeat("a", 300), 0, 1)
	buf.Reset()
	require.Error(t, longURI.Encode(buf))
	require.Len(t, buf.Bytes(), 22)
}

func TestAnchorDecodeEveryByte(t *testing.T) {
	valid := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 4, 300))
	require.NoError(t, decodeAnchor(t, valid).BasicCheck())

	onlyAddr := valid[:21]
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(onlyAddr)))
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(valid[:len(valid)-1])))

	del := encodeAnchor(t, &payload.AnchorPayload{From: blsAddr(), Action: payload.AnchorActionDelete})
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(del[:len(del)-1])))
}

func TestAnchorFieldChangesBytes(t *testing.T) {
	base := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 0, 300))
	otherType := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 1, 300))
	otherHash := setPayload(blsAddr(), 32, "abc", 0, 300)
	otherHash.RootHash[len(otherHash.RootHash)-1] ^= 0xff
	otherDeposit := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 0, 301))

	require.NotEqual(t, base, otherType)
	require.NotEqual(t, base, encodeAnchor(t, otherHash))
	require.NotEqual(t, base, otherDeposit)
}

func TestAnchorAssembleWithAccount(t *testing.T) {
	hash32 := bytes.Repeat([]byte{0x11}, 32)
	uri := "manifest"
	anchorType := uint8(0x04)
	deposit := amount.Amount(1)
	pld := setPayload(blsAddr(), 32, uri, anchorType, deposit)
	pld.RootHash = hash32
	require.NoError(t, pld.BasicCheck())

	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      pld.RootHash,
		ManifestURI:   pld.ManifestURI,
		AnchorType:    pld.AnchorType,
		LockedDeposit: pld.Deposit,
	}))
	require.Equal(t, hash32, acc.RootHash())
	require.Equal(t, uri, acc.ManifestURI())
	require.Equal(t, anchorType, acc.AnchorType())
	require.Equal(t, deposit, acc.LockedDeposit())

	zero := setPayload(blsAddr(), 32, uri, anchorType, 0)
	require.NoError(t, zero.BasicCheck())
	require.Error(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      zero.RootHash,
		ManifestURI:   zero.ManifestURI,
		AnchorType:    zero.AnchorType,
		LockedDeposit: zero.Deposit,
	}))

	short := setPayload(blsAddr(), 31, uri, anchorType, 1)
	require.Error(t, short.BasicCheck())
	require.Error(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      short.RootHash,
		ManifestURI:   short.ManifestURI,
		LockedDeposit: 1,
	}))

	longURI := setPayload(blsAddr(), 32, strings.Repeat("a", 129), anchorType, 1)
	require.Error(t, longURI.BasicCheck())
	require.Error(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      hash32,
		ManifestURI:   longURI.ManifestURI,
		LockedDeposit: 1,
	}))

	payload128 := encodeAnchor(t, setPayload(blsAddr(), 32, "", 0, 128))
	require.Equal(t, []byte{0x80, 0x01}, payload128[len(payload128)-2:])

	locked := account.NewAccount(7)
	require.NoError(t, locked.SetAnchor(account.AnchorData{
		RootHash:      hash32,
		LockedDeposit: 128,
	}))
	accountBytes, err := locked.Bytes()
	require.NoError(t, err)
	require.True(t, bytes.Contains(accountBytes, []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	require.False(t, bytes.Contains(payload128, []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))

	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	anchored, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(t, err)
	require.False(t, anchorAccepted(plain))
	require.False(t, anchorAccepted(anchored))

	frozen, err := hex.DecodeString(frozenSetHex)
	require.NoError(t, err)
	_, err = account.FromBytes(frozen)
	require.Error(t, err)

	numbered := account.NewAccount(2)
	numberedBytes, err := numbered.Bytes()
	require.NoError(t, err)
	require.Equal(t, byte(0x02), numberedBytes[0])
	require.Len(t, numberedBytes, 12)
	require.False(t, anchorAccepted(numberedBytes))

	stored, err := account.FromBytes(plain)
	require.NoError(t, err)
	before, err := stored.Bytes()
	require.NoError(t, err)
	_ = decodeAnchor(t, frozen)
	after, err := stored.Bytes()
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.Len(t, after, 12)
	require.False(t, stored.HasAnchor())

	expected, err := hash.FromString("c3b75f08e64a66cb980fdc03c3a0b78635a7b1db049096e8bbbd9a2873f3071a")
	require.NoError(t, err)
	require.Equal(t, expected, stored.Hash())

	wide := setPayload(blsAddr(), 64, "manifest", anchorType, deposit)
	require.NoError(t, wide.BasicCheck())
	require.NoError(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      wide.RootHash,
		ManifestURI:   wide.ManifestURI,
		AnchorType:    wide.AnchorType,
		LockedDeposit: wide.Deposit,
	}))

	tooWide := setPayload(blsAddr(), 65, "manifest", anchorType, deposit)
	require.Error(t, tooWide.BasicCheck())
	require.Error(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      tooWide.RootHash,
		ManifestURI:   tooWide.ManifestURI,
		LockedDeposit: deposit,
	}))

	badText := string([]byte{0xC3})
	require.Error(t, setPayload(blsAddr(), 32, badText, anchorType, deposit).BasicCheck())
	require.Error(t, account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      hash32,
		ManifestURI:   badText,
		LockedDeposit: deposit,
	}))
}

func TestAnchorDecodeStoresTightHash(t *testing.T) {
	body := encodeAnchor(t, setPayload(blsAddr(), 64, "abc", 4, 300))
	backing := make([]byte, len(body)+4096)
	copy(backing, body)

	decoded := &payload.AnchorPayload{}
	require.NoError(t, decoded.Decode(payload.DecodeContext{}, bytes.NewReader(backing[:len(body)])))
	require.Len(t, decoded.RootHash, 64)
	require.Equal(t, 64, cap(decoded.RootHash))
}

func TestAnchorEncodeMatchesDecode(t *testing.T) {
	addr := blsAddr()
	rejected := []*payload.AnchorPayload{
		setPayload(addr, 0, "", 0, 1),
		setPayload(addr, 31, "", 0, 1),
		setPayload(addr, 65, "", 0, 1),
		setPayload(addr, 32, strings.Repeat("a", 129), 0, 1),
		setPayload(addr, 32, string([]byte{0xC3}), 0, 1),
		setPayload(addr, 32, "", 0, -1),
		setPayload(addr, 32, "", 0, amount.Amount(amount.MaxNanoPAC)+1),
		{From: addr, Action: 2, RootHash: bytes.Repeat([]byte{0x11}, 32), Deposit: 1},
	}
	for _, pld := range rejected {
		buf := bytes.NewBuffer(nil)
		require.Error(t, pld.Encode(buf))
		require.Len(t, buf.Bytes(), 22)
		got := &payload.AnchorPayload{}
		require.Error(t, got.Decode(payload.DecodeContext{}, bytes.NewReader(buf.Bytes())))
		require.Error(t, got.BasicCheck())
	}

	zero := setPayload(addr, 32, "", 0, 1)
	zero.RootHash = make([]byte, 32)
	require.NoError(t, zero.BasicCheck())
	require.Equal(t, make([]byte, 32), decodeAnchor(t, encodeAnchor(t, zero)).RootHash)

	body := encodeAnchor(t, setPayload(addr, 32, "", 0, 1))
	withTime := append(append([]byte{}, body...), bytes.Repeat([]byte{0xAB}, 16)...)
	reader := bytes.NewReader(withTime)
	decoded := &payload.AnchorPayload{}
	require.NoError(t, decoded.Decode(payload.DecodeContext{}, reader))
	require.Equal(t, 16, reader.Len())
	require.Equal(t, amount.Amount(1), decoded.Deposit)

	aboveInt64 := append([]byte{}, body[:len(body)-1]...)
	aboveInt64 = append(aboveInt64, varInt(uint64(1)<<63)...)
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(aboveInt64)))

	treasurySet := append([]byte{0x00, payload.AnchorActionSet}, body[22:]...)
	treasury := &payload.AnchorPayload{}
	require.NoError(t, treasury.Decode(payload.DecodeContext{}, bytes.NewReader(treasurySet)))
	require.True(t, treasury.Signer().IsTreasuryAddress())
	require.Error(t, treasury.BasicCheck())

	treasuryDelete := &payload.AnchorPayload{}
	err := treasuryDelete.Decode(payload.DecodeContext{}, bytes.NewReader(
		[]byte{0x00, payload.AnchorActionDelete}))
	require.NoError(t, err)
	require.Equal(t, amount.Amount(0), treasuryDelete.Value())
	require.Error(t, treasuryDelete.BasicCheck())
}

func TestFuzzAnchorSeeds(t *testing.T) {
	setData := encodeAnchor(t, setPayload(blsAddr(), 32, "abc", 4, 300))
	delData := encodeAnchor(t, &payload.AnchorPayload{From: blsAddr(), Action: payload.AnchorActionDelete})
	require.NoError(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(setData)))
	require.NoError(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(delData)))
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(setData[:10])))

	action := append([]byte{}, delData...)
	action[len(action)-1] = 2
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(action)))

	over := append([]byte{}, setData[:len(setData)-2]...)
	over = append(over, bytes.Repeat([]byte{0x80}, 10)...)
	require.Error(t, (&payload.AnchorPayload{}).Decode(payload.DecodeContext{}, bytes.NewReader(over)))
}

func FuzzAnchorDecode(f *testing.F) {
	setData := encodeAnchor(f, setPayload(blsAddr(), 32, "abc", 4, 300))
	delData := encodeAnchor(f, &payload.AnchorPayload{From: blsAddr(), Action: payload.AnchorActionDelete})
	action := append([]byte{}, delData...)
	action[len(action)-1] = 2
	over := append([]byte{}, setData[:len(setData)-2]...)
	over = append(over, bytes.Repeat([]byte{0x80}, 10)...)

	f.Add(setData)
	f.Add(delData)
	f.Add(setData[:10])
	f.Add(action)
	f.Add(append(append([]byte{}, delData...), setData...))
	f.Add(over)

	f.Fuzz(func(t *testing.T, data []byte) {
		pld := &payload.AnchorPayload{}
		reader := bytes.NewReader(data)
		err := pld.Decode(payload.DecodeContext{}, reader)
		if err != nil {
			require.Error(t, pld.BasicCheck())

			return
		}

		buf := bytes.NewBuffer(nil)
		require.NoError(t, pld.Encode(buf))
		encoded := buf.Bytes()
		require.Len(t, encoded, pld.SerializeSize())
		require.Equal(t, data[:len(encoded)], encoded)

		again := &payload.AnchorPayload{}
		require.NoError(t, again.Decode(payload.DecodeContext{}, bytes.NewReader(data[:len(encoded)])))
		require.Equal(t, pld.Action, again.Action)
		require.Equal(t, pld.RootHash, again.RootHash)
		require.Equal(t, pld.ManifestURI, again.ManifestURI)
		require.True(t, pld.Action == payload.AnchorActionSet || pld.Action == payload.AnchorActionDelete)
		if pld.Action == payload.AnchorActionDelete {
			require.Empty(t, pld.RootHash)
		} else {
			require.GreaterOrEqual(t, len(pld.RootHash), 32)
			require.LessOrEqual(t, len(pld.RootHash), 64)
		}
		require.LessOrEqual(t, len(pld.ManifestURI), 128)
		require.True(t, utf8.ValidString(pld.ManifestURI))

		if len(pld.RootHash) > 0 {
			pld.RootHash[0] ^= 0xff
			require.NotEqual(t, pld.RootHash[0], again.RootHash[0])
		}

		if reader.Len() != 0 {
			require.NotEqual(t, len(data), len(encoded))
		} else {
			require.Equal(t, data, encoded)
		}
	})
}

func anchorAccepted(data []byte) bool {
	pld := &payload.AnchorPayload{}
	reader := bytes.NewReader(data)
	if err := pld.Decode(payload.DecodeContext{}, reader); err != nil {
		return false
	}
	if reader.Len() != 0 {
		return false
	}

	return pld.BasicCheck() == nil
}

func varInt(val uint64) []byte {
	buf := bytes.NewBuffer(nil)
	_ = encoding.WriteVarInt(buf, val)

	return buf.Bytes()
}
