package state

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestProposeBlock(t *testing.T) {
	td := setup(t)

	proposer := td.state.Proposer(0)
	lockTime := td.state.LastBlockHeight()
	dupSubsidyTx := tx.NewSubsidyTxLegacy(lockTime, proposer.Address(), td.state.params.BlockReward)
	invTransferTx := td.GenerateTestTransferTx()
	invBondTx := td.GenerateTestBondTx()
	invSortitionTx := td.GenerateTestSortitionTx()

	validTrx1 := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lockTime),
		testsuite.TransactionWithEd25519Signer(td.genAccKey))

	validTrx2 := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lockTime),
		testsuite.TransactionWithEd25519Signer(td.genAccKey))

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(invBondTx))
	assert.NoError(t, td.state.AddPendingTx(invSortitionTx))
	assert.NoError(t, td.state.AddPendingTx(dupSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTrx1))
	assert.NoError(t, td.state.AddPendingTx(validTrx2))

	blk, cert := td.makeBlockAndCertificate(t, 0)
	assert.Equal(t, td.state.LastBlockHash(), blk.Header().PrevBlockHash())
	assert.Equal(t, block.Txs{validTrx1, validTrx2}, blk.Transactions()[1:])
	assert.True(t, blk.Transactions()[0].IsSubsidyTx())
	assert.NoError(t, td.state.CommitBlock(blk, cert))
}

func TestExecuteBlock(t *testing.T) {
	td := setup(t)

	blk, cert := td.makeBlockAndCertificate(t, 0)
	assert.NoError(t, td.state.CommitBlock(blk, cert))

	proposerAddr := td.RandAccAddress()
	rewardAddr := td.RandAccAddress()
	invSubsidyTx, _ := td.state.createSubsidyTx(rewardAddr, 1001)
	validSubsidyTx, _ := td.state.createSubsidyTx(rewardAddr, 1000)
	invTransferTx := td.GenerateTestTransferTx()

	validTx1 := tx.NewTransferTx(1, td.genAccKey.PublicKeyNative().AccountAddress(),
		td.RandAccAddress(), 1, 1000)
	td.HelperSignTransaction(td.genAccKey, validTx1)

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(validSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(invSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTx1))

	t.Run("Subsidy amount is invalid", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(invSubsidyTx)
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, InvalidSubsidyAmountError{
			Expected: amount.Amount(1e9 + 1000),
			Got:      amount.Amount(1e9 + 1001),
		})
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(invTransferTx)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{
			Address: invTransferTx.Payload().Signer(),
		})
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Two subsidy transactions", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validSubsidyTx)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrDuplicatedSubsidyTransaction)
	})

	t.Run("OK", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validTx1)
		invBlock := block.MakeBlock(1, time.Now(), txs, td.state.lastInfo.BlockHash(),
			td.state.stateRoot(), td.state.lastInfo.Certificate(),
			td.state.lastInfo.SortitionSeed(), proposerAddr)
		sb := td.state.concreteSandbox()
		assert.NoError(t, td.state.executeBlock(invBlock, sb, true))

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		subsidy := td.state.params.BlockReward
		assert.Equal(t, 21*1e15-(10*subsidy), treasury.Balance()) // Two extra blocks has committed yet
	})
}

func TestSubsidyTransaction(t *testing.T) {
	td := setup(t)

	t.Run("Legacy Reward", func(t *testing.T) {
		td.state.params.SplitRewardForkHeight = 0
		trx := tx.NewSubsidyTxLegacy(td.RandHeight(), td.RandAccAddress(), td.RandAmount())

		err := td.state.checkSubsidy(trx, true)
		assert.NoError(t, err)
	})

	t.Run("Legacy Reward, Invalid transaction", func(t *testing.T) {
		td.state.params.SplitRewardForkHeight = 0
		trx := tx.NewSubsidyTx(td.RandHeight(), []payload.BatchRecipient{
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		})

		err := td.state.checkSubsidy(trx, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Legacy Reward", func(t *testing.T) {
		splitHeight := td.RandHeight()
		td.state.params.SplitRewardForkHeight = splitHeight
		trx := tx.NewSubsidyTxLegacy(splitHeight-1, td.RandAccAddress(), td.RandAmount())

		err := td.state.checkSubsidy(trx, true)
		assert.NoError(t, err)
	})

	t.Run("Legacy Reward after splitting Reward", func(t *testing.T) {
		splitHeight := td.RandHeight()
		td.state.isSplitForkEnabled = true
		td.state.params.SplitRewardForkHeight = splitHeight
		trx := tx.NewSubsidyTxLegacy(splitHeight+1, td.RandAccAddress(), td.RandAmount())

		err := td.state.checkSubsidy(trx, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward With Invalid Foundation", func(t *testing.T) {
		splitHeight := td.RandHeight()
		td.state.params.SplitRewardForkHeight = splitHeight
		td.state.params.FoundationAddress = []crypto.Address{td.RandAccAddress()}

		trx := tx.NewSubsidyTx(splitHeight+1, []payload.BatchRecipient{
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		})

		err := td.state.checkSubsidy(trx, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward With Invalid Foundation", func(t *testing.T) {
		splitHeight := td.RandHeight()
		td.state.params.SplitRewardForkHeight = splitHeight
		td.state.params.FoundationAddress = []crypto.Address{td.RandAccAddress()}

		trx := tx.NewSubsidyTx(splitHeight+1, []payload.BatchRecipient{
			{
				To:     td.RandAccAddress(),
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		})

		err := td.state.checkSubsidy(trx, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward: Ok", func(t *testing.T) {
		splitHeight := td.RandHeight()
		td.state.params.SplitRewardForkHeight = splitHeight
		td.state.params.FoundationAddress = []crypto.Address{td.RandAccAddress()}

		trx := tx.NewSubsidyTx(splitHeight+1, []payload.BatchRecipient{
			{
				To:     td.state.params.FoundationAddress[0],
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		})

		err := td.state.checkSubsidy(trx, true)
		assert.NoError(t, err)
	})
}
