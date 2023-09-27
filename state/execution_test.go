package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestProposeBlock(t *testing.T) {
	td := setup(t)

	curHeight := uint32(7)
	for i := uint32(0); i < curHeight; i++ {
		td.moveToNextHeightForAllStates(t)
	}
	b1, c1 := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3)
	assert.NoError(t, td.state1.CommitBlock(curHeight+1, b1, c1))
	assert.NoError(t, td.state2.CommitBlock(curHeight+1, b1, c1))
	assert.Equal(t, td.state1.LastBlockHeight(), curHeight+1)

	invSubsidyTx := tx.NewSubsidyTx(1, td.valKey2.Address(),
		td.state1.params.BlockReward, "duplicated subsidy transaction")
	invTransferTx, _ := td.GenerateTestTransferTx()
	invBondTx, _ := td.GenerateTestBondTx()
	invSortitionTx, _ := td.GenerateTestSortitionTx()

	pub, _ := td.RandBLSKeyPair()
	trx1 := tx.NewTransferTx(1, td.valKey1.Address(), td.valKey1.Address(), 1, 1000, "")
	td.HelperSignTransaction(td.valKey1.PrivateKey(), trx1)

	trx2 := tx.NewBondTx(2, td.valKey1.Address(), pub.ValidatorAddress(), pub, 1000000000, 100000, "")
	td.HelperSignTransaction(td.valKey1.PrivateKey(), trx2)

	assert.NoError(t, td.state1.txPool.AppendTx(invTransferTx))
	assert.NoError(t, td.state1.txPool.AppendTx(invBondTx))
	assert.NoError(t, td.state1.txPool.AppendTx(invSortitionTx))
	assert.NoError(t, td.state1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, td.state1.txPool.AppendTx(trx1))
	assert.NoError(t, td.state1.txPool.AppendTx(trx2))

	b2, c2 := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3)
	assert.Equal(t, b2.Header().PrevBlockHash(), b1.Hash())
	assert.Equal(t, b2.Transactions()[1:], block.Txs{trx1, trx2})
	assert.True(t, b2.Transactions()[0].IsSubsidyTx())
	assert.NoError(t, td.state1.CommitBlock(curHeight+2, b2, c2))

	assert.Equal(t, td.state1.TotalPower(), int64(1000000004))
	assert.Equal(t, td.state1.committee.TotalPower(), int64(4))
}

func TestExecuteBlock(t *testing.T) {
	td := setup(t)

	b1, c1 := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3)
	assert.NoError(t, td.state1.CommitBlock(1, b1, c1))

	proposerAddr := td.RandAccAddress()
	rewardAddr := td.RandAccAddress()
	invSubsidyTx := td.state1.createSubsidyTx(rewardAddr, 1001)
	validSubsidyTx := td.state1.createSubsidyTx(rewardAddr, 1000)
	invTransferTx, _ := td.GenerateTestTransferTx()

	validTx1 := tx.NewTransferTx(1, td.valKey1.Address(), td.valKey1.Address(), 1, 1000, "")
	td.HelperSignTransaction(td.valKey1.PrivateKey(), validTx1)

	assert.NoError(t, td.state1.txPool.AppendTx(invTransferTx))
	assert.NoError(t, td.state1.txPool.AppendTx(validSubsidyTx))
	assert.NoError(t, td.state1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, td.state1.txPool.AppendTx(validTx1))

	t.Run("Subsidy tx is invalid", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(invSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()

		assert.Error(t, td.state1.executeBlock(invBlock, sb))
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(invTransferTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()

		assert.Error(t, td.state1.executeBlock(invBlock, sb))
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()

		assert.Error(t, td.state1.executeBlock(invBlock, sb))
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()

		assert.Error(t, td.state1.executeBlock(invBlock, sb))
	})

	t.Run("Two subsidy transactions", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()

		assert.Error(t, td.state1.executeBlock(invBlock, sb))
	})

	t.Run("OK", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(),
			td.state1.stateRoot(), td.state1.lastInfo.Certificate(),
			td.state1.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state1.concreteSandbox()
		assert.NoError(t, td.state1.executeBlock(invBlock, sb))

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		subsidy := td.state1.params.BlockReward
		assert.Equal(t, treasury.Balance(), 21*1e15-(2*subsidy)) // Two blocks has committed yet
	})
}
