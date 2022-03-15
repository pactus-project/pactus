package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteWithdrawTx(t *testing.T) {
	setup(t)
	exe := NewWithdrawExecutor(true)

	// Let's create a validator first
	addr := crypto.GenerateTestAddress()
	pub, _ := bls.GenerateTestKeyPair()
	val := tSandbox.MakeNewValidator(pub)
	acc := tSandbox.TestStore.RandomTestAcc()
	amt, _ := randomAmountAndFee(acc.Balance())
	val.AddToStake(amt)
	acc.SubtractFromBalance(amt)
	tSandbox.UpdateAccount(acc)
	tSandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, 1, crypto.GenerateTestAddress(), addr, amt, "invalid validator")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+2, val.Address(), addr, amt, "invalid sequence")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt+1, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())

		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, "need to unbond first")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	val.UpdateUnbondingHeight(tSandbox.CurrentHeight() - tSandbox.UnbondInterval() + 1)
	tSandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())

		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, "not passed unbonding interval")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	tSandbox.TestStore.AddTestBlock(500001)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, "should be able to empty stake")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, 1, "can't withdraw empty stake")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Zero(t, exe.Fee())
	assert.Zero(t, tSandbox.Validator(val.Address()).Stake())
	assert.Equal(t, tSandbox.Account(addr).Balance(), amt)
	assert.Equal(t, tSandbox.Validator(val.Address()).Stake(), int64(0))
	assert.Equal(t, tSandbox.Validator(val.Address()).Power(), int64(0))
	assert.Equal(t, tSandbox.Account(addr).Balance(), amt)

	checkTotalCoin(t, 0)
}
