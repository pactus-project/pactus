package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteUnbondTx(t *testing.T) {
	setup(t)
	exe := NewUnbondExecutor(true)

	// Let's create a validator first
	pub, _ := bls.GenerateTestKeyPair()
	val := tSandbox.MakeNewValidator(pub)
	tSandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, crypto.GenerateTestAddress(), "invalid validator")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+2, pub.Address(), "invalid sequence")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		val := tSandbox.Committee().Proposer(0)
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, val.Address(), "inside committee")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidTx)
	})

	t.Run("Should fail, Cannot unbond if unbonded already", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		val.UpdateUnbondingHeight(tSandbox.CurrentHeight())
		tSandbox.UpdateValidator(val)

		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidHeight)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Zero(t, tSandbox.Validator(pub.Address()).Stake())
	assert.Zero(t, tSandbox.Validator(pub.Address()).Power())
	assert.Equal(t, tSandbox.Validator(pub.Address()).UnbondingHeight(), tSandbox.CurrentHeight())
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

// TestUnbondInsideCommittee checks if a validator inside the committee tries to
// unbond the stake.
// In non-strict mode it should be accepted.
func TestUnbondInsideCommittee(t *testing.T) {
	setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)

	val := tSandbox.Committee().Proposer(0)
	trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, val.Address(), "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

// TestUnbondJoiningCommittee checks if a validator tries to unbond after
// evaluating sortition.
// In non-strict mode it should be accepted.
func TestUnbondJoiningCommittee(t *testing.T) {
	setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)
	pub, _ := bls.GenerateTestKeyPair()

	val := tSandbox.MakeNewValidator(pub)
	val.UpdateLastJoinedHeight(tSandbox.CurrentHeight())
	tSandbox.UpdateValidator(val)

	trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")
	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}
