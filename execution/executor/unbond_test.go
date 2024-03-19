package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteUnbondTx(t *testing.T) {
	td := setup(t)
	exe := NewUnbondExecutor(true)

	bonderAddr, bonderAcc := td.sandbox.TestStore.RandomTestAcc()
	bonderBalance := bonderAcc.Balance()
	stake, _ := td.randomAmountAndFee(td.sandbox.TestParams.MinimumStake, bonderBalance)
	bonderAcc.SubtractFromBalance(stake)
	td.sandbox.UpdateAccount(bonderAddr, bonderAcc)

	pub, _ := td.RandBLSKeyPair()
	valAddr := pub.ValidatorAddress()
	val := td.sandbox.MakeNewValidator(pub)
	val.AddToStake(stake)
	td.sandbox.UpdateValidator(val)
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, td.RandAccAddress(), "invalid validator")
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		val0 := td.sandbox.Committee().Proposer(0)
		trx := tx.NewUnbondTx(lockTime, val0.Address(), "inside committee")
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Should fail, Cannot unbond if unbonded already", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		unbondedVal := td.sandbox.MakeNewValidator(pub)
		unbondedVal.UpdateUnbondingHeight(td.sandbox.CurrentHeight())
		td.sandbox.UpdateValidator(unbondedVal)

		trx := tx.NewUnbondTx(lockTime, pub.ValidatorAddress(), "Ok")
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr, "Ok")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)

		// Execute again, should fail
		err = exe.Execute(trx, td.sandbox)
		assert.Error(t, err)
	})

	assert.Equal(t, stake, td.sandbox.Validator(valAddr).Stake())
	assert.Zero(t, td.sandbox.Validator(valAddr).Power())
	assert.Equal(t, td.sandbox.Validator(valAddr).UnbondingHeight(), td.sandbox.CurrentHeight())
	assert.Equal(t, int64(-val.Stake()), td.sandbox.PowerDelta())

	td.checkTotalCoin(t, 0)
}

// TestUnbondInsideCommittee checks if a validator inside the committee tries to
// unbond the stake.
// In non-strict mode it should be accepted.
func TestUnbondInsideCommittee(t *testing.T) {
	td := setup(t)

	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)
	lockTime := td.sandbox.CurrentHeight()

	val := td.sandbox.Committee().Proposer(0)
	trx := tx.NewUnbondTx(lockTime, val.Address(), "")

	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

// TestUnbondJoiningCommittee checks if a validator tries to unbond after
// evaluating sortition.
// In non-strict mode it should be accepted.
func TestUnbondJoiningCommittee(t *testing.T) {
	td := setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)
	pub, _ := td.RandBLSKeyPair()

	val := td.sandbox.MakeNewValidator(pub)
	val.UpdateLastSortitionHeight(td.randHeight)
	td.sandbox.UpdateValidator(val)
	td.sandbox.JoinedToCommittee(val.Address())
	lockTime := td.sandbox.CurrentHeight()

	trx := tx.NewUnbondTx(lockTime, pub.ValidatorAddress(), "Ok")
	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

func TestPwerDeltaUnbond(t *testing.T) {
	td := setup(t)
	exe := NewUnbondExecutor(true)

	pub, _ := td.RandBLSKeyPair()
	valAddr := pub.ValidatorAddress()
	val := td.sandbox.MakeNewValidator(pub)
	amt := td.RandAmount()
	val.AddToStake(amt)
	td.sandbox.UpdateValidator(val)
	lockTime := td.sandbox.CurrentHeight()
	trx := tx.NewUnbondTx(lockTime, valAddr, "Ok")

	err := exe.Execute(trx, td.sandbox)
	assert.NoError(t, err)

	assert.Equal(t, int64(-amt), td.sandbox.PowerDelta())
}
