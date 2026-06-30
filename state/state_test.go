package state

import (
	"testing"
	"time"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slices"
)

type testData struct {
	*testsuite.TestSuite

	state      *state
	fakeTxPool *txpool.FakeTxPool
	fakeStore  *store.FakeStore
	// genValKeys []*bls.ValidatorKey
	// genAccKey  *ed25519.PrivateKey
}

func setup(t *testing.T) *testData {
	t.Helper()

	return setupWithVersion(t, protocol.ProtocolVersionLatest)
}

func setupWithVersion(t *testing.T, blockVersion protocol.Version) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	genValNum := 4
	genValKeys := make([]*bls.ValidatorKey, 0, genValNum)
	genVals := make([]*validator.Validator, 0, genValNum)
	for i := 0; i < genValNum; i++ {
		valKey := ts.RandValKey()
		val := validator.NewValidator(valKey.PublicKey(), int32(i))

		genValKeys = append(genValKeys, valKey)
		genVals = append(genVals, val)
	}

	numBlocks := types.Height(7)
	fakeTxPool := txpool.NewFakeTxPool(ts)
	fakeTxPool.EXPECT().SetNewSandboxAndRecheck(gomock.Any()).Return().AnyTimes()
	fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(int(numBlocks))
	fakeTxPool.EXPECT().HandleCommittedBlock(gomock.Any()).Return().AnyTimes()

	fakeStore := store.NewFakeStore(ts)
	genTime := util.RoundNow(10).Add(-8640 * time.Second)

	genParams := genesis.DefaultGenesisParams()
	genParams.CommitteeSize = 7
	genParams.BondInterval = 10
	genParams.BlockVersion = blockVersion

	genAcc1 := account.NewAccount(0)
	genAcc1.AddToBalance(21 * 1e15) // 21,000,000.000,000,000
	genAcc2 := account.NewAccount(1)
	genAcc2.AddToBalance(21 * 1e15) // 21,000,000.000,000,000
	genAccPubKey, genAccPrvKey := ts.RandEd25519KeyPair()

	genAccs := map[crypto.Address]*account.Account{
		crypto.TreasuryAddress:        genAcc1,
		genAccPubKey.AccountAddress(): genAcc2,
	}

	gnDoc := genesis.MakeGenesis(genTime, genAccs, genVals, genParams)

	// First validator is in the committee
	valKeys := []*bls.ValidatorKey{genValKeys[0], ts.RandValKey()}
	eventPipe := pipeline.New[any](t.Context())
	st1, err := LoadOrNewState(t.Context(), gnDoc, valKeys, fakeStore, fakeTxPool, eventPipe)
	require.NoError(t, err)

	state, _ := st1.(*state)

	td := &testData{
		TestSuite:  ts,
		state:      state,
		fakeTxPool: fakeTxPool,
		genValKeys: genValKeys,
		genAccKey:  genAccPrvKey,
	}

	td.commitBlocks(t, numBlocks)

	return td
}

func (td *testData) proposerKey(t *testing.T, round types.Round) *bls.ValidatorKey {
	t.Helper()

	blockProposer := td.state.Proposer(round)
	valKeyIndex := slices.IndexFunc(td.genValKeys, func(e *bls.ValidatorKey) bool {
		return e.Address() == blockProposer.Address()
	})

	return td.genValKeys[valKeyIndex]
}

func (td *testData) makeBlockAndCertificate(t *testing.T, round types.Round) (
	*block.Block, *certificate.Certificate,
) {
	t.Helper()

	valKey := td.proposerKey(t, round)
	blk, _ := td.state.ProposeBlock(valKey, td.RandAccAddress())
	cert := td.makeCertificateAndSign(t, blk.Hash(), round)

	return blk, cert
}

func (td *testData) makeCertificateAndSign(t *testing.T, blockHash hash.Hash,
	round types.Round,
) *certificate.Certificate {
	t.Helper()

	sigs := make([]*bls.Signature, 0, len(td.genValKeys))
	height := td.state.LastBlockHeight()
	cert := certificate.NewCertificate(height+1, round)
	signBytes := cert.SignBytesPrecommit(blockHash)
	committers := []int32{0, 1, 2, 3}
	absentees := []int32{3}

	for _, key := range td.genValKeys[:len(td.genValKeys)-1] {
		sig := key.Sign(signBytes)
		sigs = append(sigs, sig)
	}

	aggSig, _ := bls.SignatureAggregate(sigs...)
	cert.SetSignature(committers, absentees, aggSig)

	return cert
}

func (td *testData) commitBlocks(t *testing.T, count types.Height) {
	t.Helper()

	for i := types.Height(0); i < count; i++ {
		blk, cert := td.makeBlockAndCertificate(t, 0)
		require.NoError(t, td.state.CommitBlock(blk, cert))
	}
}

func (td *testData) checkBlockSubsidy(t *testing.T, blk *block.Block) {
	t.Helper()

	accumulatedFee := amount.Amount(0)
	for _, trx := range blk.Transactions() {
		accumulatedFee += trx.Fee()
	}

	subsidyTrx := blk.Transactions().Subsidy()
	reward := td.state.params.BlockReward(blk.Height())
	assert.Equal(t, reward+accumulatedFee, subsidyTrx.Payload().Value())
}

func TestClosingState(t *testing.T) {
	td := setup(t)

	td.state.Close()
}

func TestBlockSubsidyTx(t *testing.T) {
	td := setup(t)

	// Without reward address in config
	rewardAddr := td.RandAccAddress()
	proposerAddr := td.state.Proposer(0).Address()
	height1 := td.RandHeight(testsuite.HeightWithMax(8_000_000))
	height2 := td.RandHeight(testsuite.HeightWithMin(8_000_000))
	foundationAddr1 := td.state.params.FoundationAddress(height1 + 1)
	foundationAddr2 := td.state.params.FoundationAddress(height2 + 1)

	tests := []struct {
		name               string
		height             types.Height
		accumulatedFee     amount.Amount
		expectedRecipients []payload.BatchRecipient
	}{
		{
			name:           "subsidy with zero transaction fee",
			height:         height1,
			accumulatedFee: 0,
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr1, Amount: 0.3e9},
				{To: rewardAddr, Amount: 0.7e9},
			},
		},

		{
			name:           "subsidy with transaction fee",
			height:         height1,
			accumulatedFee: 0.01e9, // 0.1 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr1, Amount: 0.3e9},
				{To: rewardAddr, Amount: 0.71e9},
			},
		},

		{
			name:           "subsidy with zero transaction fee, after first halving",
			height:         height2,
			accumulatedFee: 0,
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr2, Amount: 0.15e9},
				{To: rewardAddr, Amount: 0.35e9},
			},
		},

		{
			name:           "subsidy with transaction fee, after first halving",
			height:         height2,
			accumulatedFee: 0.01e9, // 0.1 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr2, Amount: 0.15e9},
				{To: rewardAddr, Amount: 0.36e9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cert := td.GenerateTestBlock(tt.height)
			td.state.lastInfo.UpdateCertificate(cert)
			trx := td.state.createSubsidyTx(proposerAddr, rewardAddr, tt.accumulatedFee)

			err := td.state.checkSubsidy(trx, td.genValKeys[0].Address(), true)
			require.NoError(t, err)

			payload := trx.Payload().(*payload.BatchTransferPayload)

			for i, recipient := range payload.Recipients {
				assert.Equal(t, tt.expectedRecipients[i].To, recipient.To)
				assert.Equal(t, tt.expectedRecipients[i].Amount, recipient.Amount)
			}
		})
	}
}

func TestBlockSubsidyWithDelegationTx(t *testing.T) {
	td := setup(t)

	// Without reward address in config
	rewardAddr := td.RandAccAddress()
	dlgOwnerAddr := td.RandAccAddress()
	proposerAddr := td.state.Proposer(0).Address()
	foundationAddr := td.state.params.FoundationAddress(td.state.LastBlockHeight() + 1)

	tests := []struct {
		name               string
		accumulatedFee     amount.Amount
		ownerShare         amount.Amount
		expectedRecipients []payload.BatchRecipient
	}{
		{
			name:           "owner share 0, without transaction fee",
			accumulatedFee: 0,
			ownerShare:     0,
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: rewardAddr, Amount: 0.7e9},
			},
		},
		{
			name:           "owner share 0.2, without transaction fee",
			accumulatedFee: 0,
			ownerShare:     0.2e9, // 0.2 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: dlgOwnerAddr, Amount: 0.2e9},
				{To: rewardAddr, Amount: 0.5e9},
			},
		},
		{
			name:           "owner share 0.7, without transaction fee",
			accumulatedFee: 0,
			ownerShare:     0.7e9, // 0.7 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: dlgOwnerAddr, Amount: 0.7e9},
			},
		},
		{
			name:           "owner share 0, with transaction fee",
			accumulatedFee: 0.01e9, // 0.01 PAC
			ownerShare:     0,
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: rewardAddr, Amount: 0.71e9},
			},
		},
		{
			name:           "owner share 0.2, with transaction fee",
			accumulatedFee: 0.01e9, // 0.01 PAC
			ownerShare:     0.2e9,  // 0.2 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: dlgOwnerAddr, Amount: 0.2e9},
				{To: rewardAddr, Amount: 0.51e9},
			},
		},
		{
			name:           "owner share 0.7, with transaction fee",
			accumulatedFee: 0.01e9, // 0.01 PAC
			ownerShare:     0.7e9,  // 0.7 PAC
			expectedRecipients: []payload.BatchRecipient{
				{To: foundationAddr, Amount: 0.3e9},
				{To: dlgOwnerAddr, Amount: 0.71e9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, _ := td.state.ValidatorByAddress(proposerAddr)
			val.SetDelegation(dlgOwnerAddr, tt.ownerShare, td.RandHeight())
			td.state.store.UpdateValidator(val)

			trx := td.state.createSubsidyTx(proposerAddr, rewardAddr, tt.accumulatedFee)

			err := td.state.checkSubsidy(trx, proposerAddr, true)
			require.NoError(t, err)

			payload := trx.Payload().(*payload.BatchTransferPayload)

			for i, recipient := range payload.Recipients {
				assert.Equal(t, tt.expectedRecipients[i].To, recipient.To)
				assert.Equal(t, tt.expectedRecipients[i].Amount, recipient.Amount)
			}
		})
	}
}

func TestTryCommitInvalidCertificate(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)

	blk, _ := td.makeBlockAndCertificate(t, td.RandRound())
	invCert := td.GenerateTestCertificate(td.state.LastBlockHeight() + 1)

	require.Error(t, td.state.CommitBlock(blk, invCert))
}

func TestTryCommitValidBlocks(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)

	blk, cert := td.makeBlockAndCertificate(t, 0)
	require.NoError(t, td.state.CommitBlock(blk, cert))

	// Commit again
	// No error here but block is ignored, because the height is invalid
	require.NoError(t, td.state.CommitBlock(blk, cert))

	assert.Equal(t, blk.Hash(), td.state.LastBlockHash())
	assert.Equal(t, blk.Header().Version(), td.state.Params().BlockVersion)
	assert.Equal(t, blk.Header().Time(), td.state.LastBlockTime())
	assert.Equal(t, cert.Hash(), td.state.LastCertificate().Hash())
	assert.Equal(t, cert.Height(), td.state.LastBlockHeight())
}

func TestCommitSandbox(t *testing.T) {
	t.Run("Add new account", func(t *testing.T) {
		td := setup(t)

		addr := td.RandAccAddress()
		sb := td.state.concreteSandbox()
		newAcc := sb.MakeNewAccount(addr)
		newAcc.AddToBalance(td.RandAmount())
		sb.UpdateAccount(addr, newAcc)
		td.state.commitSandbox(sb, 0)

		stateAcc, _ := td.state.AccountByAddress(addr)
		assert.Equal(t, stateAcc, newAcc)
	})

	t.Run("Add new validator", func(t *testing.T) {
		td := setup(t)

		pub, _ := td.RandBLSKeyPair()
		sb := td.state.concreteSandbox()
		newVal := sb.MakeNewValidator(pub)
		newVal.AddToStake(td.RandAmount())
		sb.UpdateValidator(newVal)
		td.state.commitSandbox(sb, 0)

		stateValByNumber, _ := td.state.ValidatorByAddress(pub.ValidatorAddress())
		stateValByAddr, _ := td.state.ValidatorByAddress(pub.ValidatorAddress())
		assert.Equal(t, stateValByNumber, newVal)
		assert.Equal(t, stateValByAddr, newVal)
	})

	t.Run("Modify account", func(t *testing.T) {
		td := setup(t)

		sbx := td.state.concreteSandbox()
		addr := td.genAccKey.PublicKeyNative().AccountAddress()
		acc := sbx.Account(addr)
		bal := acc.Balance()
		amt := td.RandAmount()
		acc.SubtractFromBalance(amt)
		sbx.UpdateAccount(addr, acc)
		td.state.commitSandbox(sbx, 0)

		stateAcc, _ := td.state.AccountByAddress(addr)
		assert.Equal(t, bal-amt, stateAcc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		td := setup(t)

		sbx := td.state.concreteSandbox()
		addr := td.genValKeys[0].Address()
		val := sbx.Validator(addr)
		stake := val.Stake()
		amt := td.RandAmount()
		val.AddToStake(amt)
		sbx.UpdateValidator(val)
		td.state.commitSandbox(sbx, 0)

		stateVal, _ := td.state.ValidatorByAddress(addr)
		assert.Equal(t, stake+amt, stateVal.Stake(), "%+v", val.Stake())
	})

	t.Run("Move committee", func(t *testing.T) {
		td := setup(t)

		proposer0 := td.state.committee.Proposer(0)
		proposer1 := td.state.committee.Proposer(1)
		assert.Equal(t, proposer0, td.state.committee.Proposer(0))

		sb := td.state.concreteSandbox()
		td.state.commitSandbox(sb, 0)

		assert.Equal(t, proposer1, td.state.committee.Proposer(0))
	})

	t.Run("Move committee next round", func(t *testing.T) {
		td := setup(t)

		proposer0 := td.state.committee.Proposer(0)
		proposer1 := td.state.committee.Proposer(1)
		proposer2 := td.state.committee.Proposer(2)
		assert.Equal(t, proposer0, td.state.committee.Proposer(0))
		assert.Equal(t, proposer1, td.state.committee.Proposer(1))

		sb := td.state.concreteSandbox()
		td.state.commitSandbox(sb, 1)

		assert.Equal(t, proposer2, td.state.committee.Proposer(0))
	})
}

func TestUpdateLastCertificate(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	blk, cert := td.makeBlockAndCertificate(t, 1)
	_ = td.state.CommitBlock(blk, cert)

	// the above `cert` is not signed by the last validators
	valKey1 := td.genValKeys[0]
	valKey4 := td.genValKeys[len(td.genValKeys)-1]
	invValKey := td.RandValKey()

	vote1 := vote.NewPrepareVote(blk.Hash(), cert.Height(), cert.Round(), valKey4.Address())
	vote2 := vote.NewPrecommitVote(blk.Hash(), cert.Height()+1, cert.Round(), valKey4.Address())
	vote3 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round()-1, valKey4.Address())
	vote4 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), valKey4.Address())
	vote5 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), invValKey.Address())
	vote6 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), valKey1.Address())
	vote7 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), valKey4.Address())

	td.HelperSignVote(valKey4, vote1)
	td.HelperSignVote(valKey4, vote2)
	td.HelperSignVote(valKey4, vote3)
	td.HelperSignVote(invValKey, vote4)
	td.HelperSignVote(invValKey, vote5)
	td.HelperSignVote(valKey4, vote6)
	td.HelperSignVote(valKey4, vote7)

	tests := []struct {
		vote   *vote.Vote
		err    error
		reason string
	}{
		{vote1, InvalidVoteForCertificateError{Vote: vote1}, "invalid vote type"},
		{vote2, InvalidVoteForCertificateError{Vote: vote2}, "invalid height"},
		{vote3, InvalidVoteForCertificateError{Vote: vote3}, "invalid round"},
		{vote4, crypto.ErrInvalidSignature, "invalid signature"},
		{vote5, store.ErrNotFound, "unknown validator"},
		{vote6, InvalidVoteForCertificateError{Vote: vote6}, "not in absentee"},
		{vote7, nil, "ok"},
	}

	for no, tt := range tests {
		err := td.state.UpdateLastCertificate(tt.vote)
		require.ErrorIs(t, err, tt.err, "error not matched for test %v", no)
	}
}

func TestForkDetection(t *testing.T) {
	td := setup(t)

	t.Run("Two blocks with different previous block hashes", func(t *testing.T) {
		td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)

		assert.Panics(t, func() {
			blk0, _ := td.makeBlockAndCertificate(t, 0)
			blkFork := block.MakeBlock(
				blk0.Header().Version(),
				blk0.Header().Time(),
				blk0.Transactions(),
				td.RandHash(),
				blk0.Header().StateRoot(),
				blk0.PrevCertificate(),
				blk0.Header().SortitionSeed(),
				blk0.Header().ProposerAddress(),
			)
			certFork := td.makeCertificateAndSign(t, blkFork.Hash(), 0)

			_ = td.state.CommitBlock(blkFork, certFork)
		})
	})
}

func TestSortition(t *testing.T) {
	td := setup(t)

	valKey := td.state.valKeys[1]
	assert.False(t, td.state.evaluateSortition()) //  not a validator
	assert.Equal(t, int64(4), td.state.ChainInfo().CommitteePower)

	trx := tx.NewBondTx(1, td.genAccKey.PublicKeyNative().AccountAddress(),
		valKey.Address(), valKey.PublicKey(), 1000000000, 100000)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{trx}).Times(1)
	td.commitBlocks(t, 1)

	assert.False(t, td.state.evaluateSortition()) // bonding period
	assert.Equal(t, int64(4), td.state.ChainInfo().CommitteePower)
	assert.False(t, td.state.committee.Contains(valKey.Address())) // Not in the committee

	// Committing another 10 blocks
	var sortitionTrx *tx.Tx

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(10)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Do(func(trx *tx.Tx) {
		sortitionTrx = trx // Capture the input argument here
	}).Return(nil).Times(1)
	td.commitBlocks(t, 10)

	assert.Equal(t, payload.TypeSortition, sortitionTrx.Payload().Type())
	assert.Equal(t, valKey.Address(), sortitionTrx.Payload().Signer())
	assert.False(t, td.state.committee.Contains(valKey.Address())) // Still not in the committee

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{sortitionTrx}).Times(1)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).Times(1)

	td.commitBlocks(t, 1)

	assert.Equal(t, int64(1000000004), td.state.ChainInfo().CommitteePower)
	assert.True(t, td.state.committee.Contains(valKey.Address())) // In the committee
}

func TestValidateBlockTime(t *testing.T) {
	td := setup(t)

	t.Run("Time is not rounded", func(t *testing.T) {
		roundedNow := util.RoundNow(10)

		require.Error(t, td.state.validateBlockTime(roundedNow.Add(-15*time.Second)))
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(-5*time.Second)))
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(5*time.Second)))
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(15*time.Second)))
	})

	t.Run("Last block is committed 10 seconds ago", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(-10 * time.Second))

		// Before or same as the last block time
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(-20*time.Second)))
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(-10*time.Second)))

		// Ok
		require.NoError(t, td.state.validateBlockTime(roundedNow))
		require.NoError(t, td.state.validateBlockTime(roundedNow.Add(10*time.Second)))

		// More than the threshold
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(20*time.Second)))

		expectedProposeTime := roundedNow
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})

	t.Run("Last block is committed one minute ago", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(-1 * time.Minute)) // One minute ago
		lastBlockTime := td.state.LastBlockTime()

		// Before or same as the last block time
		require.Error(t, td.state.validateBlockTime(lastBlockTime.Add(-10*time.Second)))
		require.Error(t, td.state.validateBlockTime(lastBlockTime))

		// Ok
		require.NoError(t, td.state.validateBlockTime(roundedNow.Add(-10*time.Second)))
		require.NoError(t, td.state.validateBlockTime(roundedNow))
		require.NoError(t, td.state.validateBlockTime(roundedNow.Add(10*time.Second)))

		// More than the threshold
		require.Error(t, td.state.validateBlockTime(roundedNow.Add(30*time.Second)))

		expectedProposeTime := util.RoundNow(10)
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})

	t.Run("Last block is committed in future", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(1 * time.Minute)) // One minute later
		lastBlockTime := td.state.LastBlockTime()

		require.Error(t, td.state.validateBlockTime(lastBlockTime.Add(+1*time.Minute)))

		// Before the last block time
		require.Error(t, td.state.validateBlockTime(lastBlockTime.Add(-10*time.Second)))
		require.Error(t, td.state.validateBlockTime(lastBlockTime))

		// Ok
		require.NoError(t, td.state.validateBlockTime(lastBlockTime.Add(10*time.Second)))
		require.NoError(t, td.state.validateBlockTime(lastBlockTime.Add(20*time.Second)))

		// More than the threshold
		require.Error(t, td.state.validateBlockTime(lastBlockTime.Add(30*time.Second)))

		expectedProposeTime := roundedNow.Add(1 * time.Minute).Add(
			time.Duration(td.state.params.BlockIntervalInSecond) * time.Second,
		)
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})
}

func TestValidatorHelpers(t *testing.T) {
	td := setup(t)

	t.Run("Should return error non-existing validator by address", func(t *testing.T) {
		_, err := td.state.ValidatorByAddress(td.RandValAddress())
		require.Error(t, err)
	})

	t.Run("Should return error for non-existing validator by number", func(t *testing.T) {
		_, err := td.state.ValidatorByNumber(10)
		require.Error(t, err)
	})

	t.Run("Should return validator for existing validator by address", func(t *testing.T) {
		existingValidator, err := td.state.ValidatorByAddress(td.genValKeys[0].Address())
		require.NoError(t, err)
		assert.Zero(t, existingValidator.Number())
	})

	t.Run("Should return validator for existing validator by number", func(t *testing.T) {
		existingValidator, err := td.state.ValidatorByNumber(0)
		require.NoError(t, err)
		assert.Zero(t, existingValidator.Number())
	})
}

func TestLoadState(t *testing.T) {
	td := setup(t)

	// Add a bond transactions to change total power (stake)
	pub, _ := td.RandBLSKeyPair()
	lockTime := td.state.LastBlockHeight()
	bondTrx := tx.NewBondTx(lockTime, td.genAccKey.PublicKeyNative().AccountAddress(),
		pub.ValidatorAddress(), pub, 1000000000, 100000)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{bondTrx}).Times(1)
	blk5, cert5 := td.makeBlockAndCertificate(t, 1)
	assert.Equal(t, 2, blk5.Transactions().Len())
	require.NoError(t, td.state.CommitBlock(blk5, cert5))

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	blk6, cert6 := td.makeBlockAndCertificate(t, 0)

	// Load last state info
	eventPipe := pipeline.New[any](t.Context())
	newState, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
		td.state.store, td.fakeTxPool, eventPipe)
	require.NoError(t, err)

	assert.Equal(t, td.state.Params(), newState.Params())
	assert.ElementsMatch(t, td.state.ValidatorAddresses(), newState.ValidatorAddresses())
	assert.Equal(t, td.state.ChainInfo(), newState.ChainInfo())

	// Try committing the next block
	require.NoError(t, newState.CommitBlock(blk6, cert6))
}

func TestIsInCommittee(t *testing.T) {
	td := setup(t)

	proposer0 := td.proposerKey(t, 0).Address()
	proposer1 := td.proposerKey(t, 1).Address()
	assert.True(t, td.state.IsInCommittee(proposer0))
	assert.True(t, td.state.IsProposer(proposer0, 0))
	assert.True(t, td.state.IsInCommittee(proposer1))
	assert.True(t, td.state.IsProposer(proposer1, 1))

	addr := td.RandAccAddress()
	assert.False(t, td.state.IsInCommittee(addr))
	assert.False(t, td.state.IsProposer(addr, 0))
	assert.False(t, td.state.IsInCommittee(addr))
}

func TestCalculateFee(t *testing.T) {
	td := setup(t)

	expectedFee := td.RandFee()
	td.fakeTxPool.EXPECT().EstimatedFee(gomock.Any(), payload.TypeTransfer).Return(expectedFee).Times(1)

	fee := td.state.CalculateFee(td.RandAmount(), payload.TypeTransfer)

	assert.Equal(t, expectedFee, fee)
}

func TestCheckMaximumTransactionPerBlock(t *testing.T) {
	td := setup(t)

	txs := block.Txs{}
	td.state.params.MaxTransactionsPerBlock = 10
	lockTime := td.state.LastBlockHeight()
	senderAddr := td.genAccKey.PublicKeyNative().AccountAddress()
	for i := 0; i < td.state.params.MaxTransactionsPerBlock+2; i++ {
		trx := tx.NewTransferTx(lockTime, senderAddr,
			td.RandAccAddress(), td.RandAmount(), td.RandFee())

		txs = append(txs, trx)
	}

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(txs).Times(1)
	blk, err := td.state.ProposeBlock(td.state.valKeys[0], td.RandAccAddress())
	require.NoError(t, err)
	assert.Equal(t, td.state.params.MaxTransactionsPerBlock, blk.Transactions().Len())
}

func TestCommittedBlock(t *testing.T) {
	td := setup(t)

	t.Run("Block at 0", func(t *testing.T) {
		cBlkZero, err := td.state.CommittedBlock(0)
		require.Error(t, err)
		assert.Nil(t, cBlkZero)
		assert.Equal(t, hash.UndefHash, td.state.BlockHash(0))
		assert.Equal(t, types.Height(0), td.state.BlockHeight(hash.UndefHash))
	})

	t.Run("First block (Genesis)", func(t *testing.T) {
		cBlkOne, err := td.state.CommittedBlock(1)
		require.NoError(t, err)
		blkOne, err := cBlkOne.ToBlock()
		require.NoError(t, err)
		assert.Nil(t, blkOne.PrevCertificate())
		assert.Equal(t, hash.UndefHash, blkOne.Header().PrevBlockHash())
	})

	t.Run("Last block", func(t *testing.T) {
		cBlkLast, err := td.state.CommittedBlock(td.state.LastBlockHeight())
		require.NoError(t, err)
		blkLast, err := cBlkLast.ToBlock()
		require.NoError(t, err)
		assert.Equal(t, blkLast.Hash(), td.state.LastBlockHash())
	})
}

func TestUpdateProptocolVersion(t *testing.T) {
	td := setup(t)

	val, err := td.state.ValidatorByAddress(td.state.valKeys[0].Address())
	require.NoError(t, err)
	assert.Equal(t, protocol.ProtocolVersionLatest, val.ProtocolVersion())
}

func TestBlockVersionUpgrade(t *testing.T) {
	td1 := setupWithVersion(t, protocol.ProtocolVersionLatest-1)
	td2 := setupWithVersion(t, protocol.ProtocolVersionLatest)

	td1.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
	td2.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()

	blk1, cert1 := td1.makeBlockAndCertificate(t, 0)
	require.NoError(t, td1.state.CommitBlock(blk1, cert1))
	assert.Equal(t, protocol.ProtocolVersionLatest-1, td1.state.Params().BlockVersion)

	blk2, cert2 := td2.makeBlockAndCertificate(t, 0)
	require.NoError(t, td2.state.CommitBlock(blk2, cert2))
	assert.Equal(t, protocol.ProtocolVersionLatest, td2.state.Params().BlockVersion)
}

// func TestProposeBlockVersionUpgradeToV4(t *testing.T) {
// 	// When BlockVersion is already V4, proposeBlockVersion returns V4 directly.
// 	td := setupWithVersion(t, protocol.ProtocolVersion4)

// 	assert.Equal(t, protocol.ProtocolVersion4, td.state.proposeBlockVersion())

// 	td.mockTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
// 	valKey := td.proposerKey(t, 0)
// 	blk, err := td.state.ProposeBlock(valKey, td.RandAccAddress())
// 	require.NoError(t, err)
// 	assert.Equal(t, protocol.ProtocolVersion4, blk.Header().Version())
// }

// func TestProposeBlockVersionStaysAtV3(t *testing.T) {
// 	// When BlockVersion is V3 and committee validators don't support V4 yet,
// 	// proposeBlockVersion should return V3.
// 	td := setupWithVersion(t, protocol.ProtocolVersion3)

// 	assert.Equal(t, protocol.ProtocolVersion3, td.state.proposeBlockVersion())

// 	td.mockTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
// 	valKey := td.proposerKey(t, 0)
// 	blk, err := td.state.ProposeBlock(valKey, td.RandAccAddress())
// 	require.NoError(t, err)
// 	assert.Equal(t, protocol.ProtocolVersion3, blk.Header().Version())
// }

// func TestRewardHalvingWithV4(t *testing.T) {
// 	// Setup with V4 block version
// 	td := setupWithVersion(t, protocol.ProtocolVersion4)

// 	// At the test height (7), the coefficient is 1.0 since 7 <= 8,000,000.
// 	// So reward should be the same as pre-V4: 1 PAC.
// 	proposerAddr := td.state.Proposer(0).Address()
// 	trx := td.state.createSubsidyTx(proposerAddr, td.RandAccAddress(), 0, protocol.ProtocolVersion4)
// 	batchTrx := trx.Payload().(*payload.BatchTransferPayload)

// 	// Foundation reward should be 0.3 * 1.0 = 0.3 PAC
// 	assert.Equal(t, amount.Amount(0.3e9), batchTrx.Recipients[0].Amount)
// 	// Validator reward should be 1.0 - 0.3 = 0.7 PAC
// 	assert.Equal(t, amount.Amount(0.7e9), batchTrx.Recipients[1].Amount)

// 	// With V3, the same height gives the same result.
// 	trxV3 := td.state.createSubsidyTx(proposerAddr, td.RandAccAddress(), 0, protocol.ProtocolVersion3)
// 	batchTrxV3 := trxV3.Payload().(*payload.BatchTransferPayload)
// 	assert.Equal(t, batchTrx.Recipients[0].Amount, batchTrxV3.Recipients[0].Amount)
// 	assert.Equal(t, batchTrx.Recipients[1].Amount, batchTrxV3.Recipients[1].Amount)
// }
