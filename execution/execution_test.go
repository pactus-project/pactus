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

func TestTransferLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	exe := NewExecutor()
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1000 * 1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name     string
		lockTime uint32
		wantErr  error
	}{
		{
			name:     "Transaction has invalid LockTime  (-8641)",
			lockTime: sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1,
			wantErr:  LockTimeExpiredError{LockTime: sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1},
		},
		{
			name:     "Transaction has valid LockTime (-8640)",
			lockTime: sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval,
			wantErr:  nil,
		},
		{
			name:     "Transaction has valid LockTime (-88)",
			lockTime: sb.CurrentHeight() - 88,
			wantErr:  nil,
		},
		{
			name:     "Transaction has valid LockTime (0)",
			lockTime: sb.CurrentHeight(),
			wantErr:  nil,
		},
		{
			name:     "Transaction has invalid LockTime (+1)",
			lockTime: sb.CurrentHeight() + 1,
			wantErr:  LockTimeInFutureError{LockTime: sb.CurrentHeight() + 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewTransferTx(tc.lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000, "")
			ts.HelperSignTransaction(rndPrvKey, trx)
			err := exe.Execute(trx, sb)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortitionLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	sb.TestAcceptSortition = true
	exe := NewExecutor()
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndValAddr := rndPubKey.ValidatorAddress()
	rndVal := sb.MakeNewValidator(rndPubKey)
	rndVal.AddToStake(1000 * 1e9)
	sb.UpdateValidator(rndVal)
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name     string
		lockTime uint32
		wantErr  error
	}{
		{
			name:     "Sortition transaction has valid LockTime (-8)",
			lockTime: sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1,
			wantErr:  LockTimeExpiredError{LockTime: sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1},
		},
		{
			name:     "Sortition transaction has valid LockTime (-7)",
			lockTime: sb.CurrentHeight() - sb.TestParams.SortitionInterval,
			wantErr:  nil,
		},
		{
			name:     "Sortition transaction has valid LockTime (-1)",
			lockTime: sb.CurrentHeight() - 1,
			wantErr:  nil,
		},
		{
			name:     "Sortition transaction has valid LockTime (0)",
			lockTime: sb.CurrentHeight(),
			wantErr:  nil,
		},
		{
			name:     "Sortition transaction has invalid LockTime (+1)",
			lockTime: sb.CurrentHeight() + 1,
			wantErr:  LockTimeInFutureError{LockTime: sb.CurrentHeight() + 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewSortitionTx(tc.lockTime, rndValAddr, ts.RandProof())
			ts.HelperSignTransaction(rndPrvKey, trx)
			err := exe.Execute(trx, sb)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSubsidyLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	exe := NewExecutor()
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name     string
		lockTime uint32
		wantErr  error
	}{
		{
			name:     "Subsidy transaction has invalid LockTime (-1)",
			lockTime: sb.CurrentHeight() - 1,
			wantErr:  LockTimeExpiredError{LockTime: sb.CurrentHeight() - 1},
		},
		{
			name:     "Subsidy transaction has valid LockTime (0)",
			lockTime: sb.CurrentHeight(),
			wantErr:  nil,
		},
		{
			name:     "Subsidy transaction has invalid LockTime (+1)",
			lockTime: sb.CurrentHeight() + 1,
			wantErr:  LockTimeInFutureError{LockTime: sb.CurrentHeight() + 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewSubsidyTx(tc.lockTime, ts.RandAccAddress(), 1000, "subsidy-test")
			err := exe.Execute(trx, sb)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
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
	assert.ErrorIs(t, err, LockTimeInFutureError{LockTime: lockTime})
	err = checker.Execute(trx, sb)
	assert.NoError(t, err)
}
