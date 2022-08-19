package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/errors"
)

var (
	tSandbox     *sandbox.MockSandbox
	tStamp500000 hash.Stamp
)

func setup(t *testing.T) {
	tSandbox = sandbox.MockingSandbox()

	block500000 := tSandbox.TestStore.AddTestBlock(500000)
	tStamp500000 = block500000.Stamp()
	checkTotalCoin(t, 0)
}

func checkTotalCoin(t *testing.T, fee int64) {
	total := int64(0)
	for _, acc := range tSandbox.TestStore.Accounts {
		total += acc.Balance()
	}
	for _, val := range tSandbox.TestStore.Validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, int64(21000000*1e8))
}

func randomAmountAndFee(max int64) (int64, int64) {
	amt := util.RandInt64(max / 2)
	fee := int64(float64(amt) * tSandbox.FeeFraction())
	return amt, fee
}

func TestExecuteSendTx(t *testing.T) {
	setup(t)
	exe := NewSendExecutor(true)

	sender := tSandbox.TestStore.RandomTestAcc()
	senderBalance := sender.Balance()
	receiver := crypto.GenerateTestAddress()
	amt, fee := randomAmountAndFee(senderBalance)

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, 1, crypto.GenerateTestAddress(), receiver, amt, fee, "non-existing account")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, sender.Sequence()+1, sender.Address(), receiver, senderBalance+1, 0, "insufficient balance")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, sender.Sequence()+2, sender.Address(), receiver, amt, fee, "invalid sequence")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, sender.Sequence()+1, sender.Address(), receiver, amt, fee, "ok")
		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Account(sender.Address()).Balance(), senderBalance-(amt+fee))
	assert.Equal(t, tSandbox.Account(receiver).Balance(), amt)

	checkTotalCoin(t, fee)
}

func TestSendToSelf(t *testing.T) {
	setup(t)
	exe := NewSendExecutor(true)

	sender := tSandbox.TestStore.RandomTestAcc()
	senderBalance := sender.Balance()
	amt, fee := randomAmountAndFee(senderBalance)

	self := sender.Address()
	trx := tx.NewSendTx(tStamp500000, sender.Sequence()+1, sender.Address(), sender.Address(), amt, fee, "ok")
	assert.NoError(t, exe.Execute(trx, tSandbox))

	assert.Equal(t, tSandbox.Account(self).Balance(), senderBalance-fee) // Fee should be deducted
	assert.Equal(t, exe.Fee(), fee)
}

func TestSendNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSendExecutor(true)
	exe2 := NewSendExecutor(false)

	receiver1 := crypto.GenerateTestAddress()

	trx1 := tx.NewSubsidyTx(tStamp500000, int32(tSandbox.CurrentHeight()), receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx1, tSandbox)), errors.ErrInvalidSequence)
	assert.NoError(t, exe2.Execute(trx1, tSandbox))

	trx2 := tx.NewSubsidyTx(tStamp500000, int32(tSandbox.CurrentHeight()+1), receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx2, tSandbox)), errors.ErrInvalidSequence)
	assert.Equal(t, errors.Code(exe2.Execute(trx2, tSandbox)), errors.ErrInvalidSequence)
}
