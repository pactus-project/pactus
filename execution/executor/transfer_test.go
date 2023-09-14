package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	sandbox    *sandbox.MockSandbox
	randStamp  hash.Stamp
	randHeight uint32
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	sandbox := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	randBlock := sandbox.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite:  ts,
		sandbox:    sandbox,
		randStamp:  randBlock.Stamp(),
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
	amt := td.RandInt64NonZero(max)
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
	receiverAddr := td.RandAddress()
	amt, fee := td.randomAmountAndFee(0, senderBalance)

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewTransferTx(td.randStamp, td.sandbox.CurrentHeight()+1, td.RandAddress(),
			receiverAddr, amt, fee, "non-existing account")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewTransferTx(td.randStamp, td.sandbox.CurrentHeight()+2, senderAddr,
			receiverAddr, senderBalance+1, 0, "insufficient balance")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewTransferTx(td.randStamp, td.sandbox.CurrentHeight()+3, senderAddr,
			receiverAddr, amt, fee, "ok")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)

		// Execute again, should fail
		err = exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidLockTime)
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

	trx := tx.NewTransferTx(td.randStamp, td.sandbox.CurrentHeight()+1, senderAddr, senderAddr, amt, fee, "ok")
	err := exe.Execute(trx, td.sandbox)
	assert.NoError(t, err)

	assert.Equal(t, td.sandbox.Account(senderAddr).Balance(), senderBalance-fee) // Fee should be deducted
	assert.Equal(t, exe.Fee(), fee)
}

func TestTransferNonStrictMode(t *testing.T) {
	td := setup(t)
	exe1 := NewTransferExecutor(true)
	exe2 := NewTransferExecutor(false)

	receiver1 := td.RandAddress()

	trx1 := tx.NewSubsidyTx(td.randStamp, td.sandbox.CurrentHeight()+1, receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx1, td.sandbox)), errors.ErrInvalidLockTime)
	assert.NoError(t, exe2.Execute(trx1, td.sandbox))

	trx2 := tx.NewSubsidyTx(td.randStamp, td.sandbox.CurrentHeight()+2, receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx2, td.sandbox)), errors.ErrInvalidLockTime)
	assert.Equal(t, errors.Code(exe2.Execute(trx2, td.sandbox)), errors.ErrInvalidLockTime)
}
