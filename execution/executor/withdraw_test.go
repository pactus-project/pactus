package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteWithdrawTx(t *testing.T) {
	setup(t)
	exe := NewWithdrawExecutor(true)

	addr := crypto.GenerateTestAddress()

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), 1, crypto.GenerateTestAddress(), addr, 4999999000, 1000, "invalid validator")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+2, tVal1.Address(), addr, 4999999000, 1000, "invalid sequence")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+1, tVal1.Address(), addr, tVal1.Stake()+1, 0, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, tVal1.UnbondingHeight())

		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+1, tVal1.Address(), addr, 4999999000, 1000, "need to unbond first")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	tVal1.UpdateUnbondingHeight(tSandbox.CurHeight - tSandbox.UnbondInterval() + 1)

	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, tVal1.UnbondingHeight())

		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+1, tVal1.Address(), addr, 4999999000, 1000, "not passed unbonding interval")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	tVal1.UpdateUnbondingHeight(tSandbox.CurHeight - tSandbox.UnbondInterval())

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		total := tVal1.Stake()
		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+1, tVal1.Address(), addr, total-1000, 1000, "should be able to empty stake")

		assert.NoError(t, exe.Execute(trx, tSandbox))
		assert.Equal(t, exe.Fee(), int64(1000))

		assert.Zero(t, tSandbox.Validator(tVal1.Address()).Stake())
		assert.Equal(t, tSandbox.Account(addr).Balance(), total-1000)
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		assert.Zero(t, tVal1.Stake())
		trx := tx.NewWithdrawTx(tHash1000000.Stamp(), tVal1.Sequence()+1, tVal1.Address(), addr, 4999999000, 1000, "can't withdraw empty stake")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Stake(), int64(0))
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Power(), int64(0))

	checkTotalCoin(t, 1000)
}
