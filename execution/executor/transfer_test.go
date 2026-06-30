package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestExecuteTransferTx(t *testing.T) {
	td := setup(t)

	senderAcc, senderAddr := td.addTestAccount(t)
	senderBalance := senderAcc.Balance()
	receiverAddr := td.RandAccAddress()

	amt, fee := td.randAmountFee(senderBalance)
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randomAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, AccountNotFoundError{Address: randomAddr})
		td.check(t, trx, false, AccountNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr, receiverAddr, senderBalance+1, 0)

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr, receiverAddr, amt, fee)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sbx.Account(senderAddr)
	updatedReceiverAcc := td.sbx.Account(receiverAddr)

	assert.Equal(t, senderBalance-(amt+fee), updatedSenderAcc.Balance())
	assert.Equal(t, amt, updatedReceiverAcc.Balance())

	td.checkTotalCoin(t, fee)
}

func TestTransferToSelf(t *testing.T) {
	td := setup(t)

	senderAcc, senderAddr := td.addTestAccount(t)
	firstBalance := senderAcc.Balance()
	amt, fee := td.randAmountFee(senderAcc.Balance())
	lockTime := td.sbx.CurrentHeight()

	trx := tx.NewTransferTx(lockTime, senderAddr, senderAddr, amt, fee)
	td.check(t, trx, true, nil)
	td.check(t, trx, false, nil)
	td.execute(t, trx)

	secondBalance := senderAcc.Balance()
	assert.Equal(t, firstBalance-fee, secondBalance, "balance should only decrease by fee")

	td.checkTotalCoin(t, fee)
}

func TestTransferSecp256k1(t *testing.T) {
	td := setup(t)

	senderAcc, senderAddr := td.addTestAccount(t)
	amt, fee := td.randAmountFee(senderAcc.Balance())
	lockTime := td.sbx.CurrentHeight()

	trx := tx.NewTransferTx(lockTime, senderAddr, td.RandAccAddressSecp256k1(), amt, fee)

	td.sbx.SbxParams.BlockVersion = protocol.ProtocolVersion3
	td.check(t, trx, true, ErrSecp256k1AccountNotSupported)
	td.check(t, trx, false, ErrSecp256k1AccountNotSupported)

	td.sbx.SbxParams.BlockVersion = protocol.ProtocolVersion4
	td.check(t, trx, true, nil)
	td.check(t, trx, false, nil)
}
