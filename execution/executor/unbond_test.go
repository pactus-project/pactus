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

	pub, _ := td.RandomBLSKeyPair()
	valAddr := pub.Address()
	val := td.sandbox.MakeNewValidator(pub)
	td.sandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, td.RandomAddress(), "invalid validator")
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+2, valAddr, "invalid sequence")
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		val := td.sandbox.Committee().Proposer(0)
		trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, val.Address(), "inside committee")
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidTx)
	})

	t.Run("Should fail, Cannot unbond if unbonded already", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)
		val.UpdateUnbondingHeight(td.sandbox.CurrentHeight())
		td.sandbox.UpdateValidator(val)

		trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, pub.Address(), "Ok")
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidHeight)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, valAddr, "Ok")

		assert.NoError(t, exe.Execute(trx, td.sandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, td.sandbox))
	})

	assert.Zero(t, td.sandbox.Validator(valAddr).Stake())
	assert.Zero(t, td.sandbox.Validator(valAddr).Power())
	assert.Equal(t, td.sandbox.Validator(valAddr).UnbondingHeight(), td.sandbox.CurrentHeight())
	assert.Equal(t, td.sandbox.PowerDelta(), -1*val.Stake())
	assert.Zero(t, exe.Fee())

	td.checkTotalCoin(t, 0)
}

// TestUnbondInsideCommittee checks if a validator inside the committee tries to
// unbond the stake.
// In non-strict mode it should be accepted.
func TestUnbondInsideCommittee(t *testing.T) {
	td := setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)

	val := td.sandbox.Committee().Proposer(0)
	trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, val.Address(), "")

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
	pub, _ := td.RandomBLSKeyPair()

	val := td.sandbox.MakeNewValidator(pub)
	val.UpdateLastJoinedHeight(td.sandbox.CurrentHeight())
	td.sandbox.UpdateValidator(val)

	trx := tx.NewUnbondTx(td.stamp500000, val.Sequence()+1, pub.Address(), "Ok")
	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}
