package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestProposeBlock(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.NoError(t, tState2.CommitBlock(1, b1, c1))

	invSubsidyTx := tx.NewSubsidyTx(tState1.lastInfo.BlockHash().Stamp(), 1, tValSigner2.Address(),
		tState1.params.BlockReward, "duplicated subsidy transaction")
	invSendTx, _ := tx.GenerateTestSendTx()
	invBondTx, _ := tx.GenerateTestBondTx()
	invSortitionTx, _ := tx.GenerateTestSortitionTx()

	pub, _ := bls.GenerateTestKeyPair()
	trx1 := tx.NewSendTx(b1.Stamp(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "")
	tValSigner1.SignMsg(trx1)

	trx2 := tx.NewBondTx(b1.Stamp(), 2, tValSigner1.Address(), pub.Address(), pub, 1000, 1000, "")
	tValSigner1.SignMsg(trx2)

	assert.NoError(t, tState1.txPool.AppendTx(invSendTx))
	assert.NoError(t, tState1.txPool.AppendTx(invBondTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSortitionTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(trx1))
	assert.NoError(t, tState1.txPool.AppendTx(trx2))

	b2, c2 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.Equal(t, b2.Header().PrevBlockHash(), b1.Hash())
	assert.Equal(t, b2.Transactions()[1:], block.Txs{trx1, trx2})
	assert.True(t, b2.Transactions()[0].IsSubsidyTx())
	assert.NoError(t, tState1.CommitBlock(2, b2, c2))

	assert.Equal(t, tState1.TotalPower(), int64(1004))
	assert.Equal(t, tState1.committee.TotalPower(), int64(4))
}

func TestExecuteBlock(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))

	proposerAddr := crypto.GenerateTestAddress()
	rewardAddr := crypto.GenerateTestAddress()
	invSubsidyTx := tState1.createSubsidyTx(rewardAddr, 1001)
	validSubsidyTx := tState1.createSubsidyTx(rewardAddr, 1000)
	invSendTx, _ := tx.GenerateTestSendTx()

	validTx1 := tx.NewSendTx(b1.Stamp(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "")
	tValSigner1.SignMsg(validTx1)

	assert.NoError(t, tState1.txPool.AppendTx(invSendTx))
	assert.NoError(t, tState1.txPool.AppendTx(validSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(validTx1))

	t.Run("Subsidy tx is invalid", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(invSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.Error(t, tState1.executeBlock(invBlock, sb))
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(invSendTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.Error(t, tState1.executeBlock(invBlock, sb))
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.Error(t, tState1.executeBlock(invBlock, sb))
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.Error(t, tState1.executeBlock(invBlock, sb))
	})

	t.Run("Two subsidy transactions", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.Error(t, tState1.executeBlock(invBlock, sb))
	})

	t.Run("OK", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)
		sb := tState1.concreteSandbox()
		assert.NoError(t, tState1.executeBlock(invBlock, sb))

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		subsidy := tState1.params.BlockReward
		assert.Equal(t, treasury.Balance(), 21*1e14-(2*subsidy)) // Two blocks has committed yet
	})
}
