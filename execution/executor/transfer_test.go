package executor

import (
	"testing"

	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	sandbox *sandbox.MockSandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	_ = sb.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite: ts,
		sandbox:   sb,
	}
}

func (td *testData) checkTotalCoin(t *testing.T, fee amount.Amount) {
	t.Helper()

	total := amount.Amount(0)
	for _, acc := range td.sandbox.TestStore.Accounts {
		total += acc.Balance()
	}

	for _, val := range td.sandbox.TestStore.Validators {
		total += val.Stake()
	}
	assert.Equal(t, amount.Amount(21_000_000*1e9), total+fee)
}

func TestExecuteTransferTx(t *testing.T) {
	td := setup(t)
	exe := NewTransferExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	receiverAddr := td.RandAccAddress()
	amt := td.RandAmountRange(0, senderBalance)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, td.RandAccAddress(),
			receiverAddr, amt, fee, "non-existing account")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidAddress, errors.Code(err))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr,
			receiverAddr, senderBalance+1, 0, "insufficient balance")

		err := exe.Execute(trx, td.sandbox)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, senderAddr,
			receiverAddr, amt, fee, "ok")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)
	})

	assert.Equal(t, senderBalance-(amt+fee), td.sandbox.Account(senderAddr).Balance())
	assert.Equal(t, amt, td.sandbox.Account(receiverAddr).Balance())

	td.checkTotalCoin(t, fee)
}

func TestTransferToSelf(t *testing.T) {
	td := setup(t)
	exe := NewTransferExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	amt := td.RandAmountRange(0, senderAcc.Balance())
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	trx := tx.NewTransferTx(lockTime, senderAddr, senderAddr, amt, fee, "ok")
	err := exe.Execute(trx, td.sandbox)
	assert.NoError(t, err)

	expectedBalance := senderAcc.Balance() - fee // Fee should be deducted
	assert.Equal(t, expectedBalance, td.sandbox.Account(senderAddr).Balance())
}
