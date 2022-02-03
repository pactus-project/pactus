package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
)

func Test_WithdrawExecutor(t *testing.T) {
	setup(t)
	exe := NewWithdrawExecutor(true)

	addr := crypto.GenerateTestAddress()
	hash100 := hash.GenerateTestHash()

	tSandbox.AppendNewBlock(100, hash100)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(hash100.Stamp(), 1, addr, tAcc1.Address(), 1000, 1000, "invalid validator")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(hash100.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+2, tVal1.Address(), tAcc1.Address(), 1000, 1000, "invalid sequence")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		assert.Equal(t, int64(5000000000), tVal1.Stake())
		trx := tx.NewWithdrawTx(hash100.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+1, tVal1.Address(), tAcc1.Address(), tVal1.Stake()+1, 0, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Equal(t, 0, tVal1.UnbondingHeight())
		trx := tx.NewWithdrawTx(hash100.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+1, tVal1.Address(), tAcc1.Address(), 1000, 1000, "need to unbond first")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.Equal(t, 0, tVal1.UnbondingHeight())
		tVal1.UpdateUnbondingHeight(101)
		hash201 := hash.GenerateTestHash()
		tSandbox.AppendNewBlock(201, hash201)
		assert.Equal(t, 101, tVal1.UnbondingHeight())

		trx := tx.NewWithdrawTx(hash201.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+1, tVal1.Address(), tAcc1.Address(), 1000, 1000, "not passed unbonding interval")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		hash := hash.GenerateTestHash()
		tSandbox.AppendNewBlock(tVal1.UnbondingHeight()+tSandbox.UnbondInterval(), hash)
		assert.Equal(t, 101, tVal1.UnbondingHeight())

		trx := tx.NewWithdrawTx(hash.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+1, tVal1.Address(), tAcc1.Address(), 4999999000, 1000, "should be able to empty stake")

		assert.NoError(t, exe.Execute(trx, tSandbox))
		assert.Zero(t, tVal1.Stake())
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		assert.Zero(t, tVal1.Stake())
		trx := tx.NewWithdrawTx(hash100.Stamp(), tSandbox.Validator(tVal1.Address()).Sequence()+1, tVal1.Address(), tAcc1.Address(), 4999999000, 1000, "can't withdraw empty stake")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Stake(), int64(0))
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Power(), int64(0)) //it shouldn't return 1 but it does
	assert.Equal(t, tSandbox.Account(tAcc1.Address()).Balance(), int64(14999999000))
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).UnbondingHeight(), 101)
	assert.Equal(t, exe.Fee(), int64(1000))

	checkTotalCoin(t, 1000)
}
