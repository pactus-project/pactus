package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	sb.TestAcceptSortition = true
	exe := NewExecutor()
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(100 * 1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	rndValAddr := rndPubKey.ValidatorAddress()
	rndVal := sb.MakeNewValidator(rndPubKey)
	rndVal.AddToStake(100 * 1e9)
	sb.UpdateValidator(rndVal)
	_ = sb.TestStore.AddTestBlock(8642)

	t.Run("Future LockTime, Should return error (+1)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() + 1
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000, "future-lockTime")
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, FutureLockTimeError{LockTime: lockTime})
	})

	t.Run("Past LockTime, Should return error (-8641)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000, "past-lockTime")
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, PastLockTimeError{LockTime: lockTime})
	})

	t.Run("Transaction has valid LockTime (-8640)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000, "ok")
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})

	t.Run("Transaction has valid LockTime (0)", func(t *testing.T) {
		lockTime := sb.CurrentHeight()
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000, "ok")
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})

	t.Run("Subsidy transaction has invalid LockTime (+1)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() + 1
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000,
			"invalid-lockTime")
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, FutureLockTimeError{LockTime: lockTime})
	})

	t.Run("Subsidy transaction has invalid LockTime (-1)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() - 1
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000,
			"invalid-lockTime")
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, PastLockTimeError{LockTime: lockTime})
	})

	t.Run("Subsidy transaction has valid LockTime (0)", func(t *testing.T) {
		lockTime := sb.CurrentHeight()
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000, "ok")
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})

	t.Run("Sortition transaction has invalid LockTime (+1)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() + 1
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(lockTime, rndValAddr, proof)
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, FutureLockTimeError{LockTime: lockTime})
	})

	t.Run("Sortition transaction has invalid LockTime (-8)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(lockTime, rndValAddr, proof)
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.ErrorIs(t, err, PastLockTimeError{LockTime: lockTime})
	})

	t.Run("Sortition transaction has valid LockTime (-7)", func(t *testing.T) {
		lockTime := sb.CurrentHeight() - sb.TestParams.SortitionInterval
		proof := ts.RandProof()

		trx := tx.NewSortitionTx(lockTime, rndValAddr, proof)
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})
}

func TestExecution(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	exe := NewExecutor()

	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndValAddr := rndPubKey.ValidatorAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(100 * 1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	_ = sb.TestStore.AddTestBlock(8642)
	lockTime := sb.CurrentHeight()

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, ts.RandAccAddress(), ts.RandAccAddress(), 1000, 1000, "invalid-tx")
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1, "invalid fee")
		ts.HelperSignTransaction(rndPrvKey, trx)

		expectedErr := InvalidFeeError{Fee: 1, Expected: sb.TestParams.MinimumFee}
		assert.ErrorIs(t, exe.Execute(trx, sb), expectedErr)
		assert.ErrorIs(t, exe.checkFee(trx, sb), expectedErr)
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1002, "invalid fee")
		ts.HelperSignTransaction(rndPrvKey, trx)

		expectedErr := InvalidFeeError{Fee: 1002, Expected: sb.TestParams.MinimumFee}
		assert.ErrorIs(t, exe.Execute(trx, sb), expectedErr)
		assert.ErrorIs(t, exe.checkFee(trx, sb), expectedErr)
	})

	t.Run("Invalid fee (subsidy tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, crypto.TreasuryAddress, ts.RandAccAddress(), 1000, 1, "invalid fee")

		expectedErr := InvalidFeeError{Fee: 1, Expected: 0}
		assert.ErrorIs(t, exe.Execute(trx, sb), expectedErr)
		assert.ErrorIs(t, exe.checkFee(trx, sb), expectedErr)
	})

	t.Run("Invalid fee (transfer tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 0, "invalid fee")

		expectedErr := InvalidFeeError{Fee: 0, Expected: sb.TestParams.MinimumFee}
		assert.ErrorIs(t, exe.Execute(trx, sb), expectedErr)
		assert.ErrorIs(t, exe.checkFee(trx, sb), expectedErr)
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(lockTime, rndValAddr, proof)
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})
}

func TestReplay(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	executor := NewExecutor()
	sb := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	lockTime := sb.CurrentHeight()

	trx := tx.NewTransferTx(lockTime,
		rndAccAddr, ts.RandAccAddress(), 10000, 1000, "")
	ts.HelperSignTransaction(rndPrvKey, trx)

	err := executor.Execute(trx, sb)
	assert.NoError(t, err)
	err = executor.Execute(trx, sb)
	assert.ErrorIs(t, err, TransactionCommittedError{
		ID: trx.ID(),
	})
}

func TestChecker(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	lockTime := sb.CurrentHeight() + 1

	trx := tx.NewTransferTx(lockTime,
		rndAccAddr, ts.RandAccAddress(), 10000, 1000, "")
	ts.HelperSignTransaction(rndPrvKey, trx)

	err := executor.Execute(trx, sb)
	assert.ErrorIs(t, err, FutureLockTimeError{LockTime: lockTime})
	err = checker.Execute(trx, sb)
	assert.NoError(t, err)
}

func TestFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	exe := NewChecker()
	sb := sandbox.MockingSandbox(ts)

	tests := []struct {
		amount      amount.Amount
		fee         amount.Amount
		expectedFee amount.Amount
		expectErr   bool
	}{
		{1, 1, sb.TestParams.MinimumFee, true},
		{1, 1002, sb.TestParams.MinimumFee, true},
		{1, 998, sb.TestParams.MinimumFee, true},

		{1, 1001, sb.TestParams.MinimumFee, false},
		{1, 1000, sb.TestParams.MinimumFee, false},
		{1, 999, sb.TestParams.MinimumFee, false},

		{2 * 1e9, 100002, 200000, true},
		{2 * 1e9, 99998, 200000, true},

		{2 * 1e9, 200001, 200000, false},
		{2 * 1e9, 200000, 200000, false},
		{2 * 1e9, 199999, 200000, false},

		{1 * 1e12, 1000002, sb.TestParams.MaximumFee, true},
		{1 * 1e12, 999998, sb.TestParams.MaximumFee, true},

		{1 * 1e12, 1000001, sb.TestParams.MaximumFee, false},
		{1 * 1e12, 1000000, sb.TestParams.MaximumFee, false},
		{1 * 1e12, 999999, sb.TestParams.MaximumFee, false},

		{9_999_299_000, 999929, 999930, false}, // Block 66679
	}

	sender := ts.RandAccAddress()
	receiver := ts.RandAccAddress()
	for i, test := range tests {
		trx := tx.NewTransferTx(sb.CurrentHeight()+1, sender, receiver, test.amount, test.fee,
			"testing fee")
		err := exe.checkFee(trx, sb)

		if test.expectErr {
			assert.Error(t, err, "test %v failed. expected error", i)
		} else {
			assert.NoError(t, err, "test %v failed. unexpected error", i)
		}

		expectedFee := CalculateFee(test.amount, payload.TypeTransfer, sb.Params())
		assert.Equal(t, expectedFee, test.expectedFee, "test %v failed. invalid fee", i)
	}
}
