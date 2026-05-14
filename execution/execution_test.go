package execution

import (
	"errors"
	"testing"

	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestTransferLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.NewMockSandbox(ts.Ctrl)
	params := &param.Params{
		TransactionToLiveInterval: 10,
	}
	currentHeight := types.Height(13)

	sbx.EXPECT().Params().Return(params).AnyTimes()
	sbx.EXPECT().CurrentHeight().Return(currentHeight).AnyTimes()

	tests := []struct {
		name         string
		lockTime     types.Height
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Transaction has expired LockTime",
			lockTime:     2,
			strictErr:    LockTimeExpiredError{2},
			nonStrictErr: LockTimeExpiredError{2},
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     3,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     4,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     13,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime in future",
			lockTime:     14,
			strictErr:    LockTimeInFutureError{14},
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime in future",
			lockTime:     1014,
			strictErr:    LockTimeInFutureError{1014},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestTransferTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			require.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			require.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestSortitionLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.NewMockSandbox(ts.Ctrl)
	params := &param.Params{
		SortitionInterval: 10,
	}
	currentHeight := types.Height(13)

	sbx.EXPECT().Params().Return(params).AnyTimes()
	sbx.EXPECT().CurrentHeight().Return(currentHeight).AnyTimes()

	tests := []struct {
		name         string
		lockTime     types.Height
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Transaction has expired LockTime",
			lockTime:     2,
			strictErr:    LockTimeExpiredError{2},
			nonStrictErr: LockTimeExpiredError{2},
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     3,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     4,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime",
			lockTime:     13,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime in future",
			lockTime:     14,
			strictErr:    LockTimeInFutureError{14},
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime in future",
			lockTime:     1014,
			strictErr:    LockTimeInFutureError{1014},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestSortitionTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			require.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			require.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestSubsidyLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	params := &param.Params{
		TransactionToLiveInterval: 10,
		SortitionInterval:         10,
	}
	currentHeight := types.Height(13)

	sbx := sandbox.NewMockSandbox(ts.Ctrl)
	sbx.EXPECT().Params().Return(params).AnyTimes()
	sbx.EXPECT().CurrentHeight().Return(currentHeight).AnyTimes()

	tests := []struct {
		name         string
		lockTime     types.Height
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Subsidy transaction has expired LockTime",
			lockTime:     12,
			strictErr:    LockTimeExpiredError{12},
			nonStrictErr: LockTimeExpiredError{12},
		},
		{
			name:         "Subsidy transaction has valid LockTime",
			lockTime:     13,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Subsidy transaction has future LockTime (+1)",
			lockTime:     14,
			strictErr:    LockTimeInFutureError{14},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestSubsidyTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			require.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			require.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestExecute(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	sbx := sandbox.NewMockSandbox(ts.Ctrl)

	t.Run("Invalid transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()

		executor.DefaultFactory = func(*tx.Tx, sandbox.Sandbox) (executor.Executor, error) {
			return nil, errors.New("some Error")
		}

		err := Execute(trx, sbx)
		require.Error(t, err)
	})

	mockExe := executor.NewMockExecutor(ts.Ctrl)
	executor.DefaultFactory = func(*tx.Tx, sandbox.Sandbox) (executor.Executor, error) {
		return mockExe, nil
	}

	t.Run("Valid transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()

		mockExe.EXPECT().Execute().Return().Times(1)
		sbx.EXPECT().CommitTransaction(trx).Return().Times(1)

		err := Execute(trx, sbx)
		require.NoError(t, err)
	})
}

func TestCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	sbx := sandbox.NewMockSandbox(ts.Ctrl)
	lockTime := ts.RandHeight()

	t.Run("Invalid transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()

		executor.DefaultFactory = func(*tx.Tx, sandbox.Sandbox) (executor.Executor, error) {
			return nil, errors.New("some Error")
		}

		err := CheckAndExecute(trx, sbx, true)
		require.Error(t, err)
	})

	mockExe := executor.NewMockExecutor(ts.Ctrl)
	executor.DefaultFactory = func(*tx.Tx, sandbox.Sandbox) (executor.Executor, error) {
		return mockExe, nil
	}

	t.Run("Banned account", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()

		sbx.EXPECT().IsBanned(trx.Payload().Signer()).Return(true).Times(1)

		err := CheckAndExecute(trx, sbx, true)
		require.ErrorIs(t, err, SignerBannedError{Address: trx.Payload().Signer()})
	})

	t.Run("Replay transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime))

		sbx.EXPECT().IsBanned(trx.Payload().Signer()).Return(false).Times(1)
		sbx.EXPECT().RecentTransaction(trx.ID()).Return(true).Times(1)

		err := CheckAndExecute(trx, sbx, true)
		require.ErrorIs(t, err, TransactionCommittedError{ID: trx.ID()})
	})

	t.Run("Invalid lock-time", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime))

		sbx.EXPECT().IsBanned(trx.Payload().Signer()).Return(false).Times(1)
		sbx.EXPECT().RecentTransaction(trx.ID()).Return(false).Times(1)
		sbx.EXPECT().CurrentHeight().Return(lockTime - 1).Times(3)
		sbx.EXPECT().Params().Return(&param.Params{}).Times(1)

		err := CheckAndExecute(trx, sbx, true)
		require.ErrorIs(t, err, LockTimeInFutureError{LockTime: lockTime})
	})

	t.Run("Check fails", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime))

		sbx.EXPECT().IsBanned(trx.Payload().Signer()).Return(false).Times(1)
		sbx.EXPECT().RecentTransaction(trx.ID()).Return(false).Times(1)
		sbx.EXPECT().Params().Return(&param.Params{}).Times(1)
		sbx.EXPECT().CurrentHeight().Return(lockTime).Times(3)
		mockExe.EXPECT().Check(true).Return(errors.New("Some error")).Times(1)

		err := CheckAndExecute(trx, sbx, true)
		require.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime))

		sbx.EXPECT().IsBanned(trx.Payload().Signer()).Return(false).Times(1)
		sbx.EXPECT().RecentTransaction(trx.ID()).Return(false).Times(1)
		sbx.EXPECT().Params().Return(&param.Params{}).Times(1)
		sbx.EXPECT().CurrentHeight().Return(lockTime).Times(3)
		mockExe.EXPECT().Check(true).Return(nil).Times(1)
		mockExe.EXPECT().Execute().Return().Times(1)
		sbx.EXPECT().CommitTransaction(trx).Return().Times(1)

		err := CheckAndExecute(trx, sbx, true)
		require.NoError(t, err)
	})
}
