package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

var (
	tSandbox     *sandbox.MockSandbox
	tStamp500000 hash.Stamp
)

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	tSandbox = sandbox.MockingSandbox()

	block500000 := block.GenerateTestBlock(nil, nil)
	tSandbox.AddTestBlock(500000, block500000)

	tStamp500000 = block500000.Stamp()
	assert.Equal(t, tSandbox.CurHeight, 500001)
}

func checkTotalCoin(t *testing.T, fee int64) {
	total := int64(0)
	for _, acc := range tSandbox.TestAccounts {
		total += acc.Balance()
	}
	for _, val := range tSandbox.TestValidators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, int64(21000000*1e8))
}

func randomAmountandFee(max int64) (int64, int64) {
	amt := util.RandInt64(max / 2)
	fee := int64(float64(amt) * tSandbox.FeeFraction())
	return amt, fee
}

func TestExecuteSendTx(t *testing.T) {
	setup(t)
	exe := NewSendExecutor(true)

	sender := tSandbox.RandomTestAcc()
	senderBalance := sender.Balance()
	receiver := crypto.GenerateTestAddress()
	amt, fee := randomAmountandFee(senderBalance)

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, 1, crypto.GenerateTestAddress(), receiver, amt, fee, "non-existing account")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, sender.Sequence()+1, sender.Address(), receiver, senderBalance+1, 0, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSendTx(tStamp500000, sender.Sequence()+2, sender.Address(), receiver, amt, fee, "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
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

	sender := tSandbox.RandomTestAcc()
	senderBalance := sender.Balance()
	amt, fee := randomAmountandFee(senderBalance)

	self := sender.Address()
	trx := tx.NewSendTx(tStamp500000, sender.Sequence()+1, sender.Address(), sender.Address(), amt, fee, "ok")
	assert.NoError(t, exe.Execute(trx, tSandbox))

	assert.Equal(t, tSandbox.Account(self).Balance(), senderBalance-fee) /// Fee should be deducted
	assert.Equal(t, exe.Fee(), fee)
}

func TestSendNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSendExecutor(true)
	exe2 := NewSendExecutor(false)

	receiver1 := crypto.GenerateTestAddress()

	trx1 := tx.NewMintbaseTx(tStamp500000, tSandbox.CurHeight, receiver1, 1, "")
	assert.Error(t, exe1.Execute(trx1, tSandbox)) // Invalid sequence
	assert.NoError(t, exe2.Execute(trx1, tSandbox))

	trx2 := tx.NewMintbaseTx(tStamp500000, tSandbox.CurHeight+1, receiver1, 1, "")
	assert.Error(t, exe1.Execute(trx2, tSandbox)) // Invalid height
	assert.Error(t, exe2.Execute(trx2, tSandbox)) // Invalid height

}
