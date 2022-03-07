package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteUnbondTx(t *testing.T) {
	setup(t)
	exe := NewUnbondExecutor(true)

	addr := crypto.GenerateTestAddress()

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, 1, addr, "invalid validator")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewUnbondTx(tStamp500000, tSandbox.TestValSeq(tVal1.Address())+2, tVal1.Address(), "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		tSandbox.InCommittee = true
		trx := tx.NewUnbondTx(tStamp500000, tSandbox.TestValSeq(tVal1.Address())+1, tVal1.Address(), "inside committee")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		tSandbox.InCommittee = false
		trx := tx.NewUnbondTx(tStamp500000, tSandbox.TestValSeq(tVal1.Address())+1, tVal1.Address(), "Ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Replay
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Cannot unbond if unbonded already", func(t *testing.T) {
		tSandbox.InCommittee = false
		trx := tx.NewUnbondTx(tStamp500000, tSandbox.TestValSeq(tVal1.Address())+1, tVal1.Address(), "Ok")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Stake(), tVal1Stake)
	assert.Zero(t, tSandbox.Validator(tVal1.Address()).Power())
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).UnbondingHeight(), tSandbox.CurHeight)
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

func TestUnbondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewUnbondExecutor(true)
	exe2 := NewUnbondExecutor(false)

	tSandbox.InCommittee = true

	trx := tx.NewUnbondTx(tStamp500000, tSandbox.TestValSeq(tVal1.Address())+1, tVal1.Address(), "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}
