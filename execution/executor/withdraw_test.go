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

	pub, _ := td.RandBLSKeyPair()
	val := td.sandbox.MakeNewValidator(pub)
	stake := td.RandInt64(1000000000000)
	val.AddToStake(stake) // MaximumStake
	td.sandbox.UpdateValidator(val)

	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	acc.SubtractFromBalance(acc.Balance())
	td.sandbox.UpdateAccount(accAddr, acc)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, 1, td.RandAddress(), accAddr,
			0, 0, "invalid validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+2, val.Address(), accAddr,
			0, 0, "invalid sequence")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), accAddr,
			0, 0, "need to unbond first")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	val.UpdateUnbondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().UnbondInterval + 1)
	td.sandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), accAddr,
			0, 0, "not passed unbonding interval")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})
	td.sandbox.TestStore.AddTestBlock(td.randHeight + 1)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(td.randStamp, val.Sequence()+1, val.Address(), accAddr,
			0, 0, "should be able to empty stake")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)
	})

	assert.Zero(t, td.sandbox.Validator(val.Address()).Stake())
	assert.Equal(t, td.sandbox.Account(accAddr).Balance(), stake)
	assert.Equal(t, td.sandbox.Validator(val.Address()).Stake(), int64(0))
	assert.Equal(t, td.sandbox.Validator(val.Address()).Power(), int64(0))
}
