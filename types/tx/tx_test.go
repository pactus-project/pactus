package tx

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCBORMarshaling(t *testing.T) {
	tx1, _ := GenerateTestSendTx()
	bz, err := cbor.Marshal(tx1)
	assert.NoError(t, err)
	tx2 := new(Tx)
	assert.NoError(t, cbor.Unmarshal(bz, tx2))
	assert.Equal(t, tx1.ID(), tx2.ID())

	assert.Error(t, cbor.Unmarshal([]byte{1}, tx2))
}

func TestEncodingTx(t *testing.T) {
	tx1, _ := GenerateTestSendTx()
	len := tx1.SerializeSize()

	for i := 0; i < len; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, tx1.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(len)
	assert.NoError(t, tx1.Encode(w))

	for i := 0; i < len; i++ {
		tx2 := new(Tx)
		r := util.NewFixedReader(i, w.Bytes())
		assert.Error(t, tx2.Decode(r), "decode test %v failed", i)
	}

	tx2 := new(Tx)
	r := util.NewFixedReader(len, w.Bytes())
	assert.NoError(t, tx2.Decode(r))
	assert.Equal(t, tx1.ID(), tx2.ID())
}

func TestTxFromBytes(t *testing.T) {
	trx1, _ := GenerateTestSendTx()
	trx2, _ := GenerateTestBondTx()
	trx3, _ := GenerateTestUnbondTx()
	trx4, _ := GenerateTestWithdrawTx()
	trx5, _ := GenerateTestSortitionTx()
	tests := []*Tx{trx1, trx2, trx3, trx4, trx5}
	assert.True(t, trx1.IsSendTx())
	assert.True(t, trx2.IsBondTx())
	assert.True(t, trx3.IsUnbondTx())
	assert.True(t, trx4.IsWithdrawTx())
	assert.True(t, trx5.IsSortitionTx())

	for _, tx := range tests {
		assert.NoError(t, tx.SanityCheck())

		bz, err := tx.Bytes()
		assert.NoError(t, err)
		tx2, err := FromBytes(bz)
		assert.NoError(t, err)
		assert.Equal(t, tx.Version(), tx2.Version())
		assert.Equal(t, tx.Stamp(), tx2.Stamp())
		assert.Equal(t, tx.Sequence(), tx2.Sequence())
		assert.Equal(t, tx.Payload().Value(), tx2.Payload().Value())
		assert.Equal(t, tx.Payload().Signer(), tx2.Payload().Signer())
		assert.Equal(t, tx.Payload().Type(), tx2.Payload().Type())
		assert.Equal(t, tx.Fee(), tx2.Fee())
		assert.Equal(t, tx.Memo(), tx2.Memo())
		assert.Equal(t, tx.ID(), tx2.ID())
		assert.True(t, tx.PublicKey().EqualsTo(tx2.PublicKey()))
		assert.True(t, tx.Signature().EqualsTo(tx2.Signature()))
	}

	_, err := FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestTxIDNoSignatory(t *testing.T) {
	tx1, _ := GenerateTestSendTx()
	tx2 := new(Tx)
	*tx2 = *tx1
	tx2.data.PublicKey = nil
	tx2.data.Signature = nil
	require.Equal(t, tx1.ID(), tx2.ID())
	require.Equal(t, tx1.SignBytes(), tx2.SignBytes())
}

func TestTxSanityCheck(t *testing.T) {
	t.Run("Invalid sequence", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Sequence = -1
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", maxMemoLength+1)

		trx, _ := GenerateTestSendTx()
		trx.data.Memo = bigMemo
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidMemo)
	})

	t.Run("Invalid payload, Should returns error", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Payload.(*payload.SendPayload).Sender[0] = 0x2
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})
}

func TestInvalidFee(t *testing.T) {
	t.Run("Invalid fee", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewSubsidyTx(stamp, 88, crypto.GenerateTestAddress(), 2500, "subsidy")
		assert.True(t, trx.IsSubsidyTx())
		trx.data.Fee = 1
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid send fee", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Fee = 0
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid sortition fee", func(t *testing.T) {
		trx, _ := GenerateTestSortitionTx()
		trx.data.Fee = 1
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})
}

func TestSubsidyTx(t *testing.T) {
	pub, prv := bls.GenerateTestKeyPair()

	t.Run("Has signature", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewSubsidyTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})
}

func TestInvalidSignature(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Signature = nil
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("No public key", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.PublicKey = nil
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	pbInv, pvInv := bls.GenerateTestKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		sig := pvInv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
	})

	t.Run("Invalid public key", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.SetPublicKey(pbInv)
		err := trx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Memo = "Hello"
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Zero signature", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.Signature = &bls.Signature{}
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Zero public key", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		trx.data.PublicKey = &bls.PublicKey{}
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSendSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSendTx()
		assert.NoError(t, trx.SanityCheck())
		assert.True(t, trx.sanityChecked)
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid amount", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Amount = -1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
	t.Run("Invalid amount", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Amount = 21*1e14 + 1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid sender", func(t *testing.T) {
		trx, signer := GenerateTestSendTx()
		pld := trx.data.Payload.(*payload.SendPayload)
		pld.Sender = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestBondSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestBondTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid version", func(t *testing.T) {
		trx, signer := GenerateTestBondTx()
		trx.data.Version = 2
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid stake", func(t *testing.T) {
		trx, signer := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Stake = -1
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Invalid sender", func(t *testing.T) {
		trx, signer := GenerateTestBondTx()
		pld := trx.data.Payload.(*payload.BondPayload)
		pld.Sender = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSortitionSanityCheck(t *testing.T) {
	invAddr := crypto.GenerateTestAddress()
	t.Run("Ok", func(t *testing.T) {
		trx, _ := GenerateTestSortitionTx()
		assert.NoError(t, trx.SanityCheck())
	})

	t.Run("Invalid address", func(t *testing.T) {
		trx, signer := GenerateTestSortitionTx()
		pld := trx.data.Payload.(*payload.SortitionPayload)
		pld.Address = invAddr
		signer.SignMsg(trx)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestSignBytes(t *testing.T) {
	d, _ := hex.DecodeString("01f10c077fcc04f5ef819fc9d6080101d3e45d249a39d806a1faec2fd85820db340b98e30168fc72a1a961933e694439b2e3c8751d27de5ad3b9c3dc91b9c9b59b010c746573742073656e642d7478b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a68d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e0611776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc1")
	h, _ := hash.FromString("2a04aef409194ff72e942346525428f6c030e2875be27205cb2ce46065ec543f")
	trx, err := FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, trx.SerializeSize(), len(d))

	sb := d[:len(d)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, sb, trx.SignBytes())
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.ID(), hash.CalcHash(sb))
}
