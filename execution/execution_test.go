package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
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
		trx := tx.NewTransferTx(lockTime, ts.RandAccAddress(), ts.RandAccAddress(), 1000, 0.1e9, "invalid-tx")
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.ErrInvalidAddress, errors.Code(err))
	})

	t.Run("Invalid fee (subsidy tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, crypto.TreasuryAddress, ts.RandAccAddress(), 1000, 1, "invalid fee")

		expectedErr := InvalidFeeError{Fee: 1, Expected: 0}
		assert.ErrorIs(t, exe.Execute(trx, sb), expectedErr)
		assert.ErrorIs(t, exe.checkFee(trx), expectedErr)
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(lockTime, rndValAddr, proof)
		ts.HelperSignTransaction(rndPrvKey, trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.ErrInvalidAddress, errors.Code(err))
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
