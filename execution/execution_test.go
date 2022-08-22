package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
)

func TestExecution(t *testing.T) {
	sb := sandbox.MockingSandbox()
	exe := NewExecutor()

	signer1 := bls.GenerateTestSigner()
	addr1 := signer1.Address()
	acc1 := sb.MakeNewAccount(addr1)
	acc1.AddToBalance(10000000000)
	sb.UpdateAccount(acc1)

	rcvAddr := crypto.GenerateTestAddress()
	block1 := sb.TestStore.AddTestBlock(1)
	block3 := sb.TestStore.AddTestBlock(3)
	block8635 := sb.TestStore.AddTestBlock(8635)
	block8641 := sb.TestStore.AddTestBlock(8641)
	block8642 := sb.TestStore.AddTestBlock(8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, exe.Execute(trx, sb))
		assert.Zero(t, exe.AccumulatedFee())
	})

	t.Run("Genesis stamp (expired), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block1.Stamp(), 1, addr1, rcvAddr, 1000, 1000,
			"expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("stamp is valid", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
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
		trx := tx.NewSendTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Invalid fee (subsidy tx), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 2, crypto.TreasuryAddress, rcvAddr, 1000, 1, "invalid fee")
		assert.Error(t, exe.Execute(trx, sb))
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Invalid fee (send tx), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 2, addr1, rcvAddr, 1000, 0, "invalid fee")
		assert.Error(t, exe.Execute(trx, sb))
		assert.Error(t, exe.checkFee(trx, sb))
	})

	t.Run("Sortition tx - Expired stamp, Should returns error", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(block8635.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(block8642.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, exe.Execute(trx, sb))
	})
}

func TestChecker(t *testing.T) {
	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox()

	block1000 := sb.TestStore.AddTestBlock(1000)

	t.Run("In strict mode transaction should be rejected.", func(t *testing.T) {
		pub, prv := bls.GenerateTestKeyPair()
		signer := crypto.NewSigner(prv)
		acc := sb.MakeNewAccount(pub.Address())
		acc.AddToBalance(10000)
		sb.UpdateAccount(acc)
		valPub := sb.TestCommitteeSigners[0].PublicKey()

		trx := tx.NewBondTx(block1000.Stamp(), acc.Sequence()+1, acc.Address(),
			valPub.Address(), valPub.(*bls.PublicKey), 1000, 1000, "")
		signer.SignMsg(trx)
		assert.Error(t, executor.Execute(trx, sb))
		assert.NoError(t, checker.Execute(trx, sb))
	})
}

func TestLockTime(t *testing.T) {
	executor := NewExecutor()
	checker := NewChecker()
	sb := sandbox.MockingSandbox()

	curHeight := 2 * sb.TestParams.TransactionToLiveInterval
	sb.TestStore.AddTestBlock(curHeight)

	t.Run("Should reject sortition transaxtions with lock time", func(t *testing.T) {
		pub, prv := bls.GenerateTestKeyPair()
		signer := crypto.NewSigner(prv)
		val := sb.MakeNewValidator(pub)
		sb.UpdateValidator(val)

		sb.AcceptTestSortition = true
		pld := &payload.SortitionPayload{
			Address: pub.Address(),
			Proof:   sortition.GenerateRandomProof(),
		}
		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
		signer.SignMsg(trx)
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Should reject subsidy transaxtions with lock time", func(t *testing.T) {
		pld := &payload.SendPayload{
			Sender:   crypto.TreasuryAddress,
			Receiver: crypto.GenerateTestAddress(),
			Amount:   1234,
		}
		trx := tx.NewLockTimeTx(curHeight+10, 1, pld, 0, "")
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Should reject expired transaxtions", func(t *testing.T) {
		pub, prv := bls.GenerateTestKeyPair()
		signer := crypto.NewSigner(prv)
		acc := sb.MakeNewAccount(pub.Address())
		acc.AddToBalance(10000)
		sb.UpdateAccount(acc)
		pld := &payload.SendPayload{
			Sender:   acc.Address(),
			Receiver: crypto.GenerateTestAddress(),
			Amount:   1234,
		}

		trx := tx.NewLockTimeTx(curHeight-sb.TestParams.TransactionToLiveInterval, 1,
			pld, sb.TestParams.MinimumFee, "")
		signer.SignMsg(trx)
		err := executor.Execute(trx, sb)
		assert.Error(t, err)
	})

	t.Run("Not finalized transaction", func(t *testing.T) {
		pub, prv := bls.GenerateTestKeyPair()
		signer := crypto.NewSigner(prv)
		acc := sb.MakeNewAccount(pub.Address())
		acc.AddToBalance(10000)
		sb.UpdateAccount(acc)
		pld := &payload.SendPayload{
			Sender:   acc.Address(),
			Receiver: crypto.GenerateTestAddress(),
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
