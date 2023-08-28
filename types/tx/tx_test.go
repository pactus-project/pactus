package tx_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
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

	tx1, _ := ts.GenerateTestTransferTx()
	length := tx1.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, tx1.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(length)
	assert.NoError(t, tx1.Encode(w))

	for i := 0; i < length; i++ {
		tx2 := new(tx.Tx)
		r := util.NewFixedReader(i, w.Bytes())
		assert.Error(t, tx2.Decode(r), "decode test %v failed", i)
	}

	tx2 := new(tx.Tx)
	r := util.NewFixedReader(length, w.Bytes())
	assert.NoError(t, tx2.Decode(r))
	assert.Equal(t, tx1.ID(), tx2.ID())
}

func TestFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1, _ := ts.GenerateTestTransferTx()
	trx2, _ := ts.GenerateTestBondTx()
	trx3, _ := ts.GenerateTestUnbondTx()
	trx4, _ := ts.GenerateTestWithdrawTx()
	trx5, _ := ts.GenerateTestSortitionTx()
	tests := []*tx.Tx{trx1, trx2, trx3, trx4, trx5}
	assert.True(t, trx1.IsTransferTx())
	assert.True(t, trx2.IsBondTx())
	assert.True(t, trx3.IsUnbondTx())
	assert.True(t, trx4.IsWithdrawTx())
	assert.True(t, trx5.IsSortitionTx())

	for _, trx := range tests {
		assert.NoError(t, trx.BasicCheck())
		assert.NoError(t, trx.BasicCheck()) // double basic check

		bz, err := trx.Bytes()
		assert.NoError(t, err)
		tx2, err := tx.FromBytes(bz)
		assert.NoError(t, err)
		assert.Equal(t, trx.Version(), tx2.Version())
		assert.Equal(t, trx.Stamp(), tx2.Stamp())
		assert.Equal(t, trx.Sequence(), tx2.Sequence())
		assert.Equal(t, trx.Payload().Value(), tx2.Payload().Value())
		assert.Equal(t, trx.Payload().Signer(), tx2.Payload().Signer())
		assert.Equal(t, trx.Payload().Type(), tx2.Payload().Type())
		assert.Equal(t, trx.Fee(), tx2.Fee())
		assert.Equal(t, trx.Memo(), tx2.Memo())
		assert.Equal(t, trx.ID(), tx2.ID())
		assert.True(t, trx.PublicKey().EqualsTo(tx2.PublicKey()))
		assert.True(t, trx.Signature().EqualsTo(tx2.Signature()))
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
		trx := tx.NewSortitionTx(ts.RandomStamp(), -1, ts.RandomAddress(), ts.RandomProof())
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 65)

		trx := tx.NewSubsidyTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandInt64(1e9), bigMemo)

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidMemo)
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		invAddr := ts.RandomAddress()
		invAddr[0] = 2
		trx := tx.NewSubsidyTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			invAddr, 1e9, "invalid address")

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), -1, 1, "invalid amount")

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 21*1e14+1, 1, "invalid amount")

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	t.Run("Invalid signer address", func(t *testing.T) {
		signer := ts.RandomSigner()
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 1, 1, "invalid signer")
		signer.SignMsg(trx)

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid version", func(t *testing.T) {
		d := ts.DecodingHex(
			"00" + // flags
				"02" + // version
				"a1b2c3d4" + // stamp
				"01" + // sequence
				"01" + // fee
				"01" + // payload type
				"00" + // sender (treasury)
				"012222222222222222222222222222222222222222" + // receiver
				"01" + // amount
				"00") // memo
		trx, err := tx.FromBytes(d)
		assert.NoError(t, err)
		err = trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Invalid payload", func(t *testing.T) {
		d := ts.DecodingHex(
			"00" + // flags
				"01" + // version
				"a1b2c3d4" + // stamp
				"01" + // sequence
				"01" + // fee
				"06" + // payload type
				"00" + // sender (treasury)
				"012222222222222222222222222222222222222222" + // receiver
				"01" + // amount
				"00") // memo

		_, err := tx.FromBytes(d)
		assert.Error(t, err)
	})
}

func TestInvalidFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid subsidy fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			crypto.TreasuryAddress, ts.RandomAddress(), 1e9, 1, "invalid fee")
		assert.True(t, trx.IsSubsidyTx())
		err := trx.BasicCheck()

		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid transfer fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			crypto.TreasuryAddress, ts.RandomAddress(), 1e9, 1, "invalid fee")
		assert.True(t, trx.IsSubsidyTx())
		err := trx.BasicCheck()

		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid sortition fee", func(t *testing.T) {
		pld := &payload.SortitionPayload{
			Address: ts.RandomAddress(),
			Proof:   ts.RandomProof(),
		}
		trx := tx.NewTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			pld, 1, "invalid fee")

		assert.True(t, trx.IsSortitionTx())
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 1, -1, "invalid fee")

		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})
}

func TestSubsidyTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandomBLSKeyPair()

	t.Run("Has signature", func(t *testing.T) {
		stamp := ts.RandomStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := ts.RandomStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
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
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("No public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(nil)
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	pbInv, pvInv := ts.RandomBLSKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(pbInv)
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		trx0, _ := ts.GenerateTestUnbondTx()
		trx := tx.NewUnbondTx(trx0.Stamp(), trx0.Sequence(), trx0.PublicKey().Address(),
			"invalidate signature")
		trx.SetPublicKey(trx0.PublicKey())
		trx.SetSignature(trx0.Signature())
		err := trx.BasicCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
		assert.Error(t, trx.BasicCheck())
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetSignature(&bls.Signature{})
		assert.Error(t, trx.BasicCheck())
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		trx.SetPublicKey(&bls.PublicKey{})
		assert.Error(t, trx.BasicCheck())
	})
}

func TestSignBytes(t *testing.T) {
	d, _ := hex.DecodeString(
		"00" + // flags
			"01" + // version
			"a1b2c3d4" + // stamp
			"01" + // sequence
			"01" + // fee
			"01" + // payload type
			"013333333333333333333333333333333333333333" + // sender
			"012222222222222222222222222222222222222222" + // receiver
			"01" + // amount
			"00" + // memo
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6" + // sig
			"8d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e061" + // pub key
			"1776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1") // pub key

	h, _ := hash.FromString("33ad1b0533269ac4a3c919886065d0dcaf425945167d2e90ad965332445661b4")
	trx, err := tx.FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, trx.SerializeSize(), len(d))

	sb := d[1 : len(d)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, sb, trx.SignBytes())
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.ID(), hash.CalcHash(sb))
}

func TestNoPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1, _ := ts.GenerateTestTransferTx()
	trx1.SetPublicKey(nil)
	bs1, _ := trx1.Bytes()

	trx2, _ := tx.FromBytes(bs1)
	bs2, _ := trx2.Bytes()

	assert.Equal(t, bs1, bs2)
	assert.Equal(t, trx1.ID(), trx2.ID())
	assert.Nil(t, trx2.PublicKey())
}
