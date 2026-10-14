package state

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateLastCertificateAtGenesis(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)
	require.Zero(t, td.state.LastBlockHeight())
	require.Nil(t, td.state.LastCertificate(), "genesis block has no certificate")

	signer := td.genValKeys[3]
	vte := vote.NewPrecommitVote(td.RandHash(), 1, 0, signer.Address())
	td.HelperSignVote(signer, vte)

	var err error
	require.NotPanics(t, func() { err = td.state.UpdateLastCertificate(vte) })
	require.ErrorIs(t, err, InvalidVoteForCertificateError{Vote: vte})
}

func TestChainInfoSyncProgress(t *testing.T) {
	setGenesisAndLastBlock := func(td *testData, genTime, lastBlockTime time.Time) {
		td.state.genDoc = genesis.MakeGenesis(genTime,
			td.state.genDoc.Accounts(), td.state.genDoc.Validators(), td.state.genDoc.Params())
		td.state.lastInfo.UpdateBlockTime(lastBlockTime)
	}

	t.Run("half of the chain is synced", func(t *testing.T) {
		td := setup(t)
		now := time.Now()
		setGenesisAndLastBlock(td, now.Add(-1000*time.Second), now.Add(-500*time.Second))

		info := td.state.ChainInfo()
		assert.InDelta(t, 0.5, info.SyncProgress, 0.01)
		assert.Equal(t, int64(50), info.BlocksLeft) // 500 seconds with 10 seconds block interval
	})

	t.Run("last block is slightly ahead of the local clock", func(t *testing.T) {
		td := setup(t)
		now := time.Now()
		setGenesisAndLastBlock(td, now.Add(-1000*time.Second), now.Add(30*time.Second))

		info := td.state.ChainInfo()
		assert.InDelta(t, 1.0, info.SyncProgress, 1e-9)
		assert.Zero(t, info.BlocksLeft)
	})

	t.Run("genesis time is in the future", func(t *testing.T) {
		td := setup(t)
		genTime := time.Now().Add(time.Hour)
		setGenesisAndLastBlock(td, genTime, genTime)

		info := td.state.ChainInfo()
		assert.InDelta(t, 1.0, info.SyncProgress, 1e-9)
		assert.Zero(t, info.BlocksLeft)
	})

	t.Run("genesis time is now", func(t *testing.T) {
		// The elapsed time since the genesis is zero, or one second at most.
		// Dividing by zero should never produce NaN or Inf.
		td := setup(t)
		now := time.Now()
		setGenesisAndLastBlock(td, now, now)

		info := td.state.ChainInfo()
		assert.False(t, math.IsNaN(info.SyncProgress))
		assert.False(t, math.IsInf(info.SyncProgress, 0))
		assert.GreaterOrEqual(t, info.SyncProgress, 0.0)
		assert.LessOrEqual(t, info.SyncProgress, 1.0)
	})
}

func TestCommitBlockExecutionFailure(t *testing.T) {
	td := setup(t)

	// A valid block, but with an extra transaction that cannot be executed:
	// the sender account doesn't exist.
	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	valKey := td.proposerKey(t, 0)
	validBlk, err := td.state.ProposeBlock(valKey, td.RandAccAddress())
	require.NoError(t, err)

	badTrx := td.GenerateTestTransferTx()
	txs := append(block.Txs{validBlk.Transactions().Subsidy()}, badTrx)
	badBlk := block.MakeBlock(validBlk.Header().Version(), validBlk.Header().Time(), txs,
		validBlk.Header().PrevBlockHash(), validBlk.Header().StateRoot(),
		validBlk.PrevCertificate(), validBlk.Header().SortitionSeed(), validBlk.Header().ProposerAddress())
	cert := td.makeCertificateAndSign(t, badBlk.Hash(), 0)

	lastHeight := td.state.LastBlockHeight()
	lastHash := td.state.LastBlockHash()
	lastTime := td.state.LastBlockTime()
	lastCert := td.state.LastCertificate()
	stateRoot := td.state.stateRoot()
	totalPower := td.state.totalPower

	require.Error(t, td.state.CommitBlock(badBlk, cert))

	// The state is not changed by the failed commit.
	assert.Equal(t, lastHeight, td.state.LastBlockHeight())
	assert.Equal(t, lastHash, td.state.LastBlockHash())
	assert.Equal(t, lastTime, td.state.LastBlockTime())
	assert.Equal(t, lastCert.Hash(), td.state.LastCertificate().Hash())
	assert.Equal(t, stateRoot, td.state.stateRoot())
	assert.Equal(t, totalPower, td.state.totalPower)
	assert.NotContains(t, td.fakeStore.FakeBlocks, lastHeight+1, "the failed block should not be saved")
	assert.NotContains(t, td.fakeStore.FakeCertificates, lastHeight+1)

	// The state is still able to commit the valid block.
	validCert := td.makeCertificateAndSign(t, validBlk.Hash(), 0)
	require.NoError(t, td.state.CommitBlock(validBlk, validCert))
	assert.Equal(t, lastHeight+1, td.state.LastBlockHeight())
}

func TestLoadStateFailures(t *testing.T) {
	t.Run("genesis state can't be written", func(t *testing.T) {
		td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)
		errWrite := errors.New("disk is full")

		mockStore := store.NewMockStore(td.MockController())
		mockStore.EXPECT().LastCertificate().Return(nil).Times(1)
		mockStore.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).AnyTimes()
		mockStore.EXPECT().UpdateValidator(gomock.Any()).AnyTimes()
		mockStore.EXPECT().WriteBatch().Return(errWrite).Times(1)

		_, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
			mockStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		require.ErrorIs(t, err, errWrite)
	})

	t.Run("genesis without validators", func(t *testing.T) {
		td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)
		genDoc := genesis.MakeGenesis(time.Now(), td.state.genDoc.Accounts(), nil, td.state.genDoc.Params())

		// The store must not be written when the genesis is invalid.
		mockStore := store.NewMockStore(td.MockController())
		mockStore.EXPECT().LastCertificate().Return(nil).Times(1)

		var err error
		require.NotPanics(t, func() {
			_, err = LoadOrNewState(t.Context(), genDoc, td.state.valKeys,
				mockStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		})
		require.ErrorIs(t, err, ErrNoGenesisValidator)
	})

	t.Run("genesis without accounts", func(t *testing.T) {
		td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)
		genDoc := genesis.MakeGenesis(time.Now(), nil, td.state.genDoc.Validators(), td.state.genDoc.Params())

		// The store must not be written when the genesis is invalid.
		mockStore := store.NewMockStore(td.MockController())
		mockStore.EXPECT().LastCertificate().Return(nil).Times(1)

		var err error
		require.NotPanics(t, func() {
			_, err = LoadOrNewState(t.Context(), genDoc, td.state.valKeys,
				mockStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		})
		require.ErrorIs(t, err, ErrNoGenesisAccount)
	})

	t.Run("last block can't be read", func(t *testing.T) {
		td := setup(t)
		errRead := errors.New("db is corrupted")

		mockStore := store.NewMockStore(td.MockController())
		mockStore.EXPECT().LastCertificate().Return(td.state.LastCertificate()).AnyTimes()
		mockStore.EXPECT().Block(gomock.Any()).Return(nil, errRead).Times(1)

		_, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
			mockStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		require.ErrorIs(t, err, errRead)
	})

	t.Run("a block of the score window is missing", func(t *testing.T) {
		td := setup(t)
		require.Contains(t, td.fakeStore.FakeBlocks, types.Height(3))

		delete(td.fakeStore.FakeBlocks, 3)

		_, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
			td.fakeStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		require.ErrorIs(t, err, store.ErrNotFound)
	})

	t.Run("a block of the score window is corrupted", func(t *testing.T) {
		td := setup(t)

		corruptedStore := &truncatedBlockStore{FakeStore: td.fakeStore, height: 3}

		_, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
			corruptedStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		require.Error(t, err)
	})
}

// truncatedBlockStore returns the block at the given height without
// the end of its data, as a corrupted database would do.
type truncatedBlockStore struct {
	*store.FakeStore

	height types.Height
}

func (s *truncatedBlockStore) Block(height types.Height) (*store.CommittedBlock, error) {
	cBlk, err := s.FakeStore.Block(height)
	if err == nil && height == s.height {
		// A block header takes 138 bytes, followed by the certificate of the previous block.
		cBlk.Data = cBlk.Data[:140]
	}

	return cBlk, err
}

func TestEvaluateSortitionSkippedValidators(t *testing.T) {
	t.Run("unbonded validator doesn't evaluate sortition", func(t *testing.T) {
		td := setup(t)

		// Leave the bonding period, and give a stake much bigger than the total power:
		// the sortition always succeeds for this validator.
		td.state.params.BondInterval = 1
		val, err := td.state.ValidatorByAddress(td.state.valKeys[0].Address())
		require.NoError(t, err)
		val.AddToStake(1_000_000_000)

		td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).Times(1)
		require.True(t, td.state.evaluateSortition())

		// After unbonding, no sortition transaction is made.
		val.UpdateUnbondingHeight(td.state.LastBlockHeight())
		assert.False(t, td.state.evaluateSortition())
	})

	t.Run("sortition transaction rejected by the transaction pool", func(t *testing.T) {
		td := setup(t)
		valKey := td.state.valKeys[1]

		bondTrx := tx.NewBondTx(1, td.genAccKey.PublicKeyNative().AccountAddress(),
			valKey.Address(), valKey.PublicKey(), 1000000000, 100000)
		td.commitBlocks(t, 1, block.Txs{bondTrx})

		// The bond transaction is committed, and the bonding period is over after 10 blocks.
		// All the sortition transactions are rejected, so the validator never joins the committee.
		errRejected := errors.New("rejected")
		td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(errRejected).AnyTimes()
		td.commitBlocks(t, 12)

		assert.False(t, td.state.evaluateSortition())
		assert.False(t, td.state.committee.Contains(valKey.Address()))
	})
}

// PIP-54: the block at height 7,406,821 is validated with the committee of the genesis block.
func TestCommitBlockRollbackToGenesisCommittee(t *testing.T) {
	const rollbackHeight = types.Height(7_406_821)

	td := setup(t)
	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
	// The validators of the test can win the sortition at this height.
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	// Move the state to the block before the rollback.
	prevCert := certificate.NewCertificate(rollbackHeight-1, 0)
	signBytes := prevCert.SignBytesPrecommit(td.state.LastBlockHash())
	sigs := make([]*bls.Signature, 0, len(td.genValKeys))
	for _, key := range td.genValKeys[:len(td.genValKeys)-1] {
		sigs = append(sigs, key.Sign(signBytes))
	}
	aggSig, err := bls.SignatureAggregate(sigs...)
	require.NoError(t, err)
	prevCert.SetSignature([]int32{0, 1, 2, 3}, []int32{3}, aggSig)
	td.state.lastInfo.UpdateCertificate(prevCert)
	require.Equal(t, rollbackHeight-1, td.state.LastBlockHeight())

	// The committee has moved forward since the genesis block,
	// and the first validator is not the proposer anymore.
	td.state.committee.Update(2, nil)
	valKey := td.genValKeys[0]
	require.NotEqual(t, valKey.Address(), td.state.Proposer(0).Address())

	blk, err := td.state.ProposeBlock(valKey, td.RandAccAddress())
	require.NoError(t, err)
	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	require.Equal(t, rollbackHeight, cert.Height())

	// Without the rollback, the proposer of the block is not valid.
	require.ErrorIs(t, td.state.ValidateBlock(blk, 0), InvalidProposerError{
		Expected: td.state.Proposer(0).Address(),
		Got:      valKey.Address(),
	})

	// With the rollback, the first validator of the genesis is the proposer.
	require.NoError(t, td.state.CommitBlock(blk, cert))
	assert.Equal(t, rollbackHeight, td.state.LastBlockHeight())
	assert.Equal(t, blk.Hash(), td.state.LastBlockHash())
}
