package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecution(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	tSandbox := sandbox.MockingSandbox()
	tExec := NewExecutor()

	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(21*1e14 - 10000000000)
	tSandbox.UpdateAccount(acc0)

	signer1 := bls.GenerateTestSigner()
	addr1 := signer1.Address()
	acc1 := account.NewAccount(addr1, 1)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	rcvAddr := crypto.GenerateTestAddress()
	block1 := tSandbox.TestStore.AddTestBlock(1)
	block2 := tSandbox.TestStore.AddTestBlock(2)
	block3 := tSandbox.TestStore.AddTestBlock(3)
	block8635 := tSandbox.TestStore.AddTestBlock(8635)
	block8641 := tSandbox.TestStore.AddTestBlock(8641)
	block8642 := tSandbox.TestStore.AddTestBlock(8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Genesis stamp (expired), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block1.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("stamp is valid", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase has an invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewMintbaseTx(block8641.Stamp(), 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase stamp is ok", func(t *testing.T) {
		trx := tx.NewMintbaseTx(block8642.Stamp(), 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block2.Stamp(), 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block2.Stamp(), 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee (mintbase), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block2.Stamp(), 2, crypto.TreasuryAddress, rcvAddr, 1000, 1, "invalid fee")
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Error(t, tExec.checkFee(trx, tSandbox))
	})

	t.Run("Invalid fee (send), Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block2.Stamp(), 2, addr1, rcvAddr, 1000, 0, "invalid fee")
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Error(t, tExec.checkFee(trx, tSandbox))
	})

	t.Run("Sortition tx - Expired stamp, Should returns error", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(block8635.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(block8642.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})
}

func TestChecker(t *testing.T) {
	executor := NewExecutor()
	checker := NewChecker()
	tSandbox := sandbox.MockingSandbox()

	block1000 := tSandbox.TestStore.AddTestBlock(1000)

	t.Run("Accept bond transaction for future blocks", func(t *testing.T) {
		pub := tSandbox.Committee().Proposer(0).PublicKey()
		acc, signer := account.GenerateTestAccount(1)
		tSandbox.UpdateAccount(acc)

		trx := tx.NewBondTx(block1000.Stamp(), acc.Sequence()+1, acc.Address(), pub, 1000, 1000, "")
		signer.SignMsg(trx)
		assert.Error(t, executor.Execute(trx, tSandbox))
		assert.NoError(t, checker.Execute(trx, tSandbox))
	})
}
