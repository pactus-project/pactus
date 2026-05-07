package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithdrawTx(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	totalStake := val.Stake()
	fee := td.RandFee()
	amt := td.RandAmountRange(0, totalStake-fee)
	senderAddr := val.Address()
	receiverAddr := td.RandAccAddress()
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandValAddress()
		trx := tx.NewWithdrawTx(lockTime, randomAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, ValidatorNotFoundError{Address: randomAddr})
		td.check(t, trx, false, ValidatorNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, senderAddr, receiverAddr, totalStake, 1)

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Should fail, hasn't unbonded yet", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, senderAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, ErrValidatorBonded)
		td.check(t, trx, false, ErrValidatorBonded)
	})

	t.Run("Should fail, hasn't passed unbonding period", func(t *testing.T) {
		val.UpdateUnbondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().UnbondInterval - 1))
		trx := tx.NewWithdrawTx(lockTime, senderAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, ErrUnbondingPeriod)
		td.check(t, trx, false, ErrUnbondingPeriod)
	})

	val.UpdateUnbondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().UnbondInterval))

	t.Run("Should pass, Everything is Ok!", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, senderAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderVal := td.sbx.Validator(senderAddr)
	updatedReceiverAcc := td.sbx.Account(receiverAddr)

	assert.Equal(t, totalStake-amt-fee, updatedSenderVal.Stake())
	assert.Equal(t, amt, updatedReceiverAcc.Balance())

	td.checkTotalCoin(t, fee)
}

func TestExecuteDelegatedWithdrawTx(t *testing.T) {
	td := setup(t)

	valPub, _ := td.RandBLSKeyPair()
	val := td.sbx.MakeNewValidator(valPub)
	totalStake := td.params.MaximumStake
	val.AddToStake(totalStake)
	owner := td.RandAccAddress()
	val.SetDelegation(owner, amount.Amount(0.3e9), td.sbx.CurrentHeight()+10)
	val.UpdateUnbondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().UnbondInterval + 1))
	td.sbx.UpdateValidator(val)

	fee := td.RandFee()
	amt := td.RandAmountRange(0, totalStake-fee)
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, receiver must be stake owner", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), td.RandAccAddress(), amt, fee)

		td.check(t, trx, true, ErrWithdrawMustGoToStakeOwner)
		td.check(t, trx, false, ErrWithdrawMustGoToStakeOwner)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), owner, amt, fee)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedVal := td.sbx.Validator(val.Address())
	updatedOwner := td.sbx.Account(owner)
	assert.Equal(t, totalStake-amt-fee, updatedVal.Stake())
	assert.Equal(t, amt, updatedOwner.Balance())
}
