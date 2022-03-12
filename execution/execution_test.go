package execution

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
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
	block1 := block.GenerateTestBlock(nil, nil)
	block2 := block.GenerateTestBlock(nil, nil)
	block3 := block.GenerateTestBlock(nil, nil)
	block8635 := block.GenerateTestBlock(nil, nil)
	block8640 := block.GenerateTestBlock(nil, nil)
	block8641 := block.GenerateTestBlock(nil, nil)
	block8642 := block.GenerateTestBlock(nil, nil)
	tSandbox.AddTestBlock(1, block1)
	tSandbox.AddTestBlock(2, block2)
	tSandbox.AddTestBlock(3, block3)
	tSandbox.AddTestBlock(8635, block8635)
	tSandbox.AddTestBlock(8640, block8640)
	tSandbox.AddTestBlock(8641, block8641)
	tSandbox.AddTestBlock(8642, block8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(block1.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Good stamp", func(t *testing.T) {
		trx := tx.NewSendTx(block3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewMintbaseTx(block8641.Stamp(), 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase stamp is ok", func(t *testing.T) {
		trx := tx.NewMintbaseTx(block8642.Stamp(), 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 1025)
		trx := tx.NewSendTx(block8641.Stamp(), 2, addr1, rcvAddr, 1000, 1000, bigMemo)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
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

	t.Run("Sortition tx - Invalid stamp, Should returns error", func(t *testing.T) {
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

	block1000 := block.GenerateTestBlock(nil, nil)
	tSandbox.AddTestBlock(1000, block1000)

	t.Run("Accept bond transaction for future blocks", func(t *testing.T) {
		pub := tSandbox.Committee().Proposer(0).PublicKey()
		acc, signer := account.GenerateTestAccount(1)
		tSandbox.TestAccounts[acc.Address()] = acc

		trx := tx.NewBondTx(block1000.Stamp(), acc.Sequence()+1, acc.Address(), pub, 1000, 1000, "")
		signer.SignMsg(trx)
		assert.Error(t, executor.Execute(trx, tSandbox))
		assert.NoError(t, checker.Execute(trx, tSandbox))
	})
}
