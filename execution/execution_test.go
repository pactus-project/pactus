package execution

import (
	"testing"

	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTransferLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Transaction has expired LockTime  (-8641)",
			lockTime:     sbx.CurrentHeight() - sbx.TestParams.TransactionToLiveInterval - 1,
			strictErr:    LockTimeExpiredError{sbx.CurrentHeight() - sbx.TestParams.TransactionToLiveInterval - 1},
			nonStrictErr: LockTimeExpiredError{sbx.CurrentHeight() - sbx.TestParams.TransactionToLiveInterval - 1},
		},
		{
			name:         "Transaction has valid LockTime (-8640)",
			lockTime:     sbx.CurrentHeight() - sbx.TestParams.TransactionToLiveInterval,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime (-88)",
			lockTime:     sbx.CurrentHeight() - 88,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime (0)",
			lockTime:     sbx.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has future LockTime (+1)",
			lockTime:     sbx.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sbx.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestTransferTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			assert.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			assert.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestSortitionLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestAcceptSortition = true
	sbx.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Sortition transaction has expired LockTime (-8)",
			lockTime:     sbx.CurrentHeight() - sbx.TestParams.SortitionInterval - 1,
			strictErr:    LockTimeExpiredError{sbx.CurrentHeight() - sbx.TestParams.SortitionInterval - 1},
			nonStrictErr: LockTimeExpiredError{sbx.CurrentHeight() - sbx.TestParams.SortitionInterval - 1},
		},
		{
			name:         "Sortition transaction has valid LockTime (-7)",
			lockTime:     sbx.CurrentHeight() - sbx.TestParams.SortitionInterval,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has valid LockTime (-1)",
			lockTime:     sbx.CurrentHeight() - 1,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has valid LockTime (0)",
			lockTime:     sbx.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has future LockTime (+1)",
			lockTime:     sbx.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sbx.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestSortitionTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			assert.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			assert.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestSubsidyLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Subsidy transaction has expired LockTime (-1)",
			lockTime:     sbx.CurrentHeight() - 1,
			strictErr:    LockTimeExpiredError{sbx.CurrentHeight() - 1},
			nonStrictErr: LockTimeExpiredError{sbx.CurrentHeight() - 1},
		},
		{
			name:         "Subsidy transaction has valid LockTime (0)",
			lockTime:     sbx.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Subsidy transaction has future LockTime (+1)",
			lockTime:     sbx.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sbx.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := ts.GenerateTestSubsidyTx(
				testsuite.TransactionWithLockTime(tt.lockTime))

			strictErr := CheckLockTime(trx, sbx, true)
			assert.ErrorIs(t, strictErr, tt.strictErr)

			nonStrictErr := CheckLockTime(trx, sbx, false)
			assert.ErrorIs(t, nonStrictErr, tt.nonStrictErr)
		})
	}
}

func TestExecute(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestStore.AddTestBlock(8642)
	knownPub, knownSigner := ts.RandEd25519KeyPair()
	sbx.TestStore.AddTestAccount(
		testsuite.AccountWithAddress(knownPub.AccountAddress()))
	lockTime := sbx.CurrentHeight()

	t.Run("Unknown Signer", func(t *testing.T) {
		_, unknownSigner := ts.RandKeyPair()
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithSigner(unknownSigner))

		err := Execute(trx, sbx)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: trx.Payload().Signer()})
	})

	t.Run("Valid Transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithSigner(knownSigner))

		err := Execute(trx, sbx)
		assert.NoError(t, err)

		assert.True(t, sbx.RecentTransaction(trx.ID()))
	})
}

func TestCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestStore.AddTestBlock(8642)
	knownPub, knownSigner := ts.RandEd25519KeyPair()
	_, testAcc := sbx.TestStore.AddTestAccount(
		testsuite.AccountWithAddress(knownPub.AccountAddress()))
	lockTime := sbx.CurrentHeight()

	t.Run("Unknown Sender", func(t *testing.T) {
		_, unknownSigner := ts.RandKeyPair()
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithSigner(unknownSigner))

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: trx.Payload().Signer()})
	})

	t.Run("Invalid lock-time, Should return error", func(t *testing.T) {
		invalidLockTime := lockTime + 1
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(invalidLockTime),
			testsuite.TransactionWithSigner(knownSigner))

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, LockTimeInFutureError{LockTime: invalidLockTime})
	})

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithSigner(knownSigner),
			testsuite.TransactionWithAmount(testAcc.Balance()+1))

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, executor.ErrInsufficientFunds)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithSigner(knownSigner))

		err := CheckAndExecute(trx, sbx, true)
		assert.NoError(t, err)
		assert.True(t, sbx.RecentTransaction(trx.ID()))
	})
}

func TestReplay(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	sbx.TestStore.AddTestBlock(8642)
	knownPub, knownSigner := ts.RandEd25519KeyPair()
	sbx.TestStore.AddTestAccount(
		testsuite.AccountWithAddress(knownPub.AccountAddress()))

	trx := ts.GenerateTestTransferTx(
		testsuite.TransactionWithSigner(knownSigner))

	err := Execute(trx, sbx)
	assert.NoError(t, err)

	err = CheckAndExecute(trx, sbx, false)
	assert.ErrorIs(t, err, TransactionCommittedError{
		ID: trx.ID(),
	})
}
