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

	sandbox    *sandbox.MockSandbox
	randHeight uint32
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	_ = sb.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite:  ts,
		sandbox:    sb,
		randHeight: randHeight,
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
	assert.Equal(t, total+fee, amount.Amount(21_000_000*1e9))
}

func (td *testData) randomAmountAndFee(min, max amount.Amount) (amount.Amount, amount.Amount) {
	amt := amount.Amount(td.RandInt64NonZero(int64(max)))
	for amt < min {
		amt = amount.Amount(td.RandInt64NonZero(int64(max)))
	}

	fee := amt.MulF64(td.sandbox.Params().FeeFraction)
	if amt+fee > max {
		// To make sure amt+fee is less than max
		return td.randomAmountAndFee(min, max)
	}

	return amt, fee
}

func TestExecuteTransferTx(t *testing.T) {
	td := setup(t)
	exe := NewTransferExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	receiverAddr := td.RandAccAddress()
	amt, fee := td.randomAmountAndFee(0, senderBalance)
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, td.RandAccAddress(),
			receiverAddr, amt, fee, "non-existing account")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
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

	assert.Equal(t, td.sandbox.Account(senderAddr).Balance(), senderBalance-(amt+fee))
	assert.Equal(t, td.sandbox.Account(receiverAddr).Balance(), amt)

	td.checkTotalCoin(t, fee)
}

func TestTransferToSelf(t *testing.T) {
	td := setup(t)
	exe := NewTransferExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	amt, fee := td.randomAmountAndFee(0, senderBalance)
	lockTime := td.sandbox.CurrentHeight()

	trx := tx.NewTransferTx(lockTime, senderAddr, senderAddr, amt, fee, "ok")
	err := exe.Execute(trx, td.sandbox)
	assert.NoError(t, err)

	assert.Equal(t, td.sandbox.Account(senderAddr).Balance(), senderBalance-fee) // Fee should be deducted
}
