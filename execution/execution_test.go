package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTransferLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandEd25519KeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1000 * 1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Transaction has expired LockTime  (-8641)",
			lockTime:     sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1,
			strictErr:    LockTimeExpiredError{sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1},
			nonStrictErr: LockTimeExpiredError{sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval - 1},
		},
		{
			name:         "Transaction has valid LockTime (-8640)",
			lockTime:     sb.CurrentHeight() - sb.TestParams.TransactionToLiveInterval,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime (-88)",
			lockTime:     sb.CurrentHeight() - 88,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has valid LockTime (0)",
			lockTime:     sb.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Transaction has future LockTime (+1)",
			lockTime:     sb.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sb.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewTransferTx(tc.lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000)
			ts.HelperSignTransaction(rndPrvKey, trx)

			strictErr := CheckLockTime(trx, sb, true)
			assert.ErrorIs(t, strictErr, tc.strictErr)

			nonStrictErr := CheckLockTime(trx, sb, false)
			assert.ErrorIs(t, nonStrictErr, tc.nonStrictErr)
		})
	}
}

func TestSortitionLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	sb.TestAcceptSortition = true
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndValAddr := rndPubKey.ValidatorAddress()
	rndVal := sb.MakeNewValidator(rndPubKey)
	rndVal.AddToStake(1000 * 1e9)
	sb.UpdateValidator(rndVal)
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Sortition transaction has expired LockTime (-8)",
			lockTime:     sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1,
			strictErr:    LockTimeExpiredError{sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1},
			nonStrictErr: LockTimeExpiredError{sb.CurrentHeight() - sb.TestParams.SortitionInterval - 1},
		},
		{
			name:         "Sortition transaction has valid LockTime (-7)",
			lockTime:     sb.CurrentHeight() - sb.TestParams.SortitionInterval,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has valid LockTime (-1)",
			lockTime:     sb.CurrentHeight() - 1,
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has valid LockTime (0)",
			lockTime:     sb.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Sortition transaction has future LockTime (+1)",
			lockTime:     sb.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sb.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewSortitionTx(tc.lockTime, rndValAddr, ts.RandProof())
			ts.HelperSignTransaction(rndPrvKey, trx)

			strictErr := CheckLockTime(trx, sb, true)
			assert.ErrorIs(t, strictErr, tc.strictErr)

			nonStrictErr := CheckLockTime(trx, sb, false)
			assert.ErrorIs(t, nonStrictErr, tc.nonStrictErr)
		})
	}
}

func TestSubsidyLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	_ = sb.TestStore.AddTestBlock(8642)

	tests := []struct {
		name         string
		lockTime     uint32
		strictErr    error
		nonStrictErr error
	}{
		{
			name:         "Subsidy transaction has expired LockTime (-1)",
			lockTime:     sb.CurrentHeight() - 1,
			strictErr:    LockTimeExpiredError{sb.CurrentHeight() - 1},
			nonStrictErr: LockTimeExpiredError{sb.CurrentHeight() - 1},
		},
		{
			name:         "Subsidy transaction has valid LockTime (0)",
			lockTime:     sb.CurrentHeight(),
			strictErr:    nil,
			nonStrictErr: nil,
		},
		{
			name:         "Subsidy transaction has future LockTime (+1)",
			lockTime:     sb.CurrentHeight() + 1,
			strictErr:    LockTimeInFutureError{sb.CurrentHeight() + 1},
			nonStrictErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trx := tx.NewSubsidyTx(tc.lockTime, ts.RandAccAddress(), 1000)

			strictErr := CheckLockTime(trx, sb, true)
			assert.ErrorIs(t, strictErr, tc.strictErr)

			nonStrictErr := CheckLockTime(trx, sb, false)
			assert.ErrorIs(t, nonStrictErr, tc.nonStrictErr)
		})
	}
}

func TestCheckFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		name        string
		trx         *tx.Tx
		expectedErr error
	}{
		{
			name: "Subsidy transaction with fee",
			trx: tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(),
				ts.RandAmount(), 1),
			expectedErr: InvalidFeeError{Fee: 1, Expected: 0},
		},
		{
			name: "Subsidy transaction without fee",
			trx: tx.NewTransferTx(ts.RandHeight(), crypto.TreasuryAddress, ts.RandAccAddress(),
				ts.RandAmount(), 0),
			expectedErr: nil,
		},
		{
			name: "Transfer transaction with fee",
			trx: tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(),
				ts.RandAmount(), 0),
			expectedErr: nil,
		},
		{
			name: "Transfer transaction without fee",
			trx: tx.NewTransferTx(ts.RandHeight(), ts.RandAccAddress(), ts.RandAccAddress(),
				ts.RandAmount(), 0),
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckFee(tc.trx)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestExecute(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	_ = sb.TestStore.AddTestBlock(8642)
	lockTime := sb.CurrentHeight()

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		randAddr := ts.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randAddr, ts.RandAccAddress(),
			ts.RandAmount(), ts.RandFee())

		err := Execute(trx, sb)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: randAddr})
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000)
		err := Execute(trx, sb)
		assert.NoError(t, err)

		assert.True(t, sb.RecentTransaction(trx.ID()))
	})
}

func TestCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	_ = sb.TestStore.AddTestBlock(8642)
	lockTime := sb.CurrentHeight()

	t.Run("Invalid lock-time, Should return error", func(t *testing.T) {
		invalidLocoTme := lockTime + 1
		trx := tx.NewTransferTx(invalidLocoTme, crypto.TreasuryAddress, ts.RandAccAddress(), ts.RandAmount(), 0)

		err := CheckAndExecute(trx, sb, true)
		assert.ErrorIs(t, err, LockTimeInFutureError{LockTime: invalidLocoTme})
	})

	t.Run("Invalid fee, Should return error", func(t *testing.T) {
		invalidFee := amount.Amount(1)
		trx := tx.NewTransferTx(lockTime, crypto.TreasuryAddress, ts.RandAccAddress(), ts.RandAmount(), invalidFee)

		err := CheckAndExecute(trx, sb, true)
		assert.ErrorIs(t, err, InvalidFeeError{Fee: invalidFee, Expected: 0})
	})

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		randAddr := ts.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randAddr, ts.RandAccAddress(), ts.RandAmount(), ts.RandFee())

		err := CheckAndExecute(trx, sb, true)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: randAddr})
	})

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		valAddr := sb.TestCommittee.Validators()[0].Address()
		sb.TestAcceptSortition = false
		trx := tx.NewSortitionTx(lockTime, valAddr, ts.RandProof())

		err := CheckAndExecute(trx, sb, true)
		assert.ErrorIs(t, err, executor.ErrInvalidSortitionProof)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000)
		err := CheckAndExecute(trx, sb, true)
		assert.NoError(t, err)

		assert.True(t, sb.RecentTransaction(trx.ID()))
	})
}

func TestReplay(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandEd25519KeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sb.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1e9)
	sb.UpdateAccount(rndAccAddr, rndAcc)
	lockTime := sb.CurrentHeight()

	trx := tx.NewTransferTx(lockTime,
		rndAccAddr, ts.RandAccAddress(), 10000, 1000)
	ts.HelperSignTransaction(rndPrvKey, trx)

	err := Execute(trx, sb)
	assert.NoError(t, err)

	err = CheckAndExecute(trx, sb, false)
	assert.ErrorIs(t, err, TransactionCommittedError{
		ID: trx.ID(),
	})
}
