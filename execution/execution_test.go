package execution

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestExecution(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	exe := NewExecutor()

	signer1 := ts.RandSigner()
	addr1 := signer1.Address()
	acc1 := sb.MakeNewAccount(addr1)
	acc1.AddToBalance(100 * 1e9)
	sb.UpdateAccount(addr1, acc1)

	rcvAddr := ts.RandAddress()
	block1 := sb.TestStore.AddTestBlock(1)
	block3 := sb.TestStore.AddTestBlock(3)
	block8635 := sb.TestStore.AddTestBlock(8635)
	block8641 := sb.TestStore.AddTestBlock(8641)
	block8642 := sb.TestStore.AddTestBlock(8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
		assert.Zero(t, exe.AccumulatedFee())
	})

	t.Run("Genesis stamp (expired), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block1.Stamp(), 1, addr1, rcvAddr, 1000, 1000,
			"expired-stamp")
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("stamp is valid", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})

	t.Run("Subsidy transaction has an invalid stamp", func(t *testing.T) {
		trx := tx.NewSubsidyTx(block8641.Stamp(), 1, rcvAddr, 1000,
			"expired-stamp")
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Subsidy stamp is ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(block8642.Stamp(), 1, rcvAddr, 1000, "ok")
		err := exe.Execute(trx, sb)
		assert.NoError(t, err)
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
	})

	t.Run("Invalid fee (subsidy tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block8642.Stamp(), 2, crypto.TreasuryAddress, rcvAddr, 1000, 1, "invalid fee")
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Invalid fee (send tx), Should returns error", func(t *testing.T) {
		trx := tx.NewTransferTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 0, "invalid fee")
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidFee)
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Sortition tx - Expired stamp, Should returns error", func(t *testing.T) {
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(block8635.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := ts.RandProof()
		trx := tx.NewSortitionTx(block8642.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		err := exe.Execute(trx, sb)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})
}

func TestChecker(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox(ts)

	height := uint32(1000)
	block1000 := sb.TestStore.AddTestBlock(height)

	t.Run("In strict mode transaction should be rejected.", func(t *testing.T) {
		signer := ts.RandSigner()
		acc := sb.MakeNewAccount(signer.Address())
		acc.AddToBalance(10000000000)
		sb.UpdateAccount(signer.Address(), acc)
		valPub := sb.TestCommitteeSigners[0].PublicKey()

		trx := tx.NewBondTx(block1000.Stamp(), height+1, signer.Address(),
			valPub.Address(), nil, 1000000000, 100000, "")
		signer.SignMsg(trx)
		assert.Error(t, executor.Execute(trx, sb))
		assert.NoError(t, checker.Execute(trx, sb))
	})
}

// func TestLockTime(t *testing.T) {
// 	ts := testsuite.NewTestSuite(t)

// 	executor := NewExecutor()
// 	checker := NewChecker()
// 	sb := sandbox.MockingSandbox(ts)

// 	curHeight := 2 * sb.TestParams.TransactionToLiveInterval
// 	sb.TestStore.AddTestBlock(curHeight)

// 	t.Run("Should reject sortition transactions with lock time", func(t *testing.T) {
// 		pub, prv := ts.RandomBLSKeyPair()
// 		signer := crypto.NewSigner(prv)
// 		val := sb.MakeNewValidator(pub)
// 		sb.UpdateValidator(val)

// 		sb.TestAcceptSortition = true
// 		pld := &payload.SortitionPayload{
// 			Address: pub.Address(),
// 			Proof:   ts.RandomProof(),
// 		}
// 		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
// 		signer.SignMsg(trx)
// 		err := executor.Execute(trx, sb)
// 		assert.Error(t, err)
// 	})

// 	t.Run("Should reject subsidy transactions with lock time", func(t *testing.T) {
// 		pld := &payload.TransferPayload{
// 			Sender:   crypto.TreasuryAddress,
// 			Receiver: ts.RandomAddress(),
// 			Amount:   1234,
// 		}
// 		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
// 		err := executor.Execute(trx, sb)
// 		assert.Error(t, err)
// 	})

// 	t.Run("Should reject expired transactions", func(t *testing.T) {
// 		signer := ts.RandSigner()
// 		acc := sb.MakeNewAccount(signer.Address())
// 		acc.AddToBalance(10000)
// 		sb.UpdateAccount(signer.Address(), acc)
// 		pld := &payload.TransferPayload{
// 			Sender:   signer.Address(),
// 			Receiver: ts.RandAddress(),
// 			Amount:   1234,
// 		}

// 		trx := tx.NewLockTimeTx(curHeight-sb.TestParams.TransactionToLiveInterval, 1,
// 			pld, sb.TestParams.MinimumFee, "")
// 		signer.SignMsg(trx)
// 		err := executor.Execute(trx, sb)
// 		assert.Error(t, err)
// 	})

// 	t.Run("Not finalized transaction", func(t *testing.T) {
// 		signer := ts.RandomSigner()
// 		acc := sb.MakeNewAccount(signer.Address())
// 		acc.AddToBalance(10000)
// 		sb.UpdateAccount(signer.Address(), acc)
// 		pld := &payload.TransferPayload{
// 			Sender:   signer.Address(),
// 			Receiver: ts.RandomAddress(),
// 			Amount:   1234,
// 		}

// 		trx1 := tx.NewLockTimeTx(curHeight+sb.TestParams.TransactionToLiveInterval, 1,
// 			pld, sb.TestParams.MinimumFee, "")
// 		signer.SignMsg(trx1)
// 		err := executor.Execute(trx1, sb)
// 		assert.Error(t, err)

// 		// In non-strict mode this transaction remains in pool
// 		err = checker.Execute(trx1, sb)
// 		assert.NoError(t, err)
// 	})
// }

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

	sender := ts.RandAddress()
	receiver := ts.RandAddress()
	stamp := ts.RandStamp()
	for i, test := range tests {
		trx := tx.NewTransferTx(stamp, sb.CurrentHeight()+1, sender, receiver, test.amount, test.fee,
			"testing fee")
		err := exe.checkFee(trx, sb)

		assert.Equal(t, errors.Code(err), test.expectedErrCode,
			"test %v failed. unexpected error", i)

		assert.Equal(t, CalculateFee(test.amount, sb.Params()), test.expectedFee,
			"test %v failed. invalid fee", i)
	}
}
