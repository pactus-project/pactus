package account_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
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

func TestNoAnchorRecord(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(2)

	require.False(t, acc.HasAnchor())
	require.Equal(t, amount.Amount(0), acc.LockedDeposit())
	require.Nil(t, acc.RootHash())
	require.Empty(t, acc.ManifestURI())
	require.Equal(t, uint8(0), acc.AnchorType())
	require.Equal(t, types.Height(0), acc.CreatedAtHeight())
	require.Equal(t, uint32(0), acc.CreatedAtTime())
	require.Equal(t, types.Height(0), acc.UpdatedAtHeight())
	require.Equal(t, uint32(0), acc.UpdatedAtTime())

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Len(t, encoded, 12)
	require.Equal(t, 12, acc.SerializeSize())

	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	require.Equal(t, plain, encoded)
}

func TestAnchorRoundTrip(t *testing.T) {
	cases := []account.AnchorData{
		{
			RootHash:        bytes.Repeat([]byte{0x11}, 32),
			LockedDeposit:   1,
			AnchorType:      0x00,
			CreatedAtHeight: 10,
			CreatedAtTime:   20,
			UpdatedAtHeight: 30,
			UpdatedAtTime:   40,
		},
		{
			RootHash:        bytes.Repeat([]byte{0x22}, 64),
			ManifestURI:     strings.Repeat("a", 128),
			AnchorType:      0xFF,
			LockedDeposit:   amount.MaxNanoPAC,
			CreatedAtHeight: 1,
			CreatedAtTime:   2,
			UpdatedAtHeight: 3,
			UpdatedAtTime:   4,
		},
		{
			RootHash:        bytes.Repeat([]byte{0x33}, 48),
			ManifestURI:     "",
			AnchorType:      0x04,
			LockedDeposit:   1,
			CreatedAtHeight: 9,
			CreatedAtTime:   8,
			UpdatedAtHeight: 7,
			UpdatedAtTime:   6,
		},
	}

	for _, anchor := range cases {
		plain := account.NewAccount(7)
		plain.AddToBalance(9)

		anchored := plain.Clone()
		require.NoError(t, anchored.SetAnchor(anchor))
		require.True(t, anchored.HasAnchor())
		require.Equal(t, anchor.RootHash, anchored.RootHash())
		require.Equal(t, anchor.ManifestURI, anchored.ManifestURI())
		require.Equal(t, anchor.AnchorType, anchored.AnchorType())
		require.Equal(t, anchor.LockedDeposit, anchored.LockedDeposit())
		require.Equal(t, anchor.CreatedAtHeight, anchored.CreatedAtHeight())
		require.Equal(t, anchor.CreatedAtTime, anchored.CreatedAtTime())
		require.Equal(t, anchor.UpdatedAtHeight, anchored.UpdatedAtHeight())
		require.Equal(t, anchor.UpdatedAtTime, anchored.UpdatedAtTime())

		encoded, err := anchored.Bytes()
		require.NoError(t, err)
		require.Len(t, encoded, anchored.SerializeSize())
		require.Greater(t, len(encoded), 12)
		require.Equal(t, byte(1), encoded[12])

		decoded, err := account.FromBytes(encoded)
		require.NoError(t, err)
		require.Equal(t, anchored, decoded)

		again, err := decoded.Bytes()
		require.NoError(t, err)
		require.Equal(t, encoded, again)
		require.NotEqual(t, plain.Hash(), anchored.Hash())
	}
}

func TestAnchorExactSizes(t *testing.T) {
	short := account.NewAccount(1)
	require.NoError(t, short.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x11}, 32),
		LockedDeposit: 1,
	}))
	shortBytes, err := short.Bytes()
	require.NoError(t, err)
	require.Len(t, shortBytes, 72)
	require.Equal(t, 72, short.SerializeSize())

	longURI := strings.Repeat("a", 128)
	longAcc := account.NewAccount(1)
	require.NoError(t, longAcc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x22}, 64),
		ManifestURI:   longURI,
		LockedDeposit: 1,
	}))
	longBytes, err := longAcc.Bytes()
	require.NoError(t, err)
	require.Len(t, longBytes, 232)
	require.Equal(t, 232, longAcc.SerializeSize())
}

func TestAnchorGoldenVector(t *testing.T) {
	data, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(t, err)
	require.Len(t, data, 72)

	acc, err := account.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, int32(1), acc.Number())
	require.Equal(t, amount.Amount(2), acc.Balance())
	require.True(t, acc.HasAnchor())
	require.Equal(t, bytes.Repeat([]byte{0x11}, 32), acc.RootHash())
	require.Empty(t, acc.ManifestURI())
	require.Equal(t, uint8(0), acc.AnchorType())
	require.Equal(t, amount.Amount(1), acc.LockedDeposit())
	require.Equal(t, types.Height(10), acc.CreatedAtHeight())
	require.Equal(t, uint32(20), acc.CreatedAtTime())
	require.Equal(t, types.Height(30), acc.UpdatedAtHeight())
	require.Equal(t, uint32(40), acc.UpdatedAtTime())

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, data, encoded)
	require.Equal(t, hash.CalcHash(data), acc.Hash())
	expected, err := hash.FromString("a97eb6557a257b9c2e8eb44e950c4918997eb7fab9ead0e5b3a7b67249ae04a6")
	require.NoError(t, err)
	require.Equal(t, expected, acc.Hash())
	require.Equal(t, len(data), acc.SerializeSize())

	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	require.Equal(t, plain, data[:12])
	require.Equal(t, byte(1), data[12])
}

func TestZeroRootHashIsAnAnchor(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(2)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      make([]byte, 32),
		LockedDeposit: 1,
	}))

	require.True(t, acc.HasAnchor())
	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Greater(t, len(encoded), 12)

	decoded, err := account.FromBytes(encoded)
	require.NoError(t, err)
	require.Equal(t, make([]byte, 32), decoded.RootHash())
	require.True(t, decoded.HasAnchor())
}

func TestAnchorTypeIsOpaque(t *testing.T) {
	for _, anchorType := range []uint8{0x00, 0x04, 0xFF} {
		acc := account.NewAccount(1)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0x44}, 32),
			AnchorType:    anchorType,
			LockedDeposit: 1,
		}))

		decoded := mustRoundTrip(t, acc)
		require.Equal(t, anchorType, decoded.AnchorType())
	}
}

func TestCreatedAfterUpdatedIsAccepted(t *testing.T) {
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        bytes.Repeat([]byte{0x55}, 32),
		LockedDeposit:   1,
		CreatedAtHeight: 100,
		UpdatedAtHeight: 1,
		CreatedAtTime:   50,
		UpdatedAtTime:   2,
	}))

	decoded := mustRoundTrip(t, acc)
	require.Equal(t, types.Height(100), decoded.CreatedAtHeight())
	require.Equal(t, types.Height(1), decoded.UpdatedAtHeight())
}

func TestTimestampsRoundTripInPlace(t *testing.T) {
	anchor := account.AnchorData{
		RootHash:        bytes.Repeat([]byte{0x66}, 32),
		LockedDeposit:   1,
		CreatedAtHeight: 11,
		CreatedAtTime:   22,
		UpdatedAtHeight: 33,
		UpdatedAtTime:   44,
	}
	acc := account.NewAccount(4)
	acc.AddToBalance(5)
	require.NoError(t, acc.SetAnchor(anchor))

	decoded := mustRoundTrip(t, acc)
	require.Equal(t, types.Height(11), decoded.CreatedAtHeight())
	require.Equal(t, uint32(22), decoded.CreatedAtTime())
	require.Equal(t, types.Height(33), decoded.UpdatedAtHeight())
	require.Equal(t, uint32(44), decoded.UpdatedAtTime())

	changed := anchor
	changed.UpdatedAtTime = 45
	other := account.NewAccount(4)
	other.AddToBalance(5)
	require.NoError(t, other.SetAnchor(changed))
	require.NotEqual(t, acc.Hash(), other.Hash())
	require.Equal(t, acc.CreatedAtHeight(), other.CreatedAtHeight())
	require.Equal(t, acc.CreatedAtTime(), other.CreatedAtTime())
	require.Equal(t, acc.UpdatedAtHeight(), other.UpdatedAtHeight())
}

func TestBalanceDoesNotIncludeLockedDeposit(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(2)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x77}, 32),
		LockedDeposit: 9,
	}))

	require.Equal(t, amount.Amount(2), acc.Balance())
	require.Equal(t, amount.Amount(9), acc.LockedDeposit())

	acc.AddToBalance(3)
	require.Equal(t, amount.Amount(5), acc.Balance())
	require.Equal(t, amount.Amount(9), acc.LockedDeposit())
	require.True(t, acc.HasAnchor())
}

func TestMaxBalanceAndMaxDepositStaySeparate(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(amount.MaxNanoPAC)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x88}, 32),
		LockedDeposit: amount.MaxNanoPAC,
	}))

	decoded := mustRoundTrip(t, acc)
	require.Equal(t, amount.Amount(amount.MaxNanoPAC), decoded.Balance())
	require.Equal(t, amount.Amount(amount.MaxNanoPAC), decoded.LockedDeposit())
}

func TestNegativeBalanceRoundTrip(t *testing.T) {
	acc := account.NewAccount(1)
	acc.SubtractFromBalance(5)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x99}, 32),
		LockedDeposit: 1,
	}))

	decoded := mustRoundTrip(t, acc)
	require.Equal(t, amount.Amount(-5), decoded.Balance())
	require.Equal(t, amount.Amount(1), decoded.LockedDeposit())
}

func TestDepositBounds(t *testing.T) {
	valid := []amount.Amount{1, amount.MaxNanoPAC}
	for _, deposit := range valid {
		acc := account.NewAccount(1)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0xab}, 32),
			LockedDeposit: deposit,
		}))
		decoded := mustRoundTrip(t, acc)
		require.Equal(t, deposit, decoded.LockedDeposit())
	}

	invalid := []amount.Amount{0, -1, amount.Amount(math.MinInt64), amount.MaxNanoPAC + 1}
	for _, deposit := range invalid {
		acc := account.NewAccount(1)
		err := acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0xab}, 32),
			LockedDeposit: deposit,
		})
		require.ErrorIs(t, err, account.ErrInvalidAnchorDeposit)
		require.False(t, acc.HasAnchor())

		suffix := baseSuffix(bytes.Repeat([]byte{0xab}, 32))
		suffix.deposit = int64(deposit)
		encoded := rawAccount(t, &suffix)
		_, err = account.FromBytes(encoded)
		require.ErrorIs(t, err, account.ErrInvalidAnchorDeposit)
	}
}

func TestSetAnchorRejectsUndecodableInput(t *testing.T) {
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x01}, 32),
		ManifestURI:   "keep",
		LockedDeposit: 1,
	}))

	cases := []account.AnchorData{
		{RootHash: bytes.Repeat([]byte{0x01}, 31), LockedDeposit: 1},
		{RootHash: bytes.Repeat([]byte{0x01}, 65), LockedDeposit: 1},
		{RootHash: bytes.Repeat([]byte{0x01}, 32), ManifestURI: strings.Repeat("a", 129), LockedDeposit: 1},
		{RootHash: bytes.Repeat([]byte{0x01}, 32), ManifestURI: "\xC3", LockedDeposit: 1},
		{RootHash: bytes.Repeat([]byte{0x01}, 32), LockedDeposit: 0},
	}
	for _, anchor := range cases {
		err := acc.SetAnchor(anchor)
		require.Error(t, err)
		require.Equal(t, "keep", acc.ManifestURI())
		require.Equal(t, bytes.Repeat([]byte{0x01}, 32), acc.RootHash())
	}

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, byte(1), encoded[12])
}

func TestFromBytesRejectsAnchor(t *testing.T) {
	hash32 := bytes.Repeat([]byte{0x11}, 32)

	flagZero := baseSuffix(hash32)
	flagZero.flag = 0
	_, err := account.FromBytes(rawAccount(t, &flagZero))
	require.ErrorIs(t, err, account.ErrInvalidAnchorFlag)

	flagTwo := baseSuffix(hash32)
	flagTwo.flag = 2
	_, err = account.FromBytes(rawAccount(t, &flagTwo))
	require.ErrorIs(t, err, account.ErrInvalidAnchorFlag)

	flagHigh := baseSuffix(hash32)
	flagHigh.flag = 0xFF
	_, err = account.FromBytes(rawAccount(t, &flagHigh))
	require.ErrorIs(t, err, account.ErrInvalidAnchorFlag)

	shortHash := baseSuffix(bytes.Repeat([]byte{0x11}, 31))
	_, err = account.FromBytes(rawAccount(t, &shortHash))
	require.ErrorIs(t, err, account.ErrAnchorHashLength)

	longHash := baseSuffix(bytes.Repeat([]byte{0x11}, 65))
	_, err = account.FromBytes(rawAccount(t, &longHash))
	require.ErrorIs(t, err, account.ErrAnchorHashLength)

	longURI := baseSuffix(hash32)
	longURI.uri = bytes.Repeat([]byte{'a'}, 129)
	_, err = account.FromBytes(rawAccount(t, &longURI))
	require.ErrorIs(t, err, account.ErrAnchorURILength)

	badUTF8 := baseSuffix(hash32)
	badUTF8.uri = []byte{0xC3}
	_, err = account.FromBytes(rawAccount(t, &badUTF8))
	require.ErrorIs(t, err, account.ErrInvalidAnchorURI)

	validSuffix := baseSuffix(hash32)
	validSuffix.uri = []byte("ok")
	valid := rawAccount(t, &validSuffix)
	_, err = account.FromBytes(append(bytes.Clone(valid), 0x99))
	require.ErrorIs(t, err, account.ErrAnchorTrailingBytes)
}

func TestTruncatedAccount(t *testing.T) {
	suffix := baseSuffix(bytes.Repeat([]byte{0x11}, 32))
	suffix.uri = []byte("ok")
	suffix.anchorType = 7
	suffix.createdHeight = 10
	suffix.createdTime = 20
	suffix.updatedHeight = 30
	suffix.updatedTime = 40
	valid := rawAccount(t, &suffix)
	require.Greater(t, len(valid), 12)

	_, err := account.FromBytes(nil)
	require.Error(t, err)
	_, err = account.FromBytes(valid[:11])
	require.Error(t, err)

	header, err := account.FromBytes(valid[:12])
	require.NoError(t, err)
	require.False(t, header.HasAnchor())
	require.Equal(t, int32(1), header.Number())
	require.Equal(t, amount.Amount(2), header.Balance())

	for length := 13; length < len(valid); length++ {
		_, err = account.FromBytes(valid[:length])
		require.Error(t, err, "length %d", length)
	}

	_, err = account.FromBytes(valid)
	require.NoError(t, err)
}

func TestTrailingZeroAndTwoSuffixes(t *testing.T) {
	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	_, err = account.FromBytes(append(bytes.Clone(plain), 0x00))
	require.ErrorIs(t, err, account.ErrInvalidAnchorFlag)

	anchored, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(t, err)
	two := append(bytes.Clone(anchored), anchored[12:]...)
	_, err = account.FromBytes(two)
	require.ErrorIs(t, err, account.ErrAnchorTrailingBytes)
}

func TestDeclaredLengthWithoutPayloadDoesNotPanic(t *testing.T) {
	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)

	hashLen255 := append(bytes.Clone(plain), 1, 255)
	require.NotPanics(t, func() {
		_, err = account.FromBytes(hashLen255)
		require.ErrorIs(t, err, account.ErrAnchorHashLength)
	})

	hash32 := bytes.Repeat([]byte{0x11}, 32)
	uriLen255 := append(append(bytes.Clone(plain), 1, 32), hash32...)
	uriLen255 = append(uriLen255, 255)
	require.NotPanics(t, func() {
		_, err = account.FromBytes(uriLen255)
		require.ErrorIs(t, err, account.ErrAnchorURILength)
	})
}

func TestURILengthCountsBytes(t *testing.T) {
	oneRune := strings.Repeat("é", 64)
	require.Len(t, oneRune, 128)
	require.Len(t, []rune(oneRune), 64)

	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x12}, 32),
		ManifestURI:   oneRune,
		LockedDeposit: 1,
	}))
	decoded := mustRoundTrip(t, acc)
	require.Equal(t, oneRune, decoded.ManifestURI())

	tooManyBytes := oneRune + "a"
	require.Len(t, tooManyBytes, 129)
	err := acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x12}, 32),
		ManifestURI:   tooManyBytes,
		LockedDeposit: 1,
	})
	require.ErrorIs(t, err, account.ErrAnchorURILength)
	require.Equal(t, oneRune, acc.ManifestURI())

	manyRunes := strings.Repeat("é", 128)
	require.Len(t, []rune(manyRunes), 128)
	require.Len(t, manyRunes, 256)
	err = account.NewAccount(1).SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x12}, 32),
		ManifestURI:   manyRunes,
		LockedDeposit: 1,
	})
	require.ErrorIs(t, err, account.ErrAnchorURILength)
}

func TestHostileURIIsStoredVerbatim(t *testing.T) {
	uris := []string{
		"line\nbreak",
		"has\x00nul",
		"javascript:alert(1)",
		"caf\u00e9",
		"cafe\u0301",
	}
	require.NotEqual(t, uris[3], uris[4])

	for _, uri := range uris {
		acc := account.NewAccount(1)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0x13}, 32),
			ManifestURI:   uri,
			LockedDeposit: 1,
		}))
		decoded := mustRoundTrip(t, acc)
		require.Equal(t, uri, decoded.ManifestURI())
	}
}

func TestSetAnchorCopiesInput(t *testing.T) {
	root := bytes.Repeat([]byte{0x21}, 32)
	anchor := account.AnchorData{
		RootHash:      root,
		ManifestURI:   "first",
		AnchorType:    3,
		LockedDeposit: 4,
	}
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(anchor))

	root[0] = 0
	anchor.RootHash[0] = 0
	anchor.ManifestURI = "mutated"
	anchor.AnchorType = 9
	anchor.LockedDeposit = 8

	require.Equal(t, byte(0x21), acc.RootHash()[0])
	require.Equal(t, "first", acc.ManifestURI())
	require.Equal(t, uint8(3), acc.AnchorType())
	require.Equal(t, amount.Amount(4), acc.LockedDeposit())

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, byte(0x21), encoded[14])
}

func TestRootHashGetterReturnsACopy(t *testing.T) {
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x31}, 32),
		LockedDeposit: 1,
	}))

	got := acc.RootHash()
	got[0] = 0
	require.Equal(t, byte(0x31), acc.RootHash()[0])
}

func TestSecondSetAnchorReplacesTheSlot(t *testing.T) {
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x41}, 32),
		ManifestURI:   "old-uri",
		AnchorType:    1,
		LockedDeposit: 2,
	}))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        bytes.Repeat([]byte{0x42}, 64),
		ManifestURI:     "new-uri",
		AnchorType:      2,
		LockedDeposit:   3,
		CreatedAtHeight: 4,
		CreatedAtTime:   5,
		UpdatedAtHeight: 6,
		UpdatedAtTime:   7,
	}))

	require.Equal(t, "new-uri", acc.ManifestURI())
	require.Equal(t, bytes.Repeat([]byte{0x42}, 64), acc.RootHash())
	require.Equal(t, uint8(2), acc.AnchorType())
	require.Equal(t, amount.Amount(3), acc.LockedDeposit())
	require.Equal(t, types.Height(4), acc.CreatedAtHeight())
	require.Equal(t, uint32(7), acc.UpdatedAtTime())

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.NotContains(t, encoded, []byte("old-uri"))
}

func TestClearAnchor(t *testing.T) {
	empty := account.NewAccount(1)
	empty.AddToBalance(2)
	empty.ClearAnchor()
	emptyBytes, err := empty.Bytes()
	require.NoError(t, err)
	require.Len(t, emptyBytes, 12)
	require.False(t, empty.HasAnchor())

	anchored := empty.Clone()
	require.NoError(t, anchored.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x51}, 32),
		LockedDeposit: 1,
	}))
	anchored.ClearAnchor()
	cleared, err := anchored.Bytes()
	require.NoError(t, err)
	require.Equal(t, emptyBytes, cleared)
	require.False(t, anchored.HasAnchor())
	require.Equal(t, amount.Amount(0), anchored.LockedDeposit())

	require.NoError(t, anchored.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x52}, 32),
		ManifestURI:   "after-clear",
		LockedDeposit: 1,
	}))
	rebuilt, err := anchored.Bytes()
	require.NoError(t, err)
	require.Greater(t, len(rebuilt), 12)
	require.Equal(t, "after-clear", anchored.ManifestURI())
	require.False(t, bytes.Contains(rebuilt, bytes.Repeat([]byte{0x51}, 32)))
	require.True(t, bytes.Contains(rebuilt, bytes.Repeat([]byte{0x52}, 32)))
}

func TestCloneAnchor(t *testing.T) {
	plain := account.NewAccount(8)
	plain.AddToBalance(3)
	clonedPlain := plain.Clone()
	require.False(t, clonedPlain.HasAnchor())
	plainBytes, err := plain.Bytes()
	require.NoError(t, err)
	clonedBytes, err := clonedPlain.Bytes()
	require.NoError(t, err)
	require.Equal(t, plainBytes, clonedBytes)
	require.Len(t, clonedBytes, 12)

	original := plain.Clone()
	require.NoError(t, original.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x61}, 32),
		ManifestURI:   "origin",
		LockedDeposit: 6,
	}))
	cloned := original.Clone()
	cloned.AddToBalance(1)
	require.NoError(t, cloned.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x62}, 32),
		ManifestURI:   "clone",
		LockedDeposit: 7,
	}))

	require.Equal(t, amount.Amount(3), original.Balance())
	require.Equal(t, bytes.Repeat([]byte{0x61}, 32), original.RootHash())
	require.Equal(t, "origin", original.ManifestURI())
	require.Equal(t, amount.Amount(6), original.LockedDeposit())

	cloned.ClearAnchor()
	require.True(t, original.HasAnchor())
	require.False(t, cloned.HasAnchor())
}

func TestAppendToBytesDoesNotChangeNextEncode(t *testing.T) {
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x71}, 32),
		LockedDeposit: 1,
	}))

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	kept := bytes.Clone(encoded)
	_ = append(encoded, 0xFF)

	next, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, kept, next)
}

func TestAnchoredPrefixMatchesPlainAccount(t *testing.T) {
	plain := account.NewAccount(1)
	plain.AddToBalance(2)
	plainBytes, err := plain.Bytes()
	require.NoError(t, err)

	anchored := plain.Clone()
	require.NoError(t, anchored.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x81}, 32),
		ManifestURI:   "x",
		LockedDeposit: 1,
		AnchorType:    9,
	}))
	encoded, err := anchored.Bytes()
	require.NoError(t, err)
	require.Equal(t, plainBytes, encoded[:12])
	require.Equal(t, byte(1), encoded[12])
}

func TestAnchorWireLayout(t *testing.T) {
	number := int32(0x01020304)
	balance := int64(0x0807060504030201)
	deposit := int64(0x0000000100000001)
	root := make([]byte, 32)
	for i := range root {
		root[i] = byte(i + 1)
	}
	uri := "\u00e9"
	require.Equal(t, []byte{0xC3, 0xA9}, []byte(uri))

	createdHeight := uint32(0xA1B2C3D4)
	createdTime := uint32(0x11111111)
	updatedHeight := uint32(0x22222222)
	updatedTime := uint32(0xFFFFFFFF)
	anchorType := byte(0xAB)

	acc := account.NewAccount(number)
	acc.AddToBalance(amount.Amount(balance))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		ManifestURI:     uri,
		AnchorType:      anchorType,
		LockedDeposit:   amount.Amount(deposit),
		CreatedAtHeight: types.Height(createdHeight),
		CreatedAtTime:   createdTime,
		UpdatedAtHeight: types.Height(updatedHeight),
		UpdatedAtTime:   updatedTime,
	}))

	encoded, err := acc.Bytes()
	require.NoError(t, err)

	expected := new(bytes.Buffer)
	require.NoError(t, binary.Write(expected, binary.LittleEndian, number))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, balance))
	require.NoError(t, expected.WriteByte(1))
	require.NoError(t, expected.WriteByte(byte(len(root))))
	_, err = expected.Write(root)
	require.NoError(t, err)
	require.NoError(t, expected.WriteByte(byte(len(uri))))
	_, err = expected.WriteString(uri)
	require.NoError(t, err)
	require.NoError(t, expected.WriteByte(anchorType))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, deposit))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, createdHeight))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, createdTime))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, updatedHeight))
	require.NoError(t, binary.Write(expected, binary.LittleEndian, updatedTime))
	require.Equal(t, expected.Bytes(), encoded)

	encoded[14] ^= 0xFF
	require.Equal(t, root, acc.RootHash())
	require.Equal(t, uri, acc.ManifestURI())

	wire := bytes.Clone(expected.Bytes())
	decoded, err := account.FromBytes(wire)
	require.NoError(t, err)
	wire[14] ^= 0xFF
	wire[14+len(root)+1] ^= 0xFF
	require.Equal(t, root, decoded.RootHash())
	require.Equal(t, uri, decoded.ManifestURI())
	require.Equal(t, number, decoded.Number())
	require.Equal(t, amount.Amount(balance), decoded.Balance())
	require.Equal(t, amount.Amount(deposit), decoded.LockedDeposit())
	require.Equal(t, anchorType, decoded.AnchorType())
	require.Equal(t, types.Height(createdHeight), decoded.CreatedAtHeight())
	require.Equal(t, createdTime, decoded.CreatedAtTime())
	require.Equal(t, types.Height(updatedHeight), decoded.UpdatedAtHeight())
	require.Equal(t, updatedTime, decoded.UpdatedAtTime())
}

func TestHashLengthInteriorValues(t *testing.T) {
	for _, length := range []int{33, 63} {
		root := bytes.Repeat([]byte{0x5A}, length)
		acc := account.NewAccount(1)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      root,
			LockedDeposit: 1,
		}))
		decoded := mustRoundTrip(t, acc)
		require.Equal(t, root, decoded.RootHash())
	}
}

func TestSetAnchorRejectsEmptyHash(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(4)
	require.NotPanics(t, func() {
		err := acc.SetAnchor(account.AnchorData{LockedDeposit: 1})
		require.ErrorIs(t, err, account.ErrAnchorHashLength)
	})
	require.False(t, acc.HasAnchor())
	require.Equal(t, amount.Amount(4), acc.Balance())

	err := acc.SetAnchor(account.AnchorData{
		RootHash:      []byte{},
		LockedDeposit: 1,
	})
	require.ErrorIs(t, err, account.ErrAnchorHashLength)

	hashLenZero := baseSuffix(nil)
	decoded, err := account.FromBytes(rawAccount(t, &hashLenZero))
	require.ErrorIs(t, err, account.ErrAnchorHashLength)
	require.Nil(t, decoded)
}

func TestInvalidUTF8Shapes(t *testing.T) {
	cases := [][]byte{
		{0xC0, 0x80},
		{0xED, 0xA0, 0x80},
		{0xFF},
		append(bytes.Repeat([]byte{'a'}, 127), 0xC3),
	}
	for _, uri := range cases {
		acc := account.NewAccount(1)
		err := acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0x12}, 32),
			ManifestURI:   string(uri),
			LockedDeposit: 1,
		})
		require.ErrorIs(t, err, account.ErrInvalidAnchorURI)
		require.False(t, acc.HasAnchor())

		suffix := baseSuffix(bytes.Repeat([]byte{0x12}, 32))
		suffix.uri = uri
		decoded, err := account.FromBytes(rawAccount(t, &suffix))
		require.ErrorIs(t, err, account.ErrInvalidAnchorURI)
		require.Nil(t, decoded)
	}

	emoji := "\U0001F600"
	require.Len(t, emoji, 4)
	body := strings.Repeat("a", 124) + emoji
	require.Len(t, body, 128)
	acc := account.NewAccount(1)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x12}, 32),
		ManifestURI:   body,
		LockedDeposit: 1,
	}))
	decoded := mustRoundTrip(t, acc)
	require.Equal(t, body, decoded.ManifestURI())
}

func TestExtremeBalancesRoundTrip(t *testing.T) {
	cases := []struct {
		number  int32
		balance amount.Amount
	}{
		{number: math.MinInt32, balance: amount.Amount(math.MinInt64)},
		{number: math.MaxInt32, balance: amount.Amount(math.MaxInt64)},
	}
	for _, testCase := range cases {
		acc := account.NewAccount(testCase.number)
		acc.AddToBalance(testCase.balance)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0x77}, 32),
			LockedDeposit: 1,
		}))
		decoded := mustRoundTrip(t, acc)
		require.Equal(t, testCase.number, decoded.Number())
		require.Equal(t, testCase.balance, decoded.Balance())
		require.Equal(t, amount.Amount(1), decoded.LockedDeposit())

		encoded, err := decoded.Bytes()
		require.NoError(t, err)
		reader := bytes.NewReader(encoded[:12])
		var gotNumber int32
		var gotBalance int64
		require.NoError(t, binary.Read(reader, binary.LittleEndian, &gotNumber))
		require.NoError(t, binary.Read(reader, binary.LittleEndian, &gotBalance))
		require.Equal(t, testCase.number, gotNumber)
		require.Equal(t, int64(testCase.balance), gotBalance)
	}
}

func TestZeroHeaderWithAnchor(t *testing.T) {
	acc := account.NewAccount(0)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x10}, 32),
		LockedDeposit: 1,
	}))
	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, make([]byte, 12), encoded[:12])
	require.Equal(t, byte(1), encoded[12])

	decoded := mustRoundTrip(t, acc)
	require.Equal(t, int32(0), decoded.Number())
	require.Equal(t, amount.Amount(0), decoded.Balance())
	require.True(t, decoded.HasAnchor())
}

func TestSubtractFromBalanceKeepsAnchor(t *testing.T) {
	acc := account.NewAccount(1)
	acc.AddToBalance(10)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{0x21}, 32),
		ManifestURI:   "stay",
		LockedDeposit: 4,
	}))
	acc.SubtractFromBalance(3)
	require.Equal(t, amount.Amount(7), acc.Balance())
	require.Equal(t, amount.Amount(4), acc.LockedDeposit())
	require.Equal(t, "stay", acc.ManifestURI())
	require.True(t, acc.HasAnchor())
}

func TestFromBytesErrorReturnsNoAccount(t *testing.T) {
	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	anchored, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(t, err)
	flagFF := baseSuffix(bytes.Repeat([]byte{0x11}, 32))
	flagFF.flag = 0xFF

	inputs := [][]byte{
		nil,
		{},
		anchored[:11],
		append(bytes.Clone(plain), 0x00),
		append(bytes.Clone(anchored), anchored[12:]...),
		anchored[:len(anchored)-1],
		rawAccount(t, &flagFF),
	}
	for _, input := range inputs {
		acc, err := account.FromBytes(input)
		require.Error(t, err)
		require.Nil(t, acc)
	}
}

func TestFuzzSeeds(t *testing.T) {
	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(t, err)
	anchored, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(t, err)

	valid := [][]byte{plain, anchored}
	for _, seed := range valid {
		acc, err := account.FromBytes(seed)
		require.NoError(t, err)
		encoded, err := acc.Bytes()
		require.NoError(t, err)
		require.Equal(t, seed, encoded)
	}

	invalid := [][]byte{
		{},
		anchored[:11],
		append(bytes.Clone(plain), 0x00),
		append(bytes.Clone(anchored), anchored[12:]...),
	}
	for _, seed := range invalid {
		acc, err := account.FromBytes(seed)
		require.Error(t, err)
		require.Nil(t, acc)
	}
}

func FuzzFromBytes(f *testing.F) {
	plain, err := hex.DecodeString(plainAccountHex)
	require.NoError(f, err)
	anchored, err := hex.DecodeString(anchoredAccountHex)
	require.NoError(f, err)

	f.Add(plain)
	f.Add(anchored)
	f.Add(anchored[:11])
	f.Add([]byte{})
	f.Add(append(bytes.Clone(plain), 0x00))
	f.Add(append(bytes.Clone(anchored), anchored[12:]...))

	f.Fuzz(func(t *testing.T, data []byte) {
		acc, err := account.FromBytes(data)
		if err != nil {
			require.Nil(t, acc)

			return
		}

		require.NotNil(t, acc)
		encoded, err := acc.Bytes()
		require.NoError(t, err)
		require.Equal(t, data, encoded)
		require.Equal(t, len(encoded), acc.SerializeSize())

		again, err := account.FromBytes(encoded)
		require.NoError(t, err)
		require.Equal(t, acc, again)

		if acc.HasAnchor() {
			require.Greater(t, len(encoded), 12)
			require.Equal(t, byte(1), encoded[12])
			require.Positive(t, int64(acc.LockedDeposit()))
			require.LessOrEqual(t, acc.LockedDeposit(), amount.Amount(amount.MaxNanoPAC))

			root := acc.RootHash()
			require.GreaterOrEqual(t, len(root), 32)
			require.LessOrEqual(t, len(root), 64)
			require.LessOrEqual(t, len(acc.ManifestURI()), 128)
			require.True(t, utf8.ValidString(acc.ManifestURI()))
			root[0] ^= 0xFF
			require.NotEqual(t, root[0], acc.RootHash()[0])

			return
		}

		require.Len(t, encoded, 12)
		require.Equal(t, amount.Amount(0), acc.LockedDeposit())
	})
}

func mustRoundTrip(t *testing.T, acc *account.Account) *account.Account {
	t.Helper()

	encoded, err := acc.Bytes()
	require.NoError(t, err)
	require.Len(t, encoded, acc.SerializeSize())

	decoded, err := account.FromBytes(encoded)
	require.NoError(t, err)
	require.Equal(t, acc, decoded)

	return decoded
}

type rawSuffix struct {
	flag          byte
	root          []byte
	uri           []byte
	anchorType    byte
	deposit       int64
	createdHeight uint32
	createdTime   uint32
	updatedHeight uint32
	updatedTime   uint32
}

func baseSuffix(root []byte) rawSuffix {
	return rawSuffix{
		flag:          1,
		root:          root,
		deposit:       1,
		createdHeight: 1,
		createdTime:   2,
		updatedHeight: 3,
		updatedTime:   4,
	}
}

func rawAccount(t *testing.T, suffix *rawSuffix) []byte {
	t.Helper()

	buf := bytes.NewBuffer(make([]byte, 0, 12+len(suffix.root)+len(suffix.uri)+32))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, int32(1)))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, int64(2)))
	require.NoError(t, buf.WriteByte(suffix.flag))
	require.LessOrEqual(t, len(suffix.root), 255)
	require.NoError(t, buf.WriteByte(byte(len(suffix.root))))
	_, err := buf.Write(suffix.root)
	require.NoError(t, err)
	require.LessOrEqual(t, len(suffix.uri), 255)
	require.NoError(t, buf.WriteByte(byte(len(suffix.uri))))
	_, err = buf.Write(suffix.uri)
	require.NoError(t, err)
	require.NoError(t, buf.WriteByte(suffix.anchorType))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, suffix.deposit))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, suffix.createdHeight))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, suffix.createdTime))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, suffix.updatedHeight))
	require.NoError(t, binary.Write(buf, binary.LittleEndian, suffix.updatedTime))

	return buf.Bytes()
}
