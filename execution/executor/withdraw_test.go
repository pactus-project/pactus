package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteWithdrawTx(t *testing.T) {
	setup(t)
	exe := NewWithdrawExecutor(true)

	// Let's create a validator first
	addr := crypto.GenerateTestAddress()
	pub, _ := bls.GenerateTestKeyPair()
	val := tSandbox.MakeNewValidator(pub)
	acc := tSandbox.TestStore.RandomTestAcc()
	amt, fee := randomAmountAndFee(acc.Balance())
	val.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	tSandbox.UpdateAccount(acc)
	tSandbox.UpdateValidator(val)

	t.Run("Should fail, Invalid validator", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, 1, crypto.GenerateTestAddress(), addr, amt, fee, "invalid validator")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+2, val.Address(), addr, amt, fee, "invalid sequence")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt+1, fee, "insufficient balance")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		assert.Zero(t, val.UnbondingHeight())

		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, fee, "need to unbond first")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidHeight)
	})

	val.UpdateUnbondingHeight(tSandbox.CurrentHeight() - tSandbox.UnbondInterval() + 1)
	tSandbox.UpdateValidator(val)
	t.Run("Should fail, hasn't passed unbonding interval", func(t *testing.T) {
		assert.NotZero(t, val.UnbondingHeight())

		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, fee, "not passed unbonding interval")
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidHeight)
	})

	tSandbox.TestStore.AddTestBlock(500001)

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, amt, fee, "should be able to empty stake")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, can't withdraw empty stake", func(t *testing.T) {
		trx := tx.NewWithdrawTx(tStamp500000, val.Sequence()+1, val.Address(), addr, 1, fee, "can't withdraw empty stake")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, exe.Fee(), fee)
	assert.Zero(t, tSandbox.Validator(val.Address()).Stake())
	assert.Equal(t, tSandbox.Account(addr).Balance(), amt)
	assert.Equal(t, tSandbox.Validator(val.Address()).Stake(), int64(0))
	assert.Equal(t, tSandbox.Validator(val.Address()).Power(), int64(0))
	assert.Equal(t, tSandbox.Account(addr).Balance(), amt)

	checkTotalCoin(t, fee)
}
