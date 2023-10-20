package executor

import (
	"testing"

	"github.com/pactus-project/pactus/sandbox"
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

	sandbox := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	_ = sandbox.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite:  ts,
		sandbox:    sandbox,
		randHeight: randHeight,
	}
}

func (td *testData) checkTotalCoin(t *testing.T, fee int64) {
	t.Helper()

	total := int64(0)
	for _, acc := range td.sandbox.TestStore.Accounts {
		total += acc.Balance()
	}

	for _, val := range td.sandbox.TestStore.Validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, int64(21000000*1e9))
}

func (td *testData) randomAmountAndFee(min int64, max int64) (int64, int64) {
	amt := td.RandInt64NonZero(max - 1) // To make sure amt+fee is less than max
	for amt < min {
		amt = td.RandInt64NonZero(max)
	}

	fee := int64(float64(amt) * td.sandbox.Params().FeeFraction)
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
