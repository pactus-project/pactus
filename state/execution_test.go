package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposeBlock(t *testing.T) {
	td := setup(t)

	lockTime := td.state.LastBlockHeight()
	dupSubsidyTx := td.GenerateTestSubsidyTx(testsuite.TransactionWithLockTime(lockTime))
	invTransferTx := td.GenerateTestTransferTx()
	invBondTx := td.GenerateTestBondTx()
	invSortitionTx := td.GenerateTestSortitionTx()

	validTrx1 := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lockTime),
		testsuite.TransactionWithSigner(td.genAccKey))

	validTrx2 := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(lockTime),
		testsuite.TransactionWithSigner(td.genAccKey))

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(invBondTx))
	assert.NoError(t, td.state.AddPendingTx(invSortitionTx))
	assert.NoError(t, td.state.AddPendingTx(dupSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTrx1))
	assert.NoError(t, td.state.AddPendingTx(validTrx2))

	rewardAddr := td.RandAccAddress()
	blk, err := td.state.ProposeBlock(td.state.valKeys[0], rewardAddr)
	require.NoError(t, err)

	blockTrxs := blk.Transactions()
	rewardTrx := blockTrxs[0]

	assert.Equal(t, protocol.ProtocolVersion2, blk.Header().Version())
	assert.Equal(t, td.state.valKeys[0].Address(), blk.Header().ProposerAddress())
	assert.Equal(t, td.state.LastBlockHash(), blk.Header().PrevBlockHash())
	assert.Equal(t, block.Txs{rewardTrx, validTrx1, validTrx2}, blockTrxs)
	assert.Equal(t, td.state.params.BlockReward+validTrx1.Fee()+validTrx2.Fee(), rewardTrx.Payload().Value())
}

func TestExecuteBlock(t *testing.T) {
	td := setup(t)

	blk, cert := td.makeBlockAndCertificate(t, 0)
	assert.NoError(t, td.state.CommitBlock(blk, cert))

	invTransferTx := td.GenerateTestTransferTx()
	validTx1 := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(1),
		testsuite.TransactionWithSigner(td.genAccKey))

	invSubsidyTx := td.state.createSubsidyTx(td.genValKeys[0].Address(), td.RandAccAddress(), validTx1.Fee()+1)
	validSubsidyTx := td.state.createSubsidyTx(td.genValKeys[0].Address(), td.RandAccAddress(), validTx1.Fee())

	assert.NoError(t, td.state.AddPendingTx(invTransferTx))
	assert.NoError(t, td.state.AddPendingTx(validSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(invSubsidyTx))
	assert.NoError(t, td.state.AddPendingTx(validTx1))

	t.Run("Block has invalid subsidy amount", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(invSubsidyTx)
		txs.Append(validTx1)
		invBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, InvalidSubsidyAmountError{
			Expected: 1e9 + validTx1.Fee(),
			Got:      1e9 + validTx1.Fee() + 1,
		})
	})

	t.Run("Block has an invalid transaction", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(invTransferTx)
		invBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, executor.AccountNotFoundError{
			Address: invTransferTx.Payload().Signer(),
		})
	})

	t.Run("Subsidy is not first transaction in block", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		txs.Append(validSubsidyTx)
		invBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Block has no subsidy transaction", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validTx1)
		invBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Block has two subsidy transactions", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validSubsidyTx)
		invBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		err := td.state.executeBlock(invBlock, sb, true)
		assert.ErrorIs(t, err, ErrDuplicatedSubsidyTransaction)
	})

	t.Run("OK", func(t *testing.T) {
		txs := block.NewTxs()
		txs.Append(validSubsidyTx)
		txs.Append(validTx1)
		validBlock, _ := td.GenerateTestBlock(1,
			testsuite.BlockWithStateHash(td.state.stateRoot()),
			testsuite.BlockWithPrevCert(td.state.lastInfo.Certificate()),
			testsuite.BlockWithPrevHash(td.state.lastInfo.BlockHash()),
			testsuite.BlockWithSeed(td.state.lastInfo.SortitionSeed()),
			testsuite.BlockWithTransactions(txs))

		sb := td.state.concreteSandbox()
		assert.NoError(t, td.state.executeBlock(validBlock, sb, true))

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		assert.Equal(t, 21*1e15-(10*td.state.params.BlockReward), treasury.Balance()) // Two extra blocks has committed yet
	})
}

func TestSubsidyTransaction(t *testing.T) {
	td := setup(t)

	t.Run("Legacy Reward", func(t *testing.T) {
		trx := tx.NewTransferTx(td.RandHeight(), crypto.TreasuryAddress, td.RandAccAddress(), td.RandAmount(), 0)

		err := td.state.checkSubsidy(trx, td.RandValAddress(), true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward With No Foundation Address", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		}
		trx := td.GenerateTestSubsidyTx(testsuite.TransactionWithRecipients(recipients))

		err := td.state.checkSubsidy(trx, td.RandValAddress(), true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward With Invalid Foundation Address", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{
				To:     td.RandAccAddress(),
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		}
		trx := td.GenerateTestSubsidyTx(testsuite.TransactionWithRecipients(recipients))

		err := td.state.checkSubsidy(trx, td.RandValAddress(), true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Split Reward: Ok", func(t *testing.T) {
		lockTime := td.RandHeight()
		recipients := []payload.BatchRecipient{
			{
				To:     td.state.params.FoundationAddress[lockTime%100],
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
		}
		trx := td.GenerateTestSubsidyTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithRecipients(recipients))

		err := td.state.checkSubsidy(trx, td.RandValAddress(), true)
		assert.NoError(t, err)
	})

	t.Run("Delegated proposer accepts valid 3-recipient subsidy", func(t *testing.T) {
		td.state.params.BlockVersion = protocol.ProtocolVersion3
		proposerAddr := td.genValKeys[0].Address()
		delegateOwner := td.RandAccAddress()
		delegateShare := amount.Amount(2e8)
		lockTime := td.RandHeight()

		val, err := td.state.store.Validator(proposerAddr)
		require.NoError(t, err)
		val.SetDelegation(delegateOwner, delegateShare, lockTime+10)
		td.state.store.UpdateValidator(val)

		recipients := []payload.BatchRecipient{
			{
				To:     td.state.params.FoundationAddress[lockTime%100],
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.state.params.BlockReward - td.state.params.FoundationReward - delegateShare,
			},
			{
				To:     delegateOwner,
				Amount: delegateShare,
			},
		}
		trx := td.GenerateTestSubsidyTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithRecipients(recipients))

		err = td.state.checkSubsidy(trx, proposerAddr, true)
		assert.NoError(t, err)
	})

	t.Run("Delegated proposer rejects invalid owner amount/address in 3-recipient subsidy", func(t *testing.T) {
		td.state.params.BlockVersion = protocol.ProtocolVersion3
		proposerAddr := td.genValKeys[1].Address()
		delegateOwner := td.RandAccAddress()
		delegateShare := amount.Amount(3e8)
		lockTime := td.RandHeight()

		val, err := td.state.store.Validator(proposerAddr)
		require.NoError(t, err)
		val.SetDelegation(delegateOwner, delegateShare, lockTime+10)
		td.state.store.UpdateValidator(val)

		badRecipients := []payload.BatchRecipient{
			{
				To:     td.state.params.FoundationAddress[lockTime%100],
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.state.params.BlockReward - td.state.params.FoundationReward - delegateShare,
			},
			{
				To:     td.RandAccAddress(),
				Amount: delegateShare + 1,
			},
		}
		trx := td.GenerateTestSubsidyTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithRecipients(badRecipients))

		err = td.state.checkSubsidy(trx, proposerAddr, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Non-delegated proposer rejects 3-recipient subsidy", func(t *testing.T) {
		td.state.params.BlockVersion = protocol.ProtocolVersion3
		lockTime := td.RandHeight()
		proposerAddr := td.genValKeys[2].Address()

		val, err := td.state.store.Validator(proposerAddr)
		require.NoError(t, err)
		val.SetDelegation(crypto.TreasuryAddress, 0, 0)
		td.state.store.UpdateValidator(val)

		recipients := []payload.BatchRecipient{
			{
				To:     td.state.params.FoundationAddress[lockTime%100],
				Amount: td.state.params.FoundationReward,
			},
			{
				To:     td.RandAccAddress(),
				Amount: td.RandAmount(),
			},
			{
				To:     td.RandAccAddress(),
				Amount: 0,
			},
		}
		trx := td.GenerateTestSubsidyTx(
			testsuite.TransactionWithLockTime(lockTime),
			testsuite.TransactionWithRecipients(recipients))

		err = td.state.checkSubsidy(trx, proposerAddr, true)
		assert.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
	})

	t.Run("Delegated proposer share boundaries allow 2-recipient subsidy", func(t *testing.T) {
		td.state.params.BlockVersion = protocol.ProtocolVersion3
		proposerAddr := td.genValKeys[3].Address()
		lockTime := td.RandHeight()

		val, err := td.state.store.Validator(proposerAddr)
		require.NoError(t, err)

		for _, share := range []amount.Amount{0, param.MaxDelegateOwnerRewardShare} {
			val.SetDelegation(td.RandAccAddress(), share, lockTime+10)
			td.state.store.UpdateValidator(val)

			recipients := []payload.BatchRecipient{
				{
					To:     td.state.params.FoundationAddress[lockTime%100],
					Amount: td.state.params.FoundationReward,
				},
				{
					To:     td.RandAccAddress(),
					Amount: td.RandAmount(),
				},
			}
			trx := td.GenerateTestSubsidyTx(
				testsuite.TransactionWithLockTime(lockTime),
				testsuite.TransactionWithRecipients(recipients))

			err := td.state.checkSubsidy(trx, proposerAddr, true)
			assert.NoError(t, err)
		}
	})
}
