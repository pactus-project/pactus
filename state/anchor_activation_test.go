package state

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAnchorProposeThreshold(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion4)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	require.Len(t, td.state.committee.Validators(), 4)

	t.Run("none", func(t *testing.T) {
		td.requireProposed(t, protocol.ProtocolVersion4)
	})

	t.Run("exactly 75 percent", func(t *testing.T) {
		applyCommittee(t, td.state, []protocol.Version{
			protocol.ProtocolVersion5, protocol.ProtocolVersion5,
			protocol.ProtocolVersion5, protocol.ProtocolVersion4,
		})
		require.False(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))
		td.requireProposed(t, protocol.ProtocolVersion4)
	})

	t.Run("everyone", func(t *testing.T) {
		applyCommittee(t, td.state, []protocol.Version{
			protocol.ProtocolVersion5, protocol.ProtocolVersion5,
			protocol.ProtocolVersion5, protocol.ProtocolVersion5,
		})
		td.requireProposed(t, protocol.ProtocolVersion5)
		require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	})

	t.Run("version 6 counts as 5", func(t *testing.T) {
		applyCommittee(t, td.state, []protocol.Version{6, 6, 6, 6})
		td.requireProposed(t, protocol.ProtocolVersion5)
	})

	t.Run("version 4 does not", func(t *testing.T) {
		applyCommittee(t, td.state, []protocol.Version{
			protocol.ProtocolVersion4, protocol.ProtocolVersion4,
			protocol.ProtocolVersion4, protocol.ProtocolVersion4,
		})
		require.False(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))
		td.requireProposed(t, protocol.ProtocolVersion4)
	})

	t.Run("no downgrade after a version 5 block", func(t *testing.T) {
		applyCommittee(t, td.state, []protocol.Version{
			protocol.ProtocolVersion5, protocol.ProtocolVersion5,
			protocol.ProtocolVersion5, protocol.ProtocolVersion5,
		})
		blk := td.proposeEmpty(t)
		require.Equal(t, protocol.ProtocolVersion5, blk.Header().Version())
		require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
		td.commit(t, blk)
		require.Equal(t, protocol.ProtocolVersion5, td.state.params.BlockVersion)

		applyCommittee(t, td.state, []protocol.Version{
			protocol.ProtocolVersion4, protocol.ProtocolVersion4,
			protocol.ProtocolVersion4, protocol.ProtocolVersion4,
		})
		td.requireProposed(t, protocol.ProtocolVersion5)
	})
}

func TestAnchorPowerThreshold(t *testing.T) {
	versions := []protocol.Version{
		protocol.ProtocolVersion5, protocol.ProtocolVersion4,
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
	}

	t.Run("one validator at 75 percent", func(t *testing.T) {
		// Stake 9 against three bootstrap validators (power 1) is 9/12.
		td := setupWithVersion(t, protocol.ProtocolVersion4)
		replaceCommittee(t, td.state, []amount.Amount{9, 0, 0, 0}, versions)
		require.False(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))
		td.requireProposed(t, protocol.ProtocolVersion4)
	})

	t.Run("one validator above 75 percent", func(t *testing.T) {
		td := setupWithVersion(t, protocol.ProtocolVersion4)
		replaceCommittee(t, td.state, []amount.Amount{10, 0, 0, 0}, versions)
		require.True(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))
		td.requireProposed(t, protocol.ProtocolVersion5)
	})
}

func TestHandBuiltVersion5BelowThreshold(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion4)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
	})

	sender := td.sender()
	balance := td.mustAccount(t, sender).Balance()
	template := td.proposeEmpty(t)
	require.Equal(t, protocol.ProtocolVersion4, template.Header().Version())
	require.False(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))

	root := bytes.Repeat([]byte{0x44}, 32)
	anchor := td.anchor(template.Height(), payload.AnchorActionSet, root, "manual",
		executor.MinAnchorDeposit, 0)
	blk := rebuildBlock(template, protocol.ProtocolVersion5, appendTxs(template.Transactions(), anchor))
	td.commit(t, blk)

	require.Equal(t, protocol.ProtocolVersion5, td.state.params.BlockVersion)
	acc := td.mustAccount(t, sender)
	require.True(t, acc.HasAnchor())
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())
	require.Equal(t, balance-executor.MinAnchorDeposit, acc.Balance())
	require.Equal(t, root, acc.RootHash())
	require.False(t, td.state.committee.SupportProtocolVersion(protocol.ProtocolVersion5))
	require.Equal(t, protocol.ProtocolVersion5, td.state.proposeBlockVersion())
}

func TestAnchorVersion5DoesNotWriteAnchor(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion4)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	sender := td.sender()
	third := td.RandAccAddress()

	credit := tx.NewTransferTx(td.nextHeight(), sender, third, 5, 0)
	td.HelperSignTransaction(td.genAccKey, credit)
	td.commit(t, td.propose(t, block.Txs{credit}))

	senderHash := td.mustAccount(t, sender).Hash()
	thirdHash := td.mustAccount(t, third).Hash()
	valHashes := td.validatorHashes(t)

	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
	})

	blk := td.proposeEmpty(t)
	require.Equal(t, protocol.ProtocolVersion5, blk.Header().Version())
	require.False(t, blockHasAnchor(blk))
	require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	td.commit(t, blk)

	require.Equal(t, protocol.ProtocolVersion5, td.state.params.BlockVersion)
	require.Len(t, td.accountBytes(t, sender), 12)
	require.Len(t, td.accountBytes(t, third), 12)
	require.False(t, td.mustAccount(t, sender).HasAnchor())
	require.Equal(t, senderHash, td.mustAccount(t, sender).Hash())
	require.Equal(t, thirdHash, td.mustAccount(t, third).Hash())
	require.Equal(t, valHashes, td.validatorHashes(t))
}

func TestAnchorActivation(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion4)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	sender := td.sender()
	third := td.RandAccAddress()

	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
	})

	credit := tx.NewTransferTx(td.nextHeight(), sender, third, 5, 0)
	td.HelperSignTransaction(td.genAccKey, credit)
	td.commit(t, td.propose(t, block.Txs{credit}))
	require.Len(t, td.accountBytes(t, sender), 12)
	require.Equal(t, amount.Amount(5), td.mustAccount(t, third).Balance())

	root := bytes.Repeat([]byte{0x11}, 32)
	early := td.anchor(td.nextHeight(), payload.AnchorActionSet, root, "below", executor.MinAnchorDeposit, 1)
	senderHash := td.mustAccount(t, sender).Hash()
	earlyBlk := td.propose(t, block.Txs{early})
	require.Equal(t, protocol.ProtocolVersion4, earlyBlk.Header().Version())
	require.False(t, blockHasAnchor(earlyBlk))
	td.commit(t, earlyBlk)
	require.Equal(t, senderHash, td.mustAccount(t, sender).Hash())
	require.Len(t, td.accountBytes(t, sender), 12)
	require.False(t, td.mustAccount(t, sender).HasAnchor())

	template := td.proposeEmpty(t)
	hand := td.anchor(template.Height(), payload.AnchorActionSet, root, "hand", executor.MinAnchorDeposit, 0)
	handBlk := rebuildBlock(template, protocol.ProtocolVersion5, appendTxs(template.Transactions(), hand))
	require.NoError(t, td.state.ValidateBlock(handBlk, 0))
	require.False(t, td.mustAccount(t, sender).HasAnchor())
	require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	require.Equal(t, protocol.ProtocolVersion4, td.state.proposeBlockVersion())

	tooNew := rebuildBlock(template, protocol.ProtocolVersionLatest+1, template.Transactions())
	require.ErrorIs(t, td.state.ValidateBlock(tooNew, 0), InvalidBlockVersionError{
		Version: protocol.ProtocolVersionLatest + 1,
	})
	tooOld := rebuildBlock(template, protocol.ProtocolVersion3, template.Transactions())
	require.ErrorIs(t, td.state.ValidateBlock(tooOld, 0), InvalidBlockVersionError{
		Version: protocol.ProtocolVersion3,
	})

	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
	})

	badTemplate := td.proposeEmpty(t)
	require.Equal(t, protocol.ProtocolVersion5, badTemplate.Header().Version())
	require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	plain := td.accountBytes(t, sender)
	stateRoot := td.state.stateRoot()
	balance := td.mustAccount(t, sender).Balance()

	rejected := []struct {
		name    string
		trx     *tx.Tx
		wantErr error
	}{
		{
			name: "deposit below 1 PAC",
			trx: td.anchor(badTemplate.Height(), payload.AnchorActionSet, root, "small",
				executor.MinAnchorDeposit-1, 0),
			wantErr: executor.ErrAnchorDepositTooSmall,
		},
		{
			name: "balance too small",
			trx: td.anchor(badTemplate.Height(), payload.AnchorActionSet, root, "drain",
				balance, 1),
			wantErr: executor.ErrInsufficientFunds,
		},
		{
			name:    "delete without an anchor",
			trx:     td.anchor(badTemplate.Height(), payload.AnchorActionDelete, nil, "", 0, 0),
			wantErr: executor.ErrAnchorNotFound,
		},
		{
			name: "fee overflows the ceiling",
			trx: td.anchor(badTemplate.Height(), payload.AnchorActionSet, root, "ceil",
				1, amount.Amount(amount.MaxNanoPAC)),
			wantErr: executor.ErrAmountOverflow,
		},
	}
	for _, item := range rejected {
		blk := rebuildBlock(badTemplate, protocol.ProtocolVersion5, appendTxs(badTemplate.Transactions(), item.trx))
		require.ErrorIs(t, td.state.ValidateBlock(blk, 0), item.wantErr, item.name)
		require.Equal(t, plain, td.accountBytes(t, sender), item.name)
		require.Equal(t, stateRoot, td.state.stateRoot(), item.name)
		require.Len(t, td.accountBytes(t, sender), 12, item.name)
	}

	// The subsidy is raised by the fee of an invalid anchor. Validation must
	// reject the block, so the fee never reaches the treasury.
	fee := amount.Amount(7)
	padded := td.anchor(badTemplate.Height(), payload.AnchorActionSet, root, "padded",
		executor.MinAnchorDeposit-1, fee)
	treasury := td.mustAccount(t, crypto.TreasuryAddress).Balance()
	paddedBlk := rebuildBlock(badTemplate, protocol.ProtocolVersion5, block.Txs{
		subsidyPlus(badTemplate.Transactions().Subsidy(), fee), padded,
	})
	require.ErrorIs(t, td.state.ValidateBlock(paddedBlk, 0), executor.ErrAnchorDepositTooSmall)
	require.Equal(t, treasury, td.mustAccount(t, crypto.TreasuryAddress).Balance())
	require.Equal(t, balance, td.mustAccount(t, sender).Balance())
	require.Len(t, td.accountBytes(t, sender), 12)

	twice := td.anchor(badTemplate.Height(), payload.AnchorActionSet, root, "twice",
		executor.MinAnchorDeposit, 0)
	dupBlk := rebuildBlock(badTemplate, protocol.ProtocolVersion5,
		appendTxs(badTemplate.Transactions(), twice, twice))
	require.ErrorIs(t, td.state.ValidateBlock(dupBlk, 0), execution.TransactionCommittedError{ID: twice.ID()})
	require.False(t, td.mustAccount(t, sender).HasAnchor())
	require.Equal(t, plain, td.accountBytes(t, sender))
	require.Equal(t, stateRoot, td.state.stateRoot())

	height := td.nextHeight()
	rootA := bytes.Repeat([]byte{0x21}, 32)
	rootB := bytes.Repeat([]byte{0x22}, 64)
	set1 := td.anchor(height, payload.AnchorActionSet, rootA, "first", executor.MinAnchorDeposit, 1)
	set2 := td.anchor(height, payload.AnchorActionSet, rootB, "second", 1, 1)
	startBalance := td.mustAccount(t, sender).Balance()
	startRoot := td.state.stateRoot()
	startHash := td.mustAccount(t, sender).Hash()
	thirdHash := td.mustAccount(t, third).Hash()
	valHashes := td.validatorHashes(t)

	good := td.propose(t, block.Txs{set1, set2})
	require.Equal(t, protocol.ProtocolVersion5, good.Header().Version())
	require.Equal(t, protocol.ProtocolVersion4, td.state.params.BlockVersion)
	require.True(t, blockHasTx(good, set1.ID()))
	require.True(t, blockHasTx(good, set2.ID()))
	td.checkBlockSubsidy(t, good)

	downgrade := rebuildBlock(good, protocol.ProtocolVersion4, good.Transactions())
	require.ErrorIs(t, td.state.ValidateBlock(downgrade, 0), executor.InvalidPayloadTypeError{
		PayloadType: payload.TypeAnchor,
	})
	require.Equal(t, startRoot, td.state.stateRoot())
	require.Equal(t, startHash, td.mustAccount(t, sender).Hash())
	require.Len(t, td.accountBytes(t, sender), 12)

	td.commit(t, good)
	require.Equal(t, protocol.ProtocolVersion5, td.state.params.BlockVersion)
	require.NotEqual(t, startRoot, td.state.stateRoot())

	anchored := td.mustAccount(t, sender)
	require.True(t, anchored.HasAnchor())
	require.Equal(t, executor.MinAnchorDeposit+1, anchored.LockedDeposit())
	require.Equal(t, startBalance-(executor.MinAnchorDeposit+1)-2, anchored.Balance())
	require.Equal(t, rootB, anchored.RootHash())
	require.Equal(t, "second", anchored.ManifestURI())
	require.NotEqual(t, startHash, anchored.Hash())
	require.Equal(t, thirdHash, td.mustAccount(t, third).Hash())
	require.Equal(t, valHashes, td.validatorHashes(t))
	require.Len(t, td.accountBytes(t, third), 12)

	leaf := anchored.Hash()
	set2.Payload().(*payload.AnchorPayload).RootHash[0] ^= 0xff
	require.Equal(t, leaf, td.mustAccount(t, sender).Hash())
	require.Equal(t, rootB, td.mustAccount(t, sender).RootHash())

	lock := td.mustAccount(t, sender).LockedDeposit()
	balance = td.mustAccount(t, sender).Balance()
	replayRoot := td.state.stateRoot()
	next := td.propose(t, block.Txs{set2})
	require.False(t, blockHasTx(next, set2.ID()))
	forced := rebuildBlock(next, next.Header().Version(), appendTxs(next.Transactions(), set2))
	require.ErrorIs(t, td.state.ValidateBlock(forced, 0), execution.TransactionCommittedError{ID: set2.ID()})
	require.Equal(t, replayRoot, td.state.stateRoot())
	require.Equal(t, lock, td.mustAccount(t, sender).LockedDeposit())
	require.Equal(t, balance, td.mustAccount(t, sender).Balance())
	td.commit(t, next)
	require.Equal(t, lock, td.mustAccount(t, sender).LockedDeposit())
	require.Equal(t, balance, td.mustAccount(t, sender).Balance())

	feeFree := td.anchor(td.nextHeight(), payload.AnchorActionSet, rootB, "second", 0, 0)
	feeBlk := td.propose(t, block.Txs{feeFree})
	require.True(t, blockHasTx(feeBlk, feeFree.ID()))
	td.checkBlockSubsidy(t, feeBlk)
	td.commit(t, feeBlk)
	require.Equal(t, lock, td.mustAccount(t, sender).LockedDeposit())
	require.Equal(t, balance, td.mustAccount(t, sender).Balance())

	deleteRoot := td.state.stateRoot()
	deleteBalance := td.mustAccount(t, sender).Balance()
	deleteLock := td.mustAccount(t, sender).LockedDeposit()
	height = td.nextHeight()
	thirdHash = td.mustAccount(t, third).Hash()
	open := td.anchor(height, payload.AnchorActionSet, rootA, "temp", executor.MinAnchorDeposit, 1)
	shut := td.anchor(height, payload.AnchorActionDelete, nil, "", 0, 1)
	deleted := td.propose(t, block.Txs{open, shut})
	require.True(t, blockHasTx(deleted, open.ID()))
	require.True(t, blockHasTx(deleted, shut.ID()))
	td.commit(t, deleted)
	require.False(t, td.mustAccount(t, sender).HasAnchor())
	require.Equal(t, amount.Amount(0), td.mustAccount(t, sender).LockedDeposit())
	require.Equal(t, deleteBalance+deleteLock-2, td.mustAccount(t, sender).Balance())
	require.Len(t, td.accountBytes(t, sender), 12)
	require.NotEqual(t, deleteRoot, td.state.stateRoot())
	require.Equal(t, thirdHash, td.mustAccount(t, third).Hash())

	spendBalance := td.mustAccount(t, sender).Balance()
	spendRoot := td.state.stateRoot()
	spendBytes := td.accountBytes(t, sender)
	overTemplate := td.proposeEmpty(t)
	overSet := td.anchor(overTemplate.Height(), payload.AnchorActionSet, rootA, "lock", executor.MinAnchorDeposit, 0)
	overPay := tx.NewTransferTx(overTemplate.Height(), sender, td.RandAccAddress(), spendBalance, 0)
	td.HelperSignTransaction(td.genAccKey, overPay)
	overBlk := rebuildBlock(overTemplate, protocol.ProtocolVersion5,
		appendTxs(overTemplate.Transactions(), overSet, overPay))
	require.ErrorIs(t, td.state.ValidateBlock(overBlk, 0), executor.ErrInsufficientFunds)
	require.Equal(t, spendRoot, td.state.stateRoot())
	require.Equal(t, spendBytes, td.accountBytes(t, sender))

	recipient := td.RandAccAddress()
	require.NotEqual(t, third, recipient)
	thirdHash = td.mustAccount(t, third).Hash()
	pay := tx.NewTransferTx(overTemplate.Height(), sender, recipient, spendBalance-executor.MinAnchorDeposit, 0)
	td.HelperSignTransaction(td.genAccKey, pay)
	moved := td.propose(t, block.Txs{overSet, pay})
	require.True(t, blockHasTx(moved, overSet.ID()))
	require.True(t, blockHasTx(moved, pay.ID()))
	td.commit(t, moved)
	require.Equal(t, executor.MinAnchorDeposit, td.mustAccount(t, sender).LockedDeposit())
	require.Equal(t, amount.Amount(0), td.mustAccount(t, sender).Balance())
	require.Equal(t, spendBalance-executor.MinAnchorDeposit, td.mustAccount(t, recipient).Balance())
	require.Equal(t, thirdHash, td.mustAccount(t, third).Hash())
	require.Len(t, td.accountBytes(t, third), 12)

	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
		protocol.ProtocolVersion4, protocol.ProtocolVersion4,
	})
	require.Equal(t, protocol.ProtocolVersion5, td.state.proposeBlockVersion())
	currentRoot := td.state.stateRoot()
	currentBytes := td.accountBytes(t, sender)
	plainBlk := td.proposeEmpty(t)
	require.Equal(t, protocol.ProtocolVersion5, plainBlk.Header().Version())
	old := rebuildBlock(plainBlk, protocol.ProtocolVersion4, plainBlk.Transactions())
	td.mustNotCommit(t, old, InvalidBlockVersionError{Version: protocol.ProtocolVersion4})
	again := td.anchor(plainBlk.Height(), payload.AnchorActionSet, rootA, "back", executor.MinAnchorDeposit, 0)
	oldAnchor := rebuildBlock(plainBlk, protocol.ProtocolVersion4, appendTxs(plainBlk.Transactions(), again))
	td.mustNotCommit(t, oldAnchor, InvalidBlockVersionError{Version: protocol.ProtocolVersion4})
	require.Equal(t, currentRoot, td.state.stateRoot())
	require.Equal(t, currentBytes, td.accountBytes(t, sender))
}

// bannedStore reports one address as banned, like a newer banned list would.
type bannedStore struct {
	store.Store

	banned crypto.Address
}

func (s bannedStore) IsBanned(addr crypto.Address) bool {
	return addr == s.banned
}

// A certified block must commit even if its anchor signer is banned later.
// Commit trusts the certificate and must not re-check anchors.
func TestAnchorCommitIgnoresLaterBan(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	sender := td.sender()

	root := bytes.Repeat([]byte{0x31}, 32)
	set := td.anchor(td.nextHeight(), payload.AnchorActionSet, root, "banned later", executor.MinAnchorDeposit, 1)
	blk := td.propose(t, block.Txs{set})
	require.True(t, blockHasTx(blk, set.ID()))
	require.NoError(t, td.state.ValidateBlock(blk, 0))

	td.state.store = bannedStore{Store: td.state.store, banned: sender}
	require.ErrorIs(t, td.state.ValidateBlock(blk, 0), execution.SignerBannedError{Address: sender})

	td.commit(t, blk)
	anchored := td.mustAccount(t, sender)
	require.True(t, anchored.HasAnchor())
	require.Equal(t, root, anchored.RootHash())
	require.Equal(t, executor.MinAnchorDeposit, anchored.LockedDeposit())
}

// Anchor clocks come from the header of the block that carries the
// transaction, and anchors never mint or burn coins (PIP-50 tests 20, 21, 26).
func TestAnchorLifecycleAcrossBlocks(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	sender := td.sender()
	supply := td.totalCoins(t)

	createBlk := td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x41}, 32), "create",
			executor.MinAnchorDeposit, 1),
	})
	td.commit(t, createBlk)
	created := td.mustAccount(t, sender)
	require.Equal(t, createBlk.Height(), created.CreatedAtHeight())
	require.Equal(t, createBlk.Header().UnixTime(), created.CreatedAtTime())
	require.Equal(t, createBlk.Height(), created.UpdatedAtHeight())
	require.Equal(t, createBlk.Header().UnixTime(), created.UpdatedAtTime())
	require.Equal(t, supply, td.totalCoins(t))

	rootB := bytes.Repeat([]byte{0x42}, 64)
	updateBlk := td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionSet, rootB, "update", 0, 1),
	})
	td.commit(t, updateBlk)
	updated := td.mustAccount(t, sender)
	require.Equal(t, rootB, updated.RootHash())
	require.Equal(t, executor.MinAnchorDeposit, updated.LockedDeposit())
	require.Equal(t, createBlk.Height(), updated.CreatedAtHeight())
	require.Equal(t, createBlk.Header().UnixTime(), updated.CreatedAtTime())
	require.Equal(t, updateBlk.Height(), updated.UpdatedAtHeight())
	require.Equal(t, updateBlk.Header().UnixTime(), updated.UpdatedAtTime())
	require.Equal(t, supply, td.totalCoins(t))

	// A top-up with the same content keeps the date of the current digest.
	topUpBlk := td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionSet, rootB, "update", 3, 1),
	})
	td.commit(t, topUpBlk)
	toppedUp := td.mustAccount(t, sender)
	require.Equal(t, executor.MinAnchorDeposit+3, toppedUp.LockedDeposit())
	require.Equal(t, updateBlk.Height(), toppedUp.UpdatedAtHeight())
	require.Equal(t, updateBlk.Header().UnixTime(), toppedUp.UpdatedAtTime())
	require.Equal(t, createBlk.Height(), toppedUp.CreatedAtHeight())
	require.Equal(t, supply, td.totalCoins(t))

	deleteBlk := td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionDelete, nil, "", 0, 1),
	})
	td.commit(t, deleteBlk)
	require.False(t, td.mustAccount(t, sender).HasAnchor())
	require.Len(t, td.accountBytes(t, sender), 12)
	require.Equal(t, supply, td.totalCoins(t))
}

// A certified block is trusted. If it still carries an anchor that cannot be
// stored, the node stops instead of crediting a fee that was never debited.
func TestCommitStopsOnUnstorableAnchor(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	template := td.proposeEmpty(t)
	empty := td.anchor(template.Height(), payload.AnchorActionSet, bytes.Repeat([]byte{0x51}, 32), "empty", 0, 0)
	blk := rebuildBlock(template, protocol.ProtocolVersion5, appendTxs(template.Transactions(), empty))
	require.ErrorIs(t, td.state.ValidateBlock(blk, 0), executor.ErrAnchorDepositTooSmall)

	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	require.Panics(t, func() {
		_ = td.state.CommitBlock(blk, cert)
	})
}

func TestConcreteSandboxUsesNextBlockTime(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)

	sbx := td.state.concreteSandbox()
	next := td.state.lastInfo.BlockTime().Add(td.state.params.BlockInterval())
	require.Equal(t, td.state.params.BlockVersion, sbx.BlockVersion())
	require.Equal(t, uint32(next.Unix()), sbx.CurrentUnixTime())
	require.NotZero(t, sbx.CurrentUnixTime())
}

// totalCoins sums spendable balances, locked anchor deposits and stakes.
func (td *testData) totalCoins(t *testing.T) amount.Amount {
	t.Helper()

	return totalCoinsOf(td.state)
}

func (td *testData) requireProposed(t *testing.T, want protocol.Version) {
	t.Helper()

	require.Equal(t, want, td.state.proposeBlockVersion())
	blk := td.proposeEmpty(t)
	require.Equal(t, want, blk.Header().Version())
}

func (td *testData) proposeEmpty(t *testing.T) *block.Block {
	t.Helper()

	return td.propose(t, block.Txs{})
}

func (td *testData) propose(t *testing.T, txs block.Txs) *block.Block {
	t.Helper()

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(txs).Times(1)
	blk, err := td.state.ProposeBlock(td.proposerKey(t, 0), td.RandAccAddress())
	require.NoError(t, err)

	return blk
}

func (td *testData) commit(t *testing.T, blk *block.Block) {
	t.Helper()

	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	require.NoError(t, td.state.CommitBlock(blk, cert))
}

func (td *testData) mustNotCommit(t *testing.T, blk *block.Block, want error) {
	t.Helper()

	root := td.state.stateRoot()
	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	err := td.state.CommitBlock(blk, cert)
	require.ErrorIs(t, err, want)
	require.Equal(t, root, td.state.stateRoot())
}

func (td *testData) nextHeight() types.Height {
	return td.state.LastBlockHeight() + 1
}

func (td *testData) sender() crypto.Address {
	return td.genAccKey.PublicKeyNative().AccountAddress()
}

func (td *testData) anchor(height types.Height, action uint8, root []byte, uri string,
	deposit, fee amount.Amount,
) *tx.Tx {
	trx := tx.NewAnchorTx(height, td.sender(), action, root, uri, 1, deposit, fee)
	td.HelperSignTransaction(td.genAccKey, trx)

	return trx
}

func (td *testData) mustAccount(t *testing.T, addr crypto.Address) *account.Account {
	t.Helper()

	acc, err := td.state.AccountByAddress(addr)
	require.NoError(t, err)

	return acc
}

func (td *testData) accountBytes(t *testing.T, addr crypto.Address) []byte {
	t.Helper()

	raw, err := td.mustAccount(t, addr).Bytes()
	require.NoError(t, err)

	return raw
}

func (td *testData) validatorHashes(t *testing.T) []hash.Hash {
	t.Helper()

	hashes := make([]hash.Hash, 4)
	for num := int32(0); num < 4; num++ {
		val, err := td.state.ValidatorByNumber(num)
		require.NoError(t, err)
		hashes[num] = val.Hash()
	}

	return hashes
}

func replaceCommittee(t *testing.T, chain *state, stakes []amount.Amount, versions []protocol.Version) {
	t.Helper()

	vals := chain.committee.Validators()
	require.Len(t, vals, len(stakes))
	require.Len(t, versions, len(stakes))

	for _, val := range vals {
		number := val.Number()
		if stakes[number] > 0 {
			val.AddToStake(stakes[number])
		}

		val.UpdateProtocolVersion(versions[number])
	}

	cmt, err := committee.NewCommittee(vals, chain.params.CommitteeSize, vals[0].Address())
	require.NoError(t, err)
	chain.committee = cmt
}

func applyCommittee(t *testing.T, chain *state, versions []protocol.Version) {
	t.Helper()

	vals := chain.committee.Validators()
	require.Len(t, vals, len(versions))

	for _, val := range vals {
		val.UpdateProtocolVersion(versions[val.Number()])
	}

	chain.committee.Update(0, vals)
}

func rebuildBlock(blk *block.Block, version protocol.Version, txs block.Txs) *block.Block {
	return block.MakeBlock(version, blk.Header().Time(), txs,
		blk.Header().PrevBlockHash(),
		blk.Header().StateRoot(),
		blk.PrevCertificate(),
		blk.Header().SortitionSeed(),
		blk.Header().ProposerAddress())
}

func subsidyPlus(subsidy *tx.Tx, extra amount.Amount) *tx.Tx {
	pld := subsidy.Payload().(*payload.BatchTransferPayload)
	recipients := make([]payload.BatchRecipient, len(pld.Recipients))
	copy(recipients, pld.Recipients)
	recipients[len(recipients)-1].Amount += extra

	return tx.NewSubsidyTx(subsidy.LockTime(), recipients)
}

func appendTxs(base block.Txs, extra ...*tx.Tx) block.Txs {
	out := make(block.Txs, 0, len(base)+len(extra))
	out = append(out, base...)
	out = append(out, extra...)

	return out
}

func blockHasAnchor(blk *block.Block) bool {
	for _, trx := range blk.Transactions() {
		if trx.Payload().Type() == payload.TypeAnchor {
			return true
		}
	}

	return false
}

func blockHasTx(blk *block.Block, id tx.ID) bool {
	for _, trx := range blk.Transactions() {
		if trx.ID() == id {
			return true
		}
	}

	return false
}
