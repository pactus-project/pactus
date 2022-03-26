package tx

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx/payload"
)

// TODO: write more tests, and check the error codes

func TestJSONMarshaling(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	_, err := tx.MarshalJSON()
	require.NoError(t, err)
}

func TestEncodingTx(t *testing.T) {
	tx1, _ := GenerateTestSendTx()
	tx2, _ := GenerateTestBondTx()
	tx3, _ := GenerateTestUnbondTx()
	tx4, _ := GenerateTestWithdrawTx()
	tx5, _ := GenerateTestSortitionTx()
	tests := []*Tx{tx1, tx2, tx3, tx4, tx5}

	for _, tx := range tests {
		assert.NoError(t, tx.SanityCheck())

		bz, err := cbor.Marshal(tx)
		fmt.Printf("%x\n", bz)
		assert.NoError(t, err)
		tx2 := new(Tx)
		assert.NoError(t, cbor.Unmarshal(bz, tx2))
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
}

func TestTxIDNoSignatory(t *testing.T) {
	tx1, _ := GenerateTestSendTx()
	tx2 := new(Tx)
	*tx2 = *tx1
	tx2.data.PublicKey = nil
	tx2.data.Signature = nil
	require.Equal(t, tx1.ID(), tx2.ID())
}

func TestTxSanityCheck(t *testing.T) {
	t.Run("Invalid sequence", func(t *testing.T) {
		tx, signer := GenerateTestSendTx()
		tx.data.Sequence = -1
		signer.SignMsg(tx)
		err := tx.SanityCheck()
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})
	t.Run("Transaction ID should be same for signed and unsigned transactions", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		id1 := tx.ID()
		sb1 := tx.SignBytes()
		tx.data.PublicKey = nil
		tx.data.Signature = nil
		id2 := tx.ID()
		sb2 := tx.SignBytes()
		assert.Equal(t, id1, id2)
		assert.Equal(t, sb1, sb2)
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", maxMemoLength+1)

		tx, signer := GenerateTestSendTx()
		tx.data.Memo = bigMemo
		signer.SignMsg(tx)
		assert.Error(t, tx.SanityCheck())
	})
}

func TestInvalidFee(t *testing.T) {
	t.Run("Invalid fee", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewMintbaseTx(stamp, 88, crypto.GenerateTestAddress(), 2500, "subsidy")
		assert.True(t, trx.IsMintbaseTx())
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
		trx := NewMintbaseTx(stamp, 88, pub.Address(), 2500, "subsidy")
		sig := prv.Sign(trx.SignBytes())
		trx.SetSignature(sig)
		assert.Error(t, trx.SanityCheck())
	})

	t.Run("Has public key", func(t *testing.T) {
		stamp := hash.GenerateTestStamp()
		trx := NewMintbaseTx(stamp, 88, pub.Address(), 2500, "subsidy")
		trx.SetPublicKey(pub)
		assert.Error(t, trx.SanityCheck())
	})
}

func TestInvalidSignature(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		assert.NoError(t, tx.SanityCheck())
	})

	t.Run("No signature", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.data.Signature = nil
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("No public key", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.data.PublicKey = nil
		assert.Error(t, tx.SanityCheck())
	})

	pbInv, pvInv := bls.GenerateTestKeyPair()
	t.Run("Invalid signature", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		sig := pvInv.Sign(tx.SignBytes())
		tx.SetSignature(sig)
		assert.Error(t, tx.SanityCheck())
	})

	t.Run("Invalid public key", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.SetPublicKey(pbInv)
		assert.Error(t, tx.SanityCheck())

	})

	t.Run("Invalid sign Bytes", func(t *testing.T) {
		tx, _ := GenerateTestSendTx()
		tx.data.Memo = "Hello"
		assert.Error(t, tx.SanityCheck())
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
	trx, err := TxFromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, trx.SerializeSize(), len(d))

	sb := d[:len(d)-bls.PublicKeySize-bls.SignatureSize]
	assert.Equal(t, sb, trx.SignBytes())
	assert.Equal(t, trx.ID(), h)
	assert.Equal(t, trx.ID(), hash.CalcHash(sb))
}

func TestSendSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	addr := crypto.GenerateTestAddress()

	trx1 := NewSendTx(stamp, 1, signer.Address(), addr, 100, 10, "test send-tx")
	signer.SignMsg(trx1)

	trx2 := NewSendTx(stamp, 1, signer.Address(), addr, 100, 10, "test send-tx")
	trx3 := NewSendTx(stamp, 2, signer.Address(), addr, 100, 10, "test send-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
}

func TestBondSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	pub, _ := bls.GenerateTestKeyPair()

	trx1 := NewBondTx(stamp, 1, signer.Address(), pub, 100, 100, "test bond-tx")
	signer.SignMsg(trx1)

	trx2 := NewBondTx(stamp, 1, signer.Address(), pub, 100, 100, "test bond-tx")
	trx3 := NewBondTx(stamp, 2, signer.Address(), pub, 100, 100, "test bond-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsBondTx())
}

func TestUnbondSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()

	trx1 := NewUnbondTx(stamp, 1, signer.Address(), "test unbond-tx")
	signer.SignMsg(trx1)

	trx2 := NewUnbondTx(stamp, 1, signer.Address(), "test unbond-tx")
	trx3 := NewUnbondTx(stamp, 2, signer.Address(), "test unbond-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsUnbondTx())

}
func TestWithdrawSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	addr := crypto.GenerateTestAddress()

	trx1 := NewWithdrawTx(stamp, 1, signer.Address(), addr, 1000, 1000, "test withdraw-tx")
	signer.SignMsg(trx1)

	trx2 := NewWithdrawTx(stamp, 1, signer.Address(), addr, 1000, 1000, "test withdraw-tx")
	trx3 := NewWithdrawTx(stamp, 2, signer.Address(), addr, 1000, 1000, "test withdraw-tx")

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsWithdrawTx())

}

func TestSortitionSignBytes(t *testing.T) {
	stamp := hash.GenerateTestStamp()
	signer := bls.GenerateTestSigner()
	proof := sortition.GenerateRandomProof()

	trx1 := NewSortitionTx(stamp, 1, signer.Address(), proof)
	signer.SignMsg(trx1)

	trx2 := NewSortitionTx(stamp, 1, signer.Address(), proof)
	trx3 := NewSortitionTx(stamp, 2, signer.Address(), proof)

	assert.Equal(t, trx1.SignBytes(), trx2.SignBytes())
	assert.NotEqual(t, trx1.SignBytes(), trx3.SignBytes())
	assert.True(t, trx1.IsSortitionTx())
}
