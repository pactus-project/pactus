package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestExecuteTransferTx(t *testing.T) {
	td := setup(t)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	receiverAddr := td.RandAccAddress()

	amt := td.RandAmountRange(0, senderBalance)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randomAddr,
			receiverAddr, amt, fee, "unknown address")

		td.check(t, trx, true, AccountNotFoundError{Address: randomAddr})
		td.check(t, trx, false, AccountNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr,
			receiverAddr, senderBalance+1, 0, "insufficient balance")

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr,
			receiverAddr, amt, fee, "ok")

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sandbox.Account(senderAddr)
	updatedReceiverAcc := td.sandbox.Account(receiverAddr)

	assert.Equal(t, senderBalance-(amt+fee), updatedSenderAcc.Balance())
	assert.Equal(t, amt, updatedReceiverAcc.Balance())

	td.checkTotalCoin(t, fee)
}

func TestTransferToSelf(t *testing.T) {
	td := setup(t)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	amt := td.RandAmountRange(0, senderAcc.Balance())
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	trx := tx.NewTransferTx(lockTime, senderAddr, senderAddr, amt, fee, "ok")
	td.check(t, trx, true, nil)
	td.check(t, trx, false, nil)
	td.execute(t, trx)

	expectedBalance := senderAcc.Balance() - fee // Fee should be deducted
	updatedAcc := td.sandbox.Account(senderAddr)
	assert.Equal(t, expectedBalance, updatedAcc.Balance())

	td.checkTotalCoin(t, fee)
}
