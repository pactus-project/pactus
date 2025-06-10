package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTransferLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandEd25519KeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sbx.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1000e9)
	sbx.UpdateAccount(rndAccAddr, rndAcc)
	_ = sbx.TestStore.AddTestBlock(8642)

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
			trx := tx.NewTransferTx(tt.lockTime, rndAccAddr, ts.RandAccAddress(), 1000, 1000)
			ts.HelperSignTransaction(rndPrvKey, trx)

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
	rndPubKey, rndPrvKey := ts.RandBLSKeyPair()
	rndValAddr := rndPubKey.ValidatorAddress()
	rndVal := sbx.MakeNewValidator(rndPubKey)
	rndVal.AddToStake(1000 * 1e9)
	sbx.UpdateValidator(rndVal)
	_ = sbx.TestStore.AddTestBlock(8642)

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
			trx := tx.NewSortitionTx(tt.lockTime, rndValAddr, ts.RandProof())
			ts.HelperSignTransaction(rndPrvKey, trx)

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
	_ = sbx.TestStore.AddTestBlock(8642)

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
			trx := tx.NewSubsidyTx(tt.lockTime, ts.RandAccAddress(), 1000)

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
	_ = sbx.TestStore.AddTestBlock(8642)
	lockTime := sbx.CurrentHeight()

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		randAddr := ts.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randAddr, ts.RandAccAddress(),
			ts.RandAmount(), ts.RandFee())

		err := Execute(trx, sbx)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: randAddr})
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000)
		err := Execute(trx, sbx)
		assert.NoError(t, err)

		assert.True(t, sbx.RecentTransaction(trx.ID()))
	})
}

func TestCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	_ = sbx.TestStore.AddTestBlock(8642)
	lockTime := sbx.CurrentHeight()

	t.Run("Invalid lock-time, Should return error", func(t *testing.T) {
		invalidLocoTme := lockTime + 1
		trx := tx.NewTransferTx(invalidLocoTme, crypto.TreasuryAddress, ts.RandAccAddress(), ts.RandAmount(), 0)

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, LockTimeInFutureError{LockTime: invalidLocoTme})
	})

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		randAddr := ts.RandAccAddress()
		trx := tx.NewTransferTx(lockTime, randAddr, ts.RandAccAddress(), ts.RandAmount(), ts.RandFee())

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{Address: randAddr})
	})

	t.Run("Invalid transaction, Should return error", func(t *testing.T) {
		valAddr := sbx.TestCommittee.Validators()[0].Address()
		sbx.TestAcceptSortition = false
		trx := tx.NewSortitionTx(lockTime, valAddr, ts.RandProof())

		err := CheckAndExecute(trx, sbx, true)
		assert.ErrorIs(t, err, executor.ErrInvalidSortitionProof)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(lockTime, ts.RandAccAddress(), 1000)
		err := CheckAndExecute(trx, sbx, true)
		assert.NoError(t, err)

		assert.True(t, sbx.RecentTransaction(trx.ID()))
	})
}

func TestReplay(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandEd25519KeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sbx.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1000e9)
	sbx.UpdateAccount(rndAccAddr, rndAcc)

	trx := ts.GenerateTestTransferTx(
		testsuite.TransactionWithEd25519Signer(rndPrvKey))

	err := Execute(trx, sbx)
	assert.NoError(t, err)

	err = CheckAndExecute(trx, sbx, false)
	assert.ErrorIs(t, err, TransactionCommittedError{
		ID: trx.ID(),
	})
}

func TestBatchTransfer(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	rndPubKey, rndPrvKey := ts.RandEd25519KeyPair()
	rndAccAddr := rndPubKey.AccountAddress()
	rndAcc := sbx.MakeNewAccount(rndAccAddr)
	rndAcc.AddToBalance(1000e9)
	sbx.UpdateAccount(rndAccAddr, rndAcc)

	sbx.TestStore.AddTestBlock(4_800_000 - 2)
	trx1 := ts.GenerateTestBatchTransferTx(
		testsuite.TransactionWithLockTime(sbx.CurrentHeight()),
		testsuite.TransactionWithEd25519Signer(rndPrvKey))
	err := CheckAndExecute(trx1, sbx, false)
	assert.ErrorIs(t, err, ErrBatchTransferNotAllowed)

	sbx.TestStore.AddTestBlock(4_800_000)
	trx2 := ts.GenerateTestBatchTransferTx(
		testsuite.TransactionWithLockTime(sbx.CurrentHeight()),
		testsuite.TransactionWithEd25519Signer(rndPrvKey))
	err = CheckAndExecute(trx2, sbx, false)
	assert.NoError(t, err)
}
