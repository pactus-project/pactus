package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
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
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+2, pub.Address(), "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		val := tSandbox.Committee().Proposer(0)
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, val.Address(), "inside committee")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Cannot unbond if unbonded already", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})
	assert.Zero(t, tSandbox.Validator(pub.Address()).Stake())
	assert.Zero(t, tSandbox.Validator(pub.Address()).Power())
	assert.Equal(t, tSandbox.Validator(pub.Address()).UnbondingHeight(), tSandbox.CurrentHeight())
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

func TestUnbondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)

	val := tSandbox.Committee().Proposer(0)
	trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, val.Address(), "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

func TestUnbondJoiningCommittee(t *testing.T) {
	setup(t)
	exe := NewUnbondExecutor(true)
	pub, _ := bls.GenerateTestKeyPair()

	val := tSandbox.MakeNewValidator(pub)
	val.UpdateLastJoinedHeight(tSandbox.CurrentHeight())
	tSandbox.UpdateValidator(val)

	trx := tx.NewUnbondTx(tStamp500000, val.Sequence()+1, pub.Address(), "Ok")

	assert.Error(t, exe.Execute(trx, tSandbox))
}
