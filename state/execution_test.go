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

	proposer := td.state.Proposer(0)
	lockTime := td.state.LastBlockHeight()
	dupSubsidyTx := tx.NewSubsidyTx(lockTime, proposer.Address(),
		td.state.params.BlockReward, "duplicated subsidy transaction")
	invTransferTx, _ := td.GenerateTestTransferTx()
	invBondTx, _ := td.GenerateTestBondTx()
	invSortitionTx, _ := td.GenerateTestSortitionTx()

	pub, _ := td.RandBLSKeyPair()
	validTrx1 := tx.NewTransferTx(lockTime, td.genAccKey.PublicKeyNative().AccountAddress(),
		td.RandAccAddress(), 1, 1000, "")
	td.HelperSignTransaction(td.genAccKey, validTrx1)

	validTrx2 := tx.NewBondTx(lockTime, td.genAccKey.PublicKeyNative().AccountAddress(),
		pub.ValidatorAddress(), pub, 1000000000, 100000, "")
	td.HelperSignTransaction(td.genAccKey, validTrx2)

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(invBondTx))
	assert.NoError(t, td.state.AddPendingTx(invSortitionTx))
	assert.NoError(t, td.state.AddPendingTx(dupSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTrx1))
	assert.NoError(t, td.state.AddPendingTx(validTrx2))

	blk, cert := td.makeBlockAndCertificate(t, 0)
	assert.Equal(t, blk.Header().PrevBlockHash(), td.state.LastBlockHash())
	assert.Equal(t, blk.Transactions()[1:], block.Txs{validTrx1, validTrx2})
	assert.True(t, blk.Transactions()[0].IsSubsidyTx())
	assert.NoError(t, td.state.CommitBlock(blk, cert))

	assert.Equal(t, td.state.TotalPower(), int64(1000000004))
	assert.Equal(t, td.state.committee.TotalPower(), int64(4))
}

func TestExecuteBlock(t *testing.T) {
	td := setup(t)

	blk, cert := td.makeBlockAndCertificate(t, 0)
	assert.NoError(t, td.state.CommitBlock(blk, cert))

	proposerAddr := td.RandAccAddress()
	rewardAddr := td.RandAccAddress()
	invSubsidyTx := td.state.createSubsidyTx(rewardAddr, 1001)
	validSubsidyTx := td.state.createSubsidyTx(rewardAddr, 1000)
	invTransferTx, _ := td.GenerateTestTransferTx()

	validTx1 := tx.NewTransferTx(1, td.genAccKey.PublicKeyNative().AccountAddress(),
		td.RandAccAddress(), 1, 1000, "")
	td.HelperSignTransaction(td.genAccKey, validTx1)

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(validSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(invSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTx1))

	t.Run("Subsidy tx is invalid", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(invSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()

		assert.Error(t, td.state.executeBlock(invBlock, sb))
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(invTransferTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()

		assert.Error(t, td.state.executeBlock(invBlock, sb))
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()

		assert.Error(t, td.state.executeBlock(invBlock, sb))
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()

		assert.Error(t, td.state.executeBlock(invBlock, sb))
	})

	t.Run("Two subsidy transactions", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()

		assert.Error(t, td.state.executeBlock(invBlock, sb))
	})

	t.Run("OK", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, util.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		assert.NoError(t, td.state.executeBlock(invBlock, sb))

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		subsidy := td.state.params.BlockReward
		assert.Equal(t, treasury.Balance(), 21*1e15-(12*subsidy)) // Two blocks has committed yet
	})
}
