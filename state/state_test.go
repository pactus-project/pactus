package state

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type testData struct {
	*testsuite.TestSuite

	state        *state
	genValKeys   []*bls.ValidatorKey
	genAccKey    *ed25519.PrivateKey
	commonTxPool *txpool.MockTxPool
}

func setup(t *testing.T) *testData {
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

	mockTxPool := txpool.MockingTxPool()
	mockStore := store.MockingStore(ts)

	genTime := util.RoundNow(10).Add(-8640 * time.Second)

	genParams := genesis.DefaultGenesisParams()
	genParams.CommitteeSize = 7
	genParams.BondInterval = 10

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
	st1, err := LoadOrNewState(gnDoc, valKeys, mockStore, mockTxPool)
	require.NoError(t, err)

	state, _ := st1.(*state)

	td := &testData{
		TestSuite:    ts,
		state:        state,
		genValKeys:   genValKeys,
		genAccKey:    genAccPrvKey,
		commonTxPool: mockTxPool,
	}

	td.commitBlocks(t, 8)

	return td
}

func (td *testData) makeBlockAndCertificate(t *testing.T, round int16) (
	*block.Block, *certificate.BlockCertificate,
) {
	t.Helper()

	blockProposer := td.state.Proposer(round)
	valKeyIndex := slices.IndexFunc(td.genValKeys, func(e *bls.ValidatorKey) bool {
		return e.Address() == blockProposer.Address()
	})
	valKey := td.genValKeys[valKeyIndex]
	blk, _ := td.state.ProposeBlock(valKey, td.RandAccAddress())
	cert := td.makeCertificateAndSign(t, blk.Hash(), round)

	return blk, cert
}

func (td *testData) makeCertificateAndSign(t *testing.T, blockHash hash.Hash,
	round int16,
) *certificate.BlockCertificate {
	t.Helper()

	sigs := make([]*bls.Signature, 0, len(td.genValKeys))
	height := td.state.LastBlockHeight()
	cert := certificate.NewBlockCertificate(height+1, round)
	signBytes := cert.SignBytes(blockHash)
	committers := []int32{0, 1, 2, 3}
	absentees := []int32{3}

	for _, key := range td.genValKeys[:len(td.genValKeys)-1] {
		sig := key.Sign(signBytes)
		sigs = append(sigs, sig)
	}

	cert.SetSignature(committers, absentees, bls.SignatureAggregate(sigs...))

	return cert
}

func (td *testData) commitBlocks(t *testing.T, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		blk, cert := td.makeBlockAndCertificate(t, 0)
		assert.NoError(t, td.state.CommitBlock(blk, cert))
	}
}

func TestClosingState(t *testing.T) {
	td := setup(t)

	td.state.Close()
}

func TestBlockSubsidyTx(t *testing.T) {
	td := setup(t)

	// Without reward address in config
	rewardAddr := td.RandAccAddress()
	randAccumulatedFee := td.RandFee()
	trx := td.state.createSubsidyTx(rewardAddr, randAccumulatedFee)
	assert.True(t, trx.IsSubsidyTx())
	assert.Equal(t, td.state.params.BlockReward+randAccumulatedFee, trx.Payload().Value())
	assert.Equal(t, crypto.TreasuryAddress, trx.Payload().(*payload.TransferPayload).From)
	assert.Equal(t, rewardAddr, trx.Payload().(*payload.TransferPayload).To)
}

func TestGenesisHash(t *testing.T) {
	td := setup(t)

	gen := td.state.Genesis()
	genAccs := gen.Accounts()
	genVals := gen.Validators()

	assert.NotNil(t, genAccs, td.genAccKey.PublicKeyNative().AccountAddress())
	assert.NotNil(t, genVals, td.genValKeys[0].Address())
}

func TestTryCommitInvalidCertificate(t *testing.T) {
	td := setup(t)

	blk, _ := td.makeBlockAndCertificate(t, td.RandRound())
	invCert := td.GenerateTestBlockCertificate(td.state.LastBlockHeight() + 1)

	assert.Error(t, td.state.CommitBlock(blk, invCert))
}

func TestTryCommitValidBlocks(t *testing.T) {
	td := setup(t)

	blk, crt := td.makeBlockAndCertificate(t, 0)

	assert.NoError(t, td.state.CommitBlock(blk, crt))

	// Commit again
	// No error here but block is ignored, because the height is invalid
	assert.NoError(t, td.state.CommitBlock(blk, crt))

	assert.Equal(t, blk.Hash(), td.state.LastBlockHash())
	assert.Equal(t, blk.Header().Time(), td.state.LastBlockTime())
	assert.Equal(t, crt.Hash(), td.state.LastCertificate().Hash())
	assert.Equal(t, uint32(9), td.state.LastBlockHeight())
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

		stateAcc := td.state.AccountByAddress(addr)
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

		stateValByNumber := td.state.ValidatorByAddress(pub.ValidatorAddress())
		stateValByAddr := td.state.ValidatorByAddress(pub.ValidatorAddress())
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

		stateAcc := td.state.AccountByAddress(addr)
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

		stateVal := td.state.ValidatorByAddress(addr)
		assert.Equal(t, stake+amt, stateVal.Stake(), val.Stake())
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
		assert.ErrorIs(t, err, tt.err, "error not matched for test %v", no)
	}
}

func TestBlockProposal(t *testing.T) {
	td := setup(t)

	t.Run("validity of the proposed block", func(t *testing.T) {
		b, err := td.state.ProposeBlock(td.state.valKeys[0], td.RandAccAddress())
		assert.NoError(t, err)
		assert.NoError(t, td.state.ValidateBlock(b, 0))
	})

	t.Run("Tx pool has two subsidy transactions", func(t *testing.T) {
		trx := td.state.createSubsidyTx(td.RandAccAddress(), 0)
		assert.NoError(t, td.state.AddPendingTx(trx))

		b, err := td.state.ProposeBlock(td.state.valKeys[0], td.RandAccAddress())
		assert.NoError(t, err)
		assert.NoError(t, td.state.ValidateBlock(b, 0))
		assert.Equal(t, 1, b.Transactions().Len())
	})
}

func TestForkDetection(t *testing.T) {
	td := setup(t)

	t.Run("Two certificates with different rounds", func(t *testing.T) {
		blk, certMain := td.makeBlockAndCertificate(t, 0)
		certFork := td.makeCertificateAndSign(t, blk.Hash(), 1)

		assert.NoError(t, td.state.CommitBlock(blk, certMain))
		assert.NoError(t, td.state.CommitBlock(blk, certFork)) // TODO: should panic here
	})

	t.Run("Two blocks with different previous block hashes", func(t *testing.T) {
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
				blk0.Header().ProposerAddress())
			certFork := td.makeCertificateAndSign(t, blkFork.Hash(), 0)

			_ = td.state.CommitBlock(blkFork, certFork)
		})
	})
}

func TestSortition(t *testing.T) {
	td := setup(t)

	secValKey := td.state.valKeys[1]
	assert.False(t, td.state.evaluateSortition()) //  not a validator
	assert.False(t, td.state.IsValidator(secValKey.Address()))
	assert.Equal(t, int64(4), td.state.CommitteePower())

	trx := tx.NewBondTx(1, td.genAccKey.PublicKeyNative().AccountAddress(),
		secValKey.Address(), secValKey.PublicKey(), 1000000000, 100000)
	td.HelperSignTransaction(td.genAccKey, trx)
	assert.NoError(t, td.state.AddPendingTx(trx))

	td.commitBlocks(t, 1)

	assert.False(t, td.state.evaluateSortition()) // bonding period
	assert.True(t, td.state.IsValidator(secValKey.Address()))
	assert.Equal(t, int64(4), td.state.CommitteePower())
	assert.False(t, td.state.committee.Contains(secValKey.Address())) // Not in the committee

	// Committing another 10 blocks
	td.commitBlocks(t, 10)

	assert.True(t, td.state.evaluateSortition())                      // OK
	assert.False(t, td.state.committee.Contains(secValKey.Address())) // Still not in the committee

	td.commitBlocks(t, 1)

	assert.True(t, td.state.IsValidator(secValKey.Address()))
	assert.Equal(t, int64(1000000004), td.state.CommitteePower())
	assert.True(t, td.state.committee.Contains(secValKey.Address())) // In the committee
}

func TestValidateBlockTime(t *testing.T) {
	td := setup(t)

	t.Run("Time is not rounded", func(t *testing.T) {
		roundedNow := util.RoundNow(10)

		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(-15*time.Second)))
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(-5*time.Second)))
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(5*time.Second)))
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(15*time.Second)))
	})

	t.Run("Last block is committed 10 seconds ago", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(-10 * time.Second))

		// Before or same as the last block time
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(-20*time.Second)))
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(-10*time.Second)))

		// Ok
		assert.NoError(t, td.state.validateBlockTime(roundedNow))
		assert.NoError(t, td.state.validateBlockTime(roundedNow.Add(10*time.Second)))

		// More than the threshold
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(20*time.Second)))

		expectedProposeTime := roundedNow
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})

	t.Run("Last block is committed one minute ago", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(-1 * time.Minute)) // One minute ago
		lastBlockTime := td.state.LastBlockTime()

		// Before or same as the last block time
		assert.Error(t, td.state.validateBlockTime(lastBlockTime.Add(-10*time.Second)))
		assert.Error(t, td.state.validateBlockTime(lastBlockTime))

		// Ok
		assert.NoError(t, td.state.validateBlockTime(roundedNow.Add(-10*time.Second)))
		assert.NoError(t, td.state.validateBlockTime(roundedNow))
		assert.NoError(t, td.state.validateBlockTime(roundedNow.Add(10*time.Second)))

		// More than the threshold
		assert.Error(t, td.state.validateBlockTime(roundedNow.Add(30*time.Second)))

		expectedProposeTime := util.RoundNow(10)
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})

	t.Run("Last block is committed in future", func(t *testing.T) {
		roundedNow := util.RoundNow(10)
		td.state.lastInfo.UpdateBlockTime(roundedNow.Add(1 * time.Minute)) // One minute later
		lastBlockTime := td.state.LastBlockTime()

		assert.Error(t, td.state.validateBlockTime(lastBlockTime.Add(+1*time.Minute)))

		// Before the last block time
		assert.Error(t, td.state.validateBlockTime(lastBlockTime.Add(-10*time.Second)))
		assert.Error(t, td.state.validateBlockTime(lastBlockTime))

		// Ok
		assert.NoError(t, td.state.validateBlockTime(lastBlockTime.Add(10*time.Second)))
		assert.NoError(t, td.state.validateBlockTime(lastBlockTime.Add(20*time.Second)))

		// More than the threshold
		assert.Error(t, td.state.validateBlockTime(lastBlockTime.Add(30*time.Second)))

		expectedProposeTime := roundedNow.Add(1 * time.Minute).Add(
			time.Duration(td.state.params.BlockIntervalInSecond) * time.Second)
		assert.Equal(t, expectedProposeTime, td.state.proposeNextBlockTime())
	})
}

func TestValidatorHelpers(t *testing.T) {
	td := setup(t)

	t.Run("Should return nil for non-existing Validator Address", func(t *testing.T) {
		nonExistenceValidator := td.state.ValidatorByAddress(td.RandValAddress())
		assert.Nil(t, nonExistenceValidator, "State 1 returned non-nil For non-existing validator")
	})

	t.Run("Should return validator for valid committee Validator Address", func(t *testing.T) {
		existingValidator := td.state.ValidatorByAddress(td.genValKeys[0].Address())
		assert.NotNil(t, existingValidator)
		assert.Zero(t, existingValidator.Number())
	})

	t.Run("Should return validator for corresponding Validator number", func(t *testing.T) {
		existingValidator := td.state.ValidatorByNumber(0)
		assert.NotNil(t, existingValidator)
		assert.Zero(t, existingValidator.Number())
	})

	t.Run("Should return nil for invalid Validator number", func(t *testing.T) {
		nonExistenceValidator := td.state.ValidatorByNumber(10)
		assert.Nil(t, nonExistenceValidator)
	})
}

func TestLoadState(t *testing.T) {
	td := setup(t)

	// Add a bond transactions to change total power (stake)
	pub, _ := td.RandBLSKeyPair()
	lockTime := td.state.LastBlockHeight()
	bondTrx := tx.NewBondTx(lockTime, td.genAccKey.PublicKeyNative().AccountAddress(),
		pub.ValidatorAddress(), pub, 1000000000, 100000)
	td.HelperSignTransaction(td.genAccKey, bondTrx)

	assert.NoError(t, td.state.AddPendingTx(bondTrx))

	blk5, cert5 := td.makeBlockAndCertificate(t, 1)
	assert.Equal(t, 2, blk5.Transactions().Len())
	assert.NoError(t, td.state.CommitBlock(blk5, cert5))

	blk6, cert6 := td.makeBlockAndCertificate(t, 0)

	// Load last state info
	newState, err := LoadOrNewState(td.state.genDoc, td.state.valKeys,
		td.state.store, td.commonTxPool)
	require.NoError(t, err)

	assert.Equal(t, td.state.TotalAccounts(), newState.TotalAccounts())
	assert.Equal(t, td.state.TotalValidators(), newState.TotalValidators())
	assert.Equal(t, td.state.CommitteeValidators(), newState.CommitteeValidators())
	assert.Equal(t, td.state.CommitteePower(), newState.CommitteePower())
	assert.Equal(t, td.state.TotalPower(), newState.TotalPower())
	assert.Equal(t, td.state.Params(), newState.Params())
	assert.ElementsMatch(t, td.state.ValidatorAddresses(), newState.ValidatorAddresses())

	assert.Equal(t, int32(11), td.state.TotalAccounts()) // 9 subsidy addrs + 2 genesis addrs
	assert.Equal(t, int32(5), td.state.TotalValidators())

	// Try committing the next block
	require.NoError(t, newState.CommitBlock(blk6, cert6))
}

func TestIsValidator(t *testing.T) {
	td := setup(t)

	assert.True(t, td.state.IsInCommittee(td.genValKeys[0].Address()))
	assert.True(t, td.state.IsProposer(td.genValKeys[0].Address(), 0))
	assert.True(t, td.state.IsProposer(td.genValKeys[1].Address(), 1))
	assert.True(t, td.state.IsInCommittee(td.genValKeys[1].Address()))
	assert.True(t, td.state.IsValidator(td.genValKeys[1].Address()))

	addr := td.RandAccAddress()
	assert.False(t, td.state.IsInCommittee(addr))
	assert.False(t, td.state.IsProposer(addr, 0))
	assert.False(t, td.state.IsInCommittee(addr))
	assert.False(t, td.state.IsValidator(addr))
}

func TestCalculateFee(t *testing.T) {
	td := setup(t)

	fee := td.state.CalculateFee(td.RandAmount(), payload.TypeTransfer)
	expectedFee := td.commonTxPool.EstimatedFee(0, payload.TypeTransfer)

	assert.Equal(t, expectedFee, fee)
}

func TestCheckMaximumTransactionPerBlock(t *testing.T) {
	td := setup(t)

	td.state.params.MaxTransactionsPerBlock = 10
	lockTime := td.state.LastBlockHeight()
	senderAddr := td.genAccKey.PublicKeyNative().AccountAddress()
	for i := 0; i < td.state.params.MaxTransactionsPerBlock+2; i++ {
		amt := td.RandAmount()
		fee := td.state.CalculateFee(amt, payload.TypeTransfer)
		trx := tx.NewTransferTx(lockTime, senderAddr, td.RandAccAddress(), amt, fee)
		err := td.state.AddPendingTx(trx)
		assert.NoError(t, err)
	}

	blk, err := td.state.ProposeBlock(td.state.valKeys[0], td.RandAccAddress())
	assert.NoError(t, err)
	assert.Equal(t, td.state.params.MaxTransactionsPerBlock, blk.Transactions().Len())
}

func TestCommittedBlock(t *testing.T) {
	td := setup(t)

	t.Run("Genesis block", func(t *testing.T) {
		assert.Nil(t, td.state.CommittedBlock(0))
		assert.Equal(t, hash.UndefHash, td.state.BlockHash(0))
		assert.Equal(t, uint32(0), td.state.BlockHeight(hash.UndefHash))
	})

	t.Run("First block", func(t *testing.T) {
		cBlkOne := td.state.CommittedBlock(1)
		blkOne, err := cBlkOne.ToBlock()
		assert.NoError(t, err)
		assert.Nil(t, blkOne.PrevCertificate())
		assert.Equal(t, hash.UndefHash, blkOne.Header().PrevBlockHash())
	})

	t.Run("Last block", func(t *testing.T) {
		cBlkLast := td.state.CommittedBlock(td.state.LastBlockHeight())
		blkLast, err := cBlkLast.ToBlock()
		assert.NoError(t, err)
		assert.Equal(t, blkLast.Hash(), td.state.LastBlockHash())
	})
}
