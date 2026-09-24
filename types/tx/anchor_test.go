package tx_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func signedAnchor(t *testing.T, ts *testsuite.TestSuite, deposit, fee amount.Amount) *tx.Tx {
	t.Helper()

	pub, prv := ts.RandBLSKeyPair()
	trx := tx.NewAnchorTx(
		1,
		pub.AccountAddress(),
		payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32),
		"abc",
		0x04,
		deposit,
		fee,
	)
	ts.HelperSignTransaction(prv, trx)

	return trx
}

func TestNewAnchorTxFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	pub, _ := ts.RandBLSKeyPair()
	from := pub.AccountAddress()

	free := tx.NewAnchorTx(1, from, payload.AnchorActionSet, bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 0)
	priced := tx.NewAnchorTx(1, from, payload.AnchorActionSet, bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 5)
	require.Equal(t, amount.Amount(0), free.Fee())
	require.Equal(t, amount.Amount(5), priced.Fee())
	require.False(t, free.IsFreeTx())
	require.False(t, priced.IsFreeTx())

	backing := make([]byte, 1<<20)
	copy(backing[:32], bytes.Repeat([]byte{0x11}, 32))
	window := backing[:32]
	owned := tx.NewAnchorTx(1, from, payload.AnchorActionSet, window, "", 0, 1, 1)
	window[0] = 0
	stored := owned.Payload().(*payload.AnchorPayload).RootHash
	require.Equal(t, byte(0x11), stored[0])
	require.Len(t, stored, 32)
	require.Equal(t, 32, cap(stored))

	deleted := tx.NewAnchorTx(1, from, payload.AnchorActionDelete, bytes.Repeat([]byte{0x11}, 32), "uri", 1, 9, 5)
	require.Equal(t, amount.Amount(0), deleted.Payload().Value())

	require.True(t, tx.NewUnbondTx(1, ts.RandValAddress()).IsFreeTx())
	require.True(t, tx.NewSortitionTx(1, ts.RandValAddress(), ts.RandProof()).IsFreeTx())
}

func TestAnchorTxBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	ok := signedAnchor(t, ts, 1, 0)
	require.NoError(t, ok.BasicCheck())
	require.False(t, ok.IsFreeTx())

	maxFee := signedAnchor(t, ts, 1, amount.Amount(amount.MaxNanoPAC))
	require.NoError(t, maxFee.BasicCheck())

	negative := signedAnchor(t, ts, 1, -1)
	require.Error(t, negative.BasicCheck())
	over := signedAnchor(t, ts, 1, amount.Amount(amount.MaxNanoPAC)+1)
	require.Error(t, over.BasicCheck())

	pub, _ := ts.RandBLSKeyPair()
	unsigned := tx.NewAnchorTx(1, pub.AccountAddress(), payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 0)
	raw, err := unsigned.Bytes()
	require.NoError(t, err)
	decoded, err := tx.FromBytes(raw)
	require.NoError(t, err)
	require.Error(t, decoded.BasicCheck())

	otherPub, otherPrv := ts.RandBLSKeyPair()
	mismatch := tx.NewAnchorTx(1, pub.AccountAddress(), payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 1)
	ts.HelperSignTransaction(otherPrv, mismatch)
	require.NotEqual(t, pub.AccountAddress(), otherPub.AccountAddress())
	require.Error(t, mismatch.BasicCheck())

	for _, sign := range []func() (crypto.PublicKey, crypto.PrivateKey){
		func() (crypto.PublicKey, crypto.PrivateKey) {
			keyPub, keyPrv := ts.RandBLSKeyPair()

			return keyPub, keyPrv
		},
		func() (crypto.PublicKey, crypto.PrivateKey) {
			keyPub, keyPrv := ts.RandEd25519KeyPair()

			return keyPub, keyPrv
		},
		func() (crypto.PublicKey, crypto.PrivateKey) {
			keyPub, keyPrv := ts.RandSecp256k1KeyPair()

			return keyPub, keyPrv
		},
	} {
		pubKey, prvKey := sign()
		trx := tx.NewAnchorTx(1, pubKey.AccountAddress(), payload.AnchorActionDelete,
			nil, "", 0, 0, 1)
		ts.HelperSignTransaction(prvKey, trx)
		signedRaw, err := trx.Bytes()
		require.NoError(t, err)
		roundTrip, err := tx.FromBytes(signedRaw)
		require.NoError(t, err)
		require.NoError(t, roundTrip.BasicCheck())
		require.Equal(t, payload.TypeAnchor, roundTrip.Payload().Type())
	}

	memoPub, memoPrv := ts.RandBLSKeyPair()
	locked := tx.NewAnchorTx(0, memoPub.AccountAddress(), payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 1)
	ts.HelperSignTransaction(memoPrv, locked)
	require.Error(t, locked.BasicCheck())

	memoOK := tx.NewAnchorTx(1, memoPub.AccountAddress(), payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 1, tx.WithMemo(bytesString(64)))
	ts.HelperSignTransaction(memoPrv, memoOK)
	require.NoError(t, memoOK.BasicCheck())

	memoLong := tx.NewAnchorTx(1, memoPub.AccountAddress(), payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 1, tx.WithMemo(bytesString(65)))
	ts.HelperSignTransaction(memoPrv, memoLong)
	require.Error(t, memoLong.BasicCheck())
}

func TestAnchorTxDecode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	trx := signedAnchor(t, ts, 300, 1)
	raw, err := trx.Bytes()
	require.NoError(t, err)
	require.Equal(t, byte(0x01), raw[1])
	typeAt := anchorTypeIndex(1, "")
	require.Equal(t, byte(payload.TypeAnchor), raw[typeAt])

	decoded, err := tx.FromBytes(raw)
	require.NoError(t, err)
	require.NoError(t, decoded.BasicCheck())

	for _, version := range []byte{0, 2} {
		bad := append([]byte{}, raw...)
		bad[1] = version
		got, err := tx.FromBytes(bad)
		require.NoError(t, err)
		require.Error(t, got.BasicCheck())
	}

	for _, typ := range []byte{0, 8, 255} {
		bad := append([]byte{}, raw...)
		bad[typeAt] = typ
		_, err := tx.FromBytes(bad)
		require.Error(t, err)
	}

	pub, _ := ts.RandBLSKeyPair()
	receiver := ts.RandAccAddress()
	transfer := tx.NewTransferTx(1, pub.AccountAddress(), receiver, 1, 1)
	transferRaw, err := transfer.Bytes()
	require.NoError(t, err)
	transferAt := anchorTypeIndex(1, "")
	require.Equal(t, byte(payload.TypeTransfer), transferRaw[transferAt])
	transferRaw[transferAt] = byte(payload.TypeAnchor)
	_, err = tx.FromBytes(transferRaw)
	require.Error(t, err)
}

func TestAnchorTxIDAndMemo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	pub, _ := ts.RandBLSKeyPair()
	from := pub.AccountAddress()
	root := bytes.Repeat([]byte{0x11}, 32)

	base := tx.NewAnchorTx(1, from, payload.AnchorActionSet, root, "abc", 0, 127, 1)
	otherHash := append([]byte{}, root...)
	otherHash[len(otherHash)-1] ^= 0xff
	changedHash := tx.NewAnchorTx(1, from, payload.AnchorActionSet, otherHash, "abc", 0, 127, 1)
	changedDeposit := tx.NewAnchorTx(1, from, payload.AnchorActionSet, root, "abc", 0, 128, 1)
	changedFee := tx.NewAnchorTx(1, from, payload.AnchorActionSet, root, "abc", 0, 127, 2)
	changedMemo := tx.NewAnchorTx(1, from, payload.AnchorActionSet, root, "abc", 0, 127, 1, tx.WithMemo("note"))
	changedLock := tx.NewAnchorTx(2, from, payload.AnchorActionSet, root, "abc", 0, 127, 1)

	require.NotEqual(t, base.ID(), changedHash.ID())
	require.NotEqual(t, base.ID(), changedDeposit.ID())
	require.NotEqual(t, base.ID(), changedFee.ID())
	require.NotEqual(t, base.ID(), changedMemo.ID())
	require.NotEqual(t, base.ID(), changedLock.ID())
	require.Equal(t, base.Payload().Value(), changedMemo.Payload().Value())
	require.Equal(t, base.Payload().Value(), changedLock.Payload().Value())
}

// The Set encoding has no timestamp fields. Bytes placed after the payload
// are read as the signature, never as timestamps.
func TestAnchorPayloadCarriesNoTimestamp(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	pub, prv := ts.RandBLSKeyPair()
	from := pub.AccountAddress()
	setTx := tx.NewAnchorTx(1, from, payload.AnchorActionSet, bytes.Repeat([]byte{0x11}, 32), "abc", 4, 300, 1)

	ts.HelperSignTransaction(prv, setTx)
	signed, err := setTx.Bytes()
	require.NoError(t, err)

	inserted := make([]byte, 0, len(signed)+16)
	cut := len(signed) - bls.SignatureSize - bls.PublicKeySize
	inserted = append(inserted, signed[:cut]...)
	inserted = append(inserted, bytes.Repeat([]byte{0xAB}, 16)...)
	inserted = append(inserted, signed[cut:]...)
	forged, err := tx.FromBytes(inserted)
	require.NoError(t, err)
	require.Equal(t, setTx.Payload(), forged.Payload(), "inserted bytes must not reach the payload")
	require.Error(t, forged.BasicCheck(), "they shift the signature, so verification fails")

	delegated := append([]byte{}, signed...)
	delegated[0] |= 0x04
	decoded, err := tx.FromBytes(delegated)
	require.NoError(t, err)
	got := decoded.Payload().(*payload.AnchorPayload)
	require.Equal(t, from, got.From)
	require.Equal(t, payload.AnchorActionSet, got.Action)
	require.Equal(t, bytes.Repeat([]byte{0x11}, 32), got.RootHash)
	require.Equal(t, "abc", got.ManifestURI)
	require.Equal(t, amount.Amount(300), got.Deposit)

	flipped := append([]byte{}, signed...)
	flippedAt := anchorTypeIndex(1, "")
	require.Equal(t, byte(payload.TypeAnchor), flipped[flippedAt])
	flipped[flippedAt] = byte(payload.TypeTransfer)
	flippedTx, err := tx.FromBytes(flipped)
	if err == nil {
		require.Error(t, flippedTx.BasicCheck())
	}
}

func TestAnchorStreamAndStandaloneBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	pub, prv := ts.RandBLSKeyPair()
	from := pub.AccountAddress()
	anchorTx := tx.NewAnchorTx(1, from, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "", 0, 1, 1)
	ts.HelperSignTransaction(prv, anchorTx)
	anchorRaw, err := anchorTx.Bytes()
	require.NoError(t, err)

	transfer := tx.NewTransferTx(1, from, ts.RandAccAddress(), 1, 1)
	ts.HelperSignTransaction(prv, transfer)
	transferRaw, err := transfer.Bytes()
	require.NoError(t, err)

	stream := append(append([]byte{}, anchorRaw...), transferRaw...)
	reader := bytes.NewReader(stream)
	first := new(tx.Tx)
	require.NoError(t, first.Decode(reader))
	second := new(tx.Tx)
	require.NoError(t, second.Decode(reader))
	require.Equal(t, payload.TypeAnchor, first.Payload().Type())
	require.Equal(t, payload.TypeTransfer, second.Payload().Type())
	require.Zero(t, reader.Len())

	unsigned, err := tx.NewAnchorTx(1, from, payload.AnchorActionDelete, nil, "", 0, 0, 1).Bytes()
	require.NoError(t, err)
	unsignedStream := append(append([]byte{}, unsigned...), transferRaw...)
	unsignedReader := bytes.NewReader(unsignedStream)
	require.NoError(t, new(tx.Tx).Decode(unsignedReader))
	require.NoError(t, new(tx.Tx).Decode(unsignedReader))
	require.Zero(t, unsignedReader.Len())

	buf := bytes.NewBuffer(nil)
	require.NoError(t, anchorTx.Encode(buf, tx.StripPublicKey()))
	stripped, err := tx.FromBytes(buf.Bytes())
	require.NoError(t, err)
	require.Nil(t, stripped.PublicKey())
	require.Error(t, stripped.BasicCheck())

	wrapped, err := anchorTx.MarshalCBOR()
	require.NoError(t, err)
	unwrapped := new(tx.Tx)
	require.NoError(t, unwrapped.UnmarshalCBOR(wrapped))
	require.Equal(t, anchorTx.ID(), unwrapped.ID())
}

// FuzzAnchorTx decodes arbitrary bytes as a transaction. Decoding must never
// panic. A decoded anchor transaction must re-encode to the exact bytes it was
// read from, keep its ID, and pass the payload checks unless its signer cannot
// own an anchor.
func FuzzAnchorTx(f *testing.F) {
	blsPrv, err := bls.PrivateKeyFromString("SECRET1PQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SEZYD4L")
	require.NoError(f, err)
	edPrv, err := ed25519.PrivateKeyFromString("SECRET1RQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SW5D8X2")
	require.NoError(f, err)

	sign := func(prv crypto.PrivateKey, trx *tx.Tx) []byte {
		trx.SetSignature(prv.Sign(trx.SignBytes()))
		trx.SetPublicKey(prv.PublicKey())
		raw, err := trx.Bytes()
		require.NoError(f, err)

		return raw
	}
	blsFrom := blsPrv.PublicKeyNative().AccountAddress()
	edFrom := edPrv.PublicKeyNative().AccountAddress()

	setRaw := sign(blsPrv, tx.NewAnchorTx(7, blsFrom, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "ipfs://manifest", 2, 1e9, 1e7, tx.WithMemo("stamp")))
	delRaw := sign(edPrv, tx.NewAnchorTx(7, edFrom, payload.AnchorActionDelete, nil, "", 0, 0, 1e7))
	unsigned, err := tx.NewAnchorTx(7, blsFrom, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x22}, 64), strings.Repeat("u", 128), 0xFF, 1, 0).Bytes()
	require.NoError(f, err)
	cut := len(setRaw) - bls.SignatureSize - bls.PublicKeySize
	inserted := append(append(bytes.Clone(setRaw[:cut]), bytes.Repeat([]byte{0xAB}, 16)...), setRaw[cut:]...)

	f.Add(setRaw)
	f.Add(delRaw)
	f.Add(unsigned)
	f.Add(setRaw[:len(setRaw)/2])
	f.Add(inserted)
	f.Add(append(bytes.Clone(delRaw), setRaw...))

	f.Fuzz(func(t *testing.T, data []byte) {
		trx, err := tx.FromBytes(data)
		if err != nil || trx.Payload().Type() != payload.TypeAnchor {
			return
		}

		encoded, err := trx.Bytes()
		require.NoError(t, err)
		require.Len(t, encoded, trx.SerializeSize())
		// Decoding ignores unknown flag bits, the delegation flag on a payload
		// that is not delegated, and the stripped-key flag on an unsigned
		// transaction; Encode recomputes the flags. This is true for every
		// transaction type, not only anchors. Every byte after the flags must
		// match.
		const stripedPublicKey, notSigned = 0x01, 0x02
		wantFlags := data[0] & (stripedPublicKey | notSigned)
		if wantFlags&notSigned != 0 {
			wantFlags = notSigned // an unsigned transaction has no public key to strip
		}
		require.Equal(t, wantFlags, encoded[0])
		require.Equal(t, data[1:len(encoded)], encoded[1:], "decoding must read a canonical prefix")

		again, err := tx.FromBytes(encoded)
		require.NoError(t, err)
		require.Equal(t, trx.ID(), again.ID())

		pld := trx.Payload().(*payload.AnchorPayload)
		if err := pld.BasicCheck(); err != nil {
			require.False(t, payload.IsAnchorOwner(pld.From), "a decoded anchor only fails on its signer: %v", err)
		}
	})
}

func bytesString(n int) string {
	return string(bytes.Repeat([]byte{'a'}, n))
}

// flags, version, uint32 lock time, fee varint, memo.
func anchorTypeIndex(fee amount.Amount, memo string) int {
	return 1 + 1 + 4 + encoding.VarIntSerializeSize(uint64(fee)) + encoding.VarStringSerializeSize(memo)
}
