package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithdrawTx(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	totalStake := val.Stake()
	amt, fee := td.randAmountFee(totalStake)
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
	totalStake := td.sbx.Params().MaximumStake
	val.AddToStake(totalStake)
	owner := td.RandAccAddress()
	val.SetDelegation(owner, amount.Amount(0.3e9), td.sbx.CurrentHeight()+10)
	val.UpdateUnbondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().UnbondInterval + 1))
	td.sbx.UpdateValidator(val)

	amt, fee := td.randAmountFee(totalStake)
	lockTime := td.sbx.CurrentHeight()

	// makeWithdrawTx builds a withdraw transaction signed by delegateOwner.
	makeWithdrawTx := func(to, delegateOwner crypto.Address) *tx.Tx {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), to, amt, fee)
		trx.Payload().(*payload.WithdrawPayload).DelegateOwner = delegateOwner

		return trx
	}

	t.Run("Should fail, operator cannot withdraw", func(t *testing.T) {
		trx := tx.NewWithdrawTx(lockTime, val.Address(), owner, amt, fee)

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Should fail, operator cannot steal stake as fee", func(t *testing.T) {
		// amount is zero and the whole stake is set as the fee.
		trx := tx.NewWithdrawTx(lockTime, val.Address(), owner, 0, totalStake)

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Should fail, wrong delegate owner", func(t *testing.T) {
		trx := makeWithdrawTx(owner, td.RandAccAddress())

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Ok, owner can withdraw to any address", func(t *testing.T) {
		receiver := td.RandAccAddress()
		trx := makeWithdrawTx(receiver, owner)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)

		updatedVal := td.sbx.Validator(val.Address())
		updatedReceiver := td.sbx.Account(receiver)
		assert.Equal(t, totalStake-amt-fee, updatedVal.Stake())
		assert.Equal(t, amt, updatedReceiver.Balance())
	})
}

func TestExecuteNonDelegatedWithdrawWithDelegateOwner(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	totalStake := val.Stake()
	amt, fee := td.randAmountFee(totalStake)
	val.UpdateUnbondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().UnbondInterval))
	td.sbx.UpdateValidator(val)

	trx := tx.NewWithdrawTx(td.sbx.CurrentHeight(), val.Address(), td.RandAccAddress(), amt, fee)
	trx.Payload().(*payload.WithdrawPayload).DelegateOwner = td.RandAccAddress()

	td.check(t, trx, true, ErrInvalidDelegateOwner)
	td.check(t, trx, false, ErrInvalidDelegateOwner)
}
