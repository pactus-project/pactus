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

	addr := td.RandAccAddress()
	pub, _ := td.RandBLSKeyPair()
	val := td.sandbox.MakeNewValidator(pub)
	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	amt := td.RandAmountRange(0, acc.Balance())
	fee := td.RandFee()
	val.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	td.sandbox.UpdateAccount(accAddr, acc)
	td.sandbox.UpdateValidator(val)
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, td.RandAccAddress(), addr,
			amt, fee, "invalid validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidAddress, errors.Code(err))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), addr,
			amt+1, fee, "insufficient balance")

		err := exe.Execute(trx, td.sandbox)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(lockTime, val.Address(), addr,
			amt, fee, "need to unbond first")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidHeight, errors.Code(err))
	})

	val.UpdateUnbondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().UnbondInterval + 1)
	td.sandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(lockTime, val.Address(), addr,
			amt, fee, "not passed unbonding interval")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidHeight, errors.Code(err))
	})

	curHeight := td.sandbox.CurrentHeight()
	td.sandbox.TestStore.AddTestBlock(curHeight + 1)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), addr,
			amt, fee, "should be able to empty stake")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), addr,
			1, fee, "can't withdraw empty stake")

		err := exe.Execute(trx, td.sandbox)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	assert.Zero(t, td.sandbox.Validator(val.Address()).Stake())
	assert.Equal(t, amt, td.sandbox.Account(addr).Balance())
	assert.Zero(t, td.sandbox.Validator(val.Address()).Stake())
	assert.Zero(t, td.sandbox.Validator(val.Address()).Power())
	assert.Equal(t, amt, td.sandbox.Account(addr).Balance())

	td.checkTotalCoin(t, fee)
}
