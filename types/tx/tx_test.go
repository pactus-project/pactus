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

	tx1, _ := ts.GenerateTestSendTx()
	bz, err := cbor.Marshal(tx1)
	assert.NoError(t, err)
	tx2 := new(tx.Tx)
	assert.NoError(t, cbor.Unmarshal(bz, tx2))
	assert.Equal(t, tx1.ID(), tx2.ID())

	assert.Error(t, cbor.Unmarshal([]byte{1}, tx2))
}

func TestEncodingTx(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tx1, _ := ts.GenerateTestSendTx()
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

	trx1, _ := ts.GenerateTestSendTx()
	trx2, _ := ts.GenerateTestBondTx()
	trx3, _ := ts.GenerateTestUnbondTx()
	trx4, _ := ts.GenerateTestWithdrawTx()
	trx5, _ := ts.GenerateTestSortitionTx()
	tests := []*tx.Tx{trx1, trx2, trx3, trx4, trx5}
	assert.True(t, trx1.IsSendTx())
	assert.True(t, trx2.IsBondTx())
	assert.True(t, trx3.IsUnbondTx())
	assert.True(t, trx4.IsWithdrawTx())
	assert.True(t, trx5.IsSortitionTx())

	for _, trx := range tests {
		assert.NoError(t, trx.SanityCheck())
		assert.NoError(t, trx.SanityCheck()) // double sanity check

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

	tx1, _ := ts.GenerateTestSendTx()
	tx2 := new(tx.Tx)
	*tx2 = *tx1
	tx2.SetPublicKey(nil)
	tx2.SetSignature(nil)
	require.Equal(t, tx1.ID(), tx2.ID())
	require.Equal(t, tx1.SignBytes(), tx2.SignBytes())
}

func TestSanityCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(ts.RandomStamp(), -1, ts.RandomAddress(), ts.RandomProof())
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 65)

		trx := tx.NewSubsidyTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandInt64(1e9), bigMemo)

		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidMemo)
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		invAddr := ts.RandomAddress()
		invAddr[0] = 2
		trx := tx.NewSubsidyTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			invAddr, 1e9, "invalid address")

		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), -1, 1, "invalid amount")

		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 21*1e14+1, 1, "invalid amount")

		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	t.Run("Invalid signer address", func(t *testing.T) {
		signer := ts.RandomSigner()
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 1, 1, "invalid signer")
		signer.SignMsg(trx)

		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid version", func(t *testing.T) {
		d, _ := hex.DecodeString(
			"023513630b1a00010001703db2cca1f0deb29fb42b98bd9d12971b1160168094ebdc0300PASS")
		trx, err := tx.FromBytes(d)
		assert.NoError(t, err)
		err = trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})
}

func TestInvalidFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid subsidy fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			crypto.TreasuryAddress, ts.RandomAddress(), 1e9, 1, "invalid fee")
		assert.True(t, trx.IsSubsidyTx())
		err := trx.SanityCheck()

		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid transfer fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			crypto.TreasuryAddress, ts.RandomAddress(), 1e9, 1, "invalid fee")
		assert.True(t, trx.IsSubsidyTx())
		err := trx.SanityCheck()

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
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid fee", func(t *testing.T) {
		trx := tx.NewTransferTx(ts.RandomStamp(), ts.RandInt32NonZero(100),
			ts.RandomAddress(), ts.RandomAddress(), 1, -1, "invalid fee")

		err := trx.SanityCheck()
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
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := ts.RandomStamp()
		trx := tx.NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})
}

func TestInvalidSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Good", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		trx.SetSignature(nil)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("No public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		trx.SetPublicKey(nil)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	pbInv, pvInv := ts.RandomBLSKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		trx.SetPublicKey(pbInv)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		trx0, _ := ts.GenerateTestUnbondTx()
		trx := tx.NewUnbondTx(trx0.Stamp(), trx0.Sequence(), trx0.PublicKey().Address(),
			"invalidate signature")
		trx.SetPublicKey(trx0.PublicKey())
		trx.SetSignature(trx0.Signature())
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		trx.SetSignature(&bls.Signature{})
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		trx.SetPublicKey(&bls.PublicKey{})
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSignBytes(t *testing.T) {
	d, _ := hex.DecodeString(
		"01f10c077fcc04f5ef819fc9d6080101d3e45d249a39d806a1faec2fd85820db340b98e30168fc72a1a961933e694439b2e3c8751d27de5a" +
			"d3b9c3dc91b9c9b59b010c746573742073656e642d7478b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112" +
			"ac01c80d0d9d64c2320664c77fa2a68d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba159137077884" +
			"6681b7b558e0611776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1")
	h, _ := hash.FromString("2a04aef409194ff72e942346525428f6c030e2875be27205cb2ce46065ec543f")
	trx, err := tx.FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, trx.SerializeSize(), len(d))

	sb := d[:len(d)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, sb, trx.SignBytes())
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.ID(), hash.CalcHash(sb))
}
