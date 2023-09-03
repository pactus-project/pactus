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

	addr := td.RandAddress()
	pub, _ := td.RandBLSKeyPair()
	val := td.sandbox.MakeNewValidator(pub)
	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	amt, fee := td.randomAmountAndFee(0, acc.Balance())
	val.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	td.sandbox.UpdateAccount(accAddr, acc)
	td.sandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, 1, td.RandAddress(), addr,
			amt, fee, "invalid validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+2, val.Address(), addr,
			amt, fee, "invalid sequence")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), addr,
			amt, fee, "need to unbond first")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	val.UpdateUnbondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().UnbondInterval + 1)
	td.sandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), addr,
			amt, fee, "not passed unbonding interval")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("Should fail, can't withdraw less than stake amount", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), addr,
			val.Stake()-1, fee, "can't withdraw less than stake amount")
		
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	td.sandbox.TestStore.AddTestBlock(td.randHeight + 1)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), addr,
			amt, fee, "should be able to empty stake")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)
		err = exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence, "Execute again, should fail")
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+2, val.Address(), addr,
			1, fee, "can't withdraw empty stake")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
	})

	assert.Equal(t, exe.Fee(), fee)
	assert.Zero(t, td.sandbox.Validator(val.Address()).Stake())
	assert.Equal(t, td.sandbox.Account(addr).Balance(), amt)
	assert.Equal(t, td.sandbox.Validator(val.Address()).Stake(), int64(0))
	assert.Equal(t, td.sandbox.Validator(val.Address()).Power(), int64(0))
	assert.Equal(t, td.sandbox.Account(addr).Balance(), amt)

	td.checkTotalCoin(t, fee)
}
