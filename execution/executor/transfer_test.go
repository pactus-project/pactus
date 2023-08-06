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

	sandbox     *sandbox.MockSandbox
	stamp500000 hash.Stamp
}

func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)

	sandbox := sandbox.MockingSandbox(ts)
	block500000 := sandbox.TestStore.AddTestBlock(500000)
	stamp500000 := block500000.Stamp()

	return &testData{
		TestSuite:   ts,
		sandbox:     sandbox,
		stamp500000: stamp500000,
	}
}

func (td *testData) checkTotalCoin(t *testing.T, fee int64) {
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
	return fee, amt
}

func TestExecuteTransferTx(t *testing.T) {
	td := setup(t)
	exe := NewTransferExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	receiverAddr := td.RandomAddress()
	amt, fee := td.randomAmountAndFee(td.sandbox.TestParams.MinimumFee, senderBalance)

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewTransferTx(td.stamp500000, 1, td.RandomAddress(),
			receiverAddr, amt, fee, "non-existing account")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewTransferTx(td.stamp500000, senderAcc.Sequence()+1, senderAddr,
			receiverAddr, senderBalance+1, 0, "insufficient balance")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewTransferTx(td.stamp500000, senderAcc.Sequence()+2, senderAddr,
			receiverAddr, amt, fee, "invalid sequence")

		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewTransferTx(td.stamp500000, senderAcc.Sequence()+1, senderAddr,
			receiverAddr, amt, fee, "ok")

		assert.NoError(t, exe.Execute(trx, td.sandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, td.sandbox))
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
	amt, fee := td.randomAmountAndFee(td.sandbox.TestParams.MinimumFee ,senderBalance)

	trx := tx.NewTransferTx(td.stamp500000, senderAcc.Sequence()+1, senderAddr, senderAddr, amt, fee, "ok")
	assert.NoError(t, exe.Execute(trx, td.sandbox))

	assert.Equal(t, td.sandbox.Account(senderAddr).Balance(), senderBalance-fee) // Fee should be deducted
	assert.Equal(t, exe.Fee(), fee)
}

func TestTransferNonStrictMode(t *testing.T) {
	td := setup(t)
	exe1 := NewTransferExecutor(true)
	exe2 := NewTransferExecutor(false)

	receiver1 := td.RandomAddress()

	trx1 := tx.NewSubsidyTx(td.stamp500000, int32(td.sandbox.CurrentHeight()), receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx1, td.sandbox)), errors.ErrInvalidSequence)
	assert.NoError(t, exe2.Execute(trx1, td.sandbox))

	trx2 := tx.NewSubsidyTx(td.stamp500000, int32(td.sandbox.CurrentHeight()+1), receiver1, 1, "")
	assert.Equal(t, errors.Code(exe1.Execute(trx2, td.sandbox)), errors.ErrInvalidSequence)
	assert.Equal(t, errors.Code(exe2.Execute(trx2, td.sandbox)), errors.ErrInvalidSequence)
}
