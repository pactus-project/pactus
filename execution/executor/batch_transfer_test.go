package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/assert"
)

func TestExecuteBatchTransferTx(t *testing.T) {
	td := setup(t)

	senderAddr, senderAcc := td.sbx.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	receiverAddr1 := td.RandAccAddress()
	receiverAddr2 := td.RandAccAddress()

	amt := td.RandAmountRange(0, senderBalance)
	amt1 := td.RandAmount(amt / 2)
	amt2 := td.RandAmount(amt / 2)
	fee := td.RandFee()
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		recipients := []payload.BatchRecipient{
			{To: receiverAddr1, Amount: amt1},
			{To: receiverAddr2, Amount: amt2},
		}
		trx := tx.NewBatchTransferTx(lockTime, randomAddr, recipients, fee)

		td.check(t, trx, true, AccountNotFoundError{Address: randomAddr})
		td.check(t, trx, false, AccountNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{To: receiverAddr1, Amount: senderBalance + 1},
			{To: receiverAddr2, Amount: senderBalance + 1},
		}
		trx := tx.NewBatchTransferTx(lockTime, senderAddr, recipients, 0)

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{To: receiverAddr1, Amount: amt1},
			{To: receiverAddr2, Amount: amt2},
		}
		trx := tx.NewBatchTransferTx(lockTime, senderAddr, recipients, fee)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sbx.Account(senderAddr)
	updatedReceiverAcc1 := td.sbx.Account(receiverAddr1)
	updatedReceiverAcc2 := td.sbx.Account(receiverAddr2)

	assert.Equal(t, senderBalance-(amt1+amt2+fee), updatedSenderAcc.Balance())
	assert.Equal(t, amt1, updatedReceiverAcc1.Balance())
	assert.Equal(t, amt2, updatedReceiverAcc2.Balance())

	td.checkTotalCoin(t, fee)
}

func TestBatchTransferToSelf(t *testing.T) {
	td := setup(t)

	senderAddr, senderAcc := td.sbx.TestStore.RandomTestAcc()
	amt := td.RandAmountRange(0, senderAcc.Balance())
	amt1 := td.RandAmount(amt / 2)
	amt2 := td.RandAmount(amt / 2)
	fee := td.RandFee()
	lockTime := td.sbx.CurrentHeight()

	recipients := []payload.BatchRecipient{
		{To: td.RandAccAddress(), Amount: amt1},
		{To: senderAddr, Amount: amt2},
	}
	trx := tx.NewBatchTransferTx(lockTime, senderAddr, recipients, fee)
	td.check(t, trx, true, nil)
	td.check(t, trx, false, nil)
	td.execute(t, trx)

	expectedBalance := senderAcc.Balance() - amt1 - fee // Fee should be deducted
	updatedAcc := td.sbx.Account(senderAddr)
	assert.Equal(t, expectedBalance, updatedAcc.Balance())

	td.checkTotalCoin(t, fee)
}
