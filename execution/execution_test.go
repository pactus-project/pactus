package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestExecution(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	exe := NewExecutor()

	signer1 := ts.RandomSigner()
	addr1 := signer1.Address()
	acc1 := sb.MakeNewAccount(addr1)
	acc1.AddToBalance(100 * 1e9)
	sb.UpdateAccount(addr1, acc1)

	rcvAddr := ts.RandomAddress()
	block1 := sb.TestStore.AddTestBlock(1)
	block3 := sb.TestStore.AddTestBlock(3)
	block8635 := sb.TestStore.AddTestBlock(8635)
	block8641 := sb.TestStore.AddTestBlock(8641)
	block8642 := sb.TestStore.AddTestBlock(8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		assert.Error(t, exe.Execute(trx, sb))
		assert.Zero(t, exe.AccumulatedFee())
	})

	t.Run("Genesis stamp (expired), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block1.Stamp(), 1, addr1, rcvAddr, 1000, 1000,
			"expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("stamp is valid", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, exe.Execute(trx, sb))
	})

	t.Run("Subsidy transaction has an invalid stamp", func(t *testing.T) {
		trx := tx.NewSubsidyTx(block8641.Stamp(), 1, rcvAddr, 1000,
			"expired-stamp")
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Subsidy stamp is ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(block8642.Stamp(), 1, rcvAddr, 1000, "ok")
		assert.NoError(t, exe.Execute(trx, sb))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Invalid fee (subsidy tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, crypto.TreasuryAddress, rcvAddr, 1000, 1, "invalid fee")
		assert.Error(t, exe.Execute(trx, sb))
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Invalid fee (send tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 0, "invalid fee")
		assert.Error(t, exe.Execute(trx, sb))
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Sortition tx - Expired stamp, Should returns error", func(t *testing.T) {
		proof := ts.RandomProof()
		trx := tx.NewSortitionTx(block8635.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := ts.RandomProof()
		trx := tx.NewSortitionTx(block8642.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})
}

func TestChecker(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox(ts)

	block1000 := sb.TestStore.AddTestBlock(1000)

	t.Run("In strict mode transaction should be rejected.", func(t *testing.T) {
		signer := ts.RandomSigner()
		acc := sb.MakeNewAccount(signer.Address())
		acc.AddToBalance(10000000000)
		sb.UpdateAccount(signer.Address(), acc)
		valPub := sb.TestCommitteeSigners[0].PublicKey()

		trx := tx.NewBondTx(block1000.Stamp(), acc.Sequence()+1, signer.Address(),
			valPub.Address(), nil, 1000000000, 100000, "")
		signer.SignMsg(trx)
		assert.Error(t, executor.Execute(trx, sb))
		assert.NoError(t, checker.Execute(trx, sb))
	})
}

func TestLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox(ts)

	curHeight := 2 * sb.TestParams.TransactionToLiveInterval
	sb.TestStore.AddTestBlock(curHeight)

	t.Run("Should reject sortition transactions with lock time", func(t *testing.T) {
		pub, prv := ts.RandomBLSKeyPair()
		signer := crypto.NewSigner(prv)
		val := sb.MakeNewValidator(pub)
		sb.UpdateValidator(val)

		sb.TestAcceptSortition = true
		pld := &payload.SortitionPayload{
			Address: pub.Address(),
			Proof:   ts.RandomProof(),
		}
		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
		signer.SignMsg(trx)
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Should reject subsidy transactions with lock time", func(t *testing.T) {
		pld := &payload.TransferPayload{
			Sender:   crypto.TreasuryAddress,
			Receiver: ts.RandomAddress(),
			Amount:   1234,
		}
		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Should reject expired transactions", func(t *testing.T) {
		signer := ts.RandomSigner()
		acc := sb.MakeNewAccount(signer.Address())
		acc.AddToBalance(10000)
		sb.UpdateAccount(signer.Address(), acc)
		pld := &payload.TransferPayload{
			Sender:   signer.Address(),
			Receiver: ts.RandomAddress(),
			Amount:   1234,
		}

		trx := tx.NewLockTimeTx(curHeight-sb.TestParams.TransactionToLiveInterval, 1,
			pld, sb.TestParams.MinimumFee, "")
		signer.SignMsg(trx)
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Not finalized transaction", func(t *testing.T) {
		signer := ts.RandomSigner()
		acc := sb.MakeNewAccount(signer.Address())
		acc.AddToBalance(10000)
		sb.UpdateAccount(signer.Address(), acc)
		pld := &payload.TransferPayload{
			Sender:   signer.Address(),
			Receiver: ts.RandomAddress(),
			Amount:   1234,
		}

		trx1 := tx.NewLockTimeTx(curHeight+sb.TestParams.TransactionToLiveInterval, 1,
			pld, sb.TestParams.MinimumFee, "")
		signer.SignMsg(trx1)
		err := executor.Execute(trx1, sb)
		assert.Error(t, err)

		// In non-strict mode this transaction remains in pool
		err = checker.Execute(trx1, sb)
		assert.NoError(t, err)
	})
}

func TestFee(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	exe := NewChecker()
	sb := sandbox.MockingSandbox(ts)

	tests := []struct {
		amount          int64
		fee             int64
		expectedFee     int64
		expectedErrCode int
	}{
		{1, 1, sb.TestParams.MinimumFee, errors.ErrInvalidFee},
		{1, 1001, sb.TestParams.MinimumFee, errors.ErrInvalidFee},
		{1, 1000, sb.TestParams.MinimumFee, errors.ErrNone},

		{1 * 1e9, 1, 100000, errors.ErrInvalidFee},
		{1 * 1e9, 100001, 100000, errors.ErrInvalidFee},
		{1 * 1e9, 100000, 100000, errors.ErrNone},

		{1 * 1e12, 1, 1000000, errors.ErrInvalidFee},
		{1 * 1e12, 1000001, 1000000, errors.ErrInvalidFee},
		{1 * 1e12, 1000000, 1000000, errors.ErrNone},
	}

	sender := ts.RandomAddress()
	receiver := ts.RandomAddress()
	stamp := ts.RandomStamp()
	for i, test := range tests {
		trx := tx.NewTransferTx(stamp, 1, sender, receiver, test.amount, test.fee,
			"testing fee")
		err := exe.checkFee(trx, sb)

		assert.Equal(t, errors.Code(err), test.expectedErrCode,
			"test %v failed. unexpected error", i)

		assert.Equal(t, calculateFee(test.amount, sb), test.expectedFee,
			"test %v failed. invalid fee", i)
	}
}
