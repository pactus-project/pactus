package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithdrawTx(t *testing.T) {
	td := setup(t)
	exe := NewWithdrawExecutor(true)

	addr := td.RandomAddress()
	pub, _ := td.RandomBLSKeyPair()
	val := td.sandbox.MakeNewValidator(pub)
	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	amt, fee := td.randomAmountAndFee(acc.Balance())
	val.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	td.sandbox.UpdateAccount(accAddr, acc)
	td.sandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.stamp500000, 1, td.RandomAddress(), addr,
			amt, fee, "invalid validator")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+2, val.Address(), addr,
			amt, fee, "invalid sequence")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+1, val.Address(), addr,
			amt+1, fee, "insufficient balance")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+1, val.Address(), addr,
			amt, fee, "need to unbond first")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidHeight)
	})

	val.UpdateUnbondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().UnbondInterval + 1)
	td.sandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+1, val.Address(), addr,
			amt, fee, "not passed unbonding interval")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidHeight)
	})

	td.sandbox.TestStore.AddTestBlock(500001)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+1, val.Address(), addr,
			amt, fee, "should be able to empty stake")

		assert.NoError(t, exe.Execute(trx, td.sandbox))
		assert.Error(t, exe.Execute(trx, td.sandbox), "Execute again, should fail")
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.stamp500000, val.Sequence()+1, val.Address(), addr,
			1, fee, "can't withdraw empty stake")

		assert.Error(t, exe.Execute(trx, td.sandbox))
	})

	assert.Equal(t, exe.Fee(), fee)
	assert.Zero(t, td.sandbox.Validator(val.Address()).Stake())
	assert.Equal(t, td.sandbox.Account(addr).Balance(), amt)
	assert.Equal(t, td.sandbox.Validator(val.Address()).Stake(), int64(0))
	assert.Equal(t, td.sandbox.Validator(val.Address()).Power(), int64(0))
	assert.Equal(t, td.sandbox.Account(addr).Balance(), amt)

	td.checkTotalCoin(t, fee)
}
