package tx_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tx1, _ := ts.GenerateTestTransferTx()
	bz, err := cbor.Marshal(tx1)
	assert.NoError(t, err)
	tx2 := new(tx.Tx)
	assert.NoError(t, cbor.Unmarshal(bz, tx2))
	assert.Equal(t, tx1.ID(), tx2.ID())

	assert.Error(t, cbor.Unmarshal([]byte{1}, tx2))
}

func TestEncodingTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1, _ := ts.GenerateTestTransferTx()
	trx2, _ := ts.GenerateTestBondTx()
	trx3, _ := ts.GenerateTestUnbondTx()
	trx4, _ := ts.GenerateTestWithdrawTx()
	trx5, _ := ts.GenerateTestSortitionTx()
	assert.True(t, trx1.IsTransferTx())
	assert.True(t, trx2.IsBondTx())
	assert.True(t, trx3.IsUnbondTx())
	assert.True(t, trx4.IsWithdrawTx())
	assert.True(t, trx5.IsSortitionTx())

	tests := []*tx.Tx{trx1, trx2, trx3, trx4, trx5}
	for _, trx := range tests {
		assert.NoError(t, trx.BasicCheck())
		assert.NoError(t, trx.BasicCheck()) // double basic check

		length := trx.SerializeSize()
		for i := 0; i < length; i++ {
			w := util.NewFixedWriter(i)
			assert.Error(t, trx.Encode(w), "encode test %v failed", i)
		}
		w := util.NewFixedWriter(length)
		assert.NoError(t, trx.Encode(w))

		for i := 0; i < length; i++ {
			trx2 := new(tx.Tx)
			r := util.NewFixedReader(i, w.Bytes())
			assert.Error(t, trx2.Decode(r), "decode test %v failed", i)
		}

		bz, err := trx.Bytes()
		assert.NoError(t, err)
		trx2, err := tx.FromBytes(bz)
		assert.NoError(t, err)
		assert.Equal(t, trx.ID(), trx2.ID())
	}

	_, err := tx.FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestTxIDNoSignatory(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tx1, _ := ts.GenerateTestTransferTx()
	tx2 := new(tx.Tx)
	*tx2 = *tx1
	tx2.SetPublicKey(nil)
	tx2.SetSignature(nil)
	require.Equal(t, tx1.ID(), tx2.ID())
	require.Equal(t, tx1.SignBytes(), tx2.SignBytes())
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(ts.RandStamp(), -1, ts.RandAddress(), ts.RandProof())
		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid sequence: -1",
		})
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 65)

		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), ts.RandInt64(1e9), ts.RandInt64(1e6), bigMemo)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "memo length exceeded: 65",
		})
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		invAddr := ts.RandAddress()
		invAddr[0] = 2
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), invAddr, 1e9, ts.RandInt64(1e6), "invalid address")

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid payload: invalid address: invalid address type",
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), -1, 1, "invalid amount")

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: -1",
		})
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), 21*1e14+1, 1, "invalid amount")

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid amount: 2100000000000001",
		})
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), 1, -1, "invalid fee")

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid fee: -1",
		})
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), 1, 21*1e14+1, "invalid fee")

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid fee: 2100000000000001",
		})
	})

	t.Run("Invalid signer address", func(t *testing.T) {
		signer := ts.RandSigner()
		trx := tx.NewTransferTx(ts.RandStamp(), ts.RandInt32NonZero(100),
			ts.RandAddress(), ts.RandAddress(), 1, 1, "invalid signer")
		signer.SignMsg(trx)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid address: invalid address",
		})
	})

	t.Run("Invalid version", func(t *testing.T) {
		d := ts.DecodingHex(
			"00" + // Flags
				"02" + // Version
				"a1b2c3d4" + // Stamp
				"01" + // Sequence
				"01" + // Fee
				"00" + // Memo
				"01" + // PayloadType
				"00" + // Sender (treasury)
				"012222222222222222222222222222222222222222" + // Receiver
				"01") // Amount

		trx, err := tx.FromBytes(d)
		assert.NoError(t, err)
		err = trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid version: 2",
		})
	})
}

func TestInvalidPayloadType(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	d := ts.DecodingHex(
		"00" + // Flags
			"01" + // Version
			"a1b2c3d4" + // Stamp
			"01" + // Sequence
			"01" + // Fee
			"00" + // Memo
			"06" + // PayloadType
			"00" + // Sender (treasury)
			"012222222222222222222222222222222222222222" + // Receiver
			"01") // Amount

	_, err := tx.FromBytes(d)
	assert.ErrorIs(t, err, tx.InvalidPayloadTypeError{
		PayloadType: payload.Type(6),
	})
}

func TestSubsidyTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandBLSKeyPair()

	t.Run("Has signature", func(t *testing.T) {
		stamp := ts.RandStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "signature set for subsidy transaction",
		})
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := ts.RandStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "public key set for subsidy transaction",
		})
	})

	t.Run("Strip public key", func(t *testing.T) {
		stamp := ts.RandStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.StripPublicKey()

		err := trx.BasicCheck()
		assert.NoError(t, err)
		assert.False(t, trx.IsPublicKeyStriped())
	})
}

func TestInvalidSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Good", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		assert.NoError(t, trx.BasicCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetSignature(nil)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no signature",
		})
	})

	t.Run("No public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(nil)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no public key",
		})
	})

	pbInv, pvInv := ts.RandBLSKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature: " + trx.Signature().String(),
		})
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(pbInv)

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid address: invalid address",
		})
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		trx0, _ := ts.GenerateTestUnbondTx()
		trx := tx.NewUnbondTx(trx0.Stamp(), trx0.Sequence(), trx0.PublicKey().Address(),
			"invalidate signature")
		trx.SetPublicKey(trx0.PublicKey())
		trx.SetSignature(trx0.Signature())

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature: " + trx.Signature().String(),
		})
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetSignature(&bls.Signature{})

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid signature: " +
				"c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		})
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(&bls.PublicKey{})

		err := trx.BasicCheck()
		assert.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "invalid address: invalid address",
		})
	})
}

func TestSignBytes(t *testing.T) {
	d, _ := hex.DecodeString(
		"00" + // Flags
			"01" + // Version
			"a1b2c3d4" + // Stamp
			"01" + // Sequence
			"01" + // Fee
			"00" + // Memo
			"01" + // PayloadType
			"013333333333333333333333333333333333333333" + // Sender
			"012222222222222222222222222222222222222222" + // Receiver
			"01" + // Amount
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // Signature
			"8d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e061" + // PublicKey
			"1776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1")

	h, _ := hash.FromString("85b0af97280eda00256a0814452d5294128377ad27f62dd8643f403c230245e5")
	trx, err := tx.FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, trx.SerializeSize(), len(d))

	sb := d[1 : len(d)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, sb, trx.SignBytes())
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.ID(), hash.CalcHash(sb))
}

func TestStripPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1, _ := ts.GenerateTestTransferTx()
	id1 := trx1.ID()
	assert.NoError(t, trx1.BasicCheck())

	trx1.StripPublicKey()
	assert.True(t, trx1.IsPublicKeyStriped())
	assert.Equal(t, id1, trx1.ID())
	assert.ErrorIs(t, trx1.BasicCheck(),
		tx.BasicCheckError{
			Reason: "no public key",
		})

	bs1, _ := trx1.Bytes()
	trx2, _ := tx.FromBytes(bs1)
	bs2, _ := trx2.Bytes()

	assert.Equal(t, bs1, bs2)
	assert.Equal(t, trx1.ID(), trx2.ID())
	assert.Nil(t, trx2.PublicKey())
}
