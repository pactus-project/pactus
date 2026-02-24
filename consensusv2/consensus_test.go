package consensusv2

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ezex-io/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

const (
	tIndexX = 0
	tIndexY = 1
	tIndexB = 2
	tIndexP = 3
)

type consMessage struct {
	sender   crypto.Address
	receiver crypto.Address
	message  message.Message
}
type testData struct {
	*testsuite.TestSuite

	valKeys []*bls.ValidatorKey
	txPool  *txpool.MockTxPool
	genDoc  *genesis.Genesis
	consX   *consensusV2  // Good peer
	consY   *consensusV2  // Good peer
	consB   *consensusV2  // Byzantine or offline peer
	consP   *consensusV2  // Partitioned peer
	network []consMessage // Network messages
}

func testConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 1 * time.Hour, // Disabling timers
		ChangeProposerDelta:   1 * time.Hour, // Disabling timers
		QueryVoteTimeout:      1 * time.Hour, // Disabling timers
	}
}

func setup(t *testing.T) *testData {
	t.Helper()

	return setupWithSeed(t, testsuite.GenerateSeed())
}

func setupWithSeed(t *testing.T, seed int64) *testData {
	t.Helper()

	fmt.Printf("=== test %s, seed: %d\n", t.Name(), seed)

	ts := testsuite.NewTestSuiteFromSeed(t, seed)

	_, valKeys := ts.GenerateTestCommittee(4)
	txPool := txpool.MockingTxPool()

	vals := make([]*validator.Validator, 4)
	for i, key := range valKeys {
		val := validator.NewValidator(key.PublicKey(), int32(i))
		vals[i] = val
	}

	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	params := genesis.DefaultGenesisParams()
	params.CommitteeSize = 4

	// To prevent triggering timers before starting the tests and
	// avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockIntervalInSecond).
		Add(time.Duration(params.BlockIntervalInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)
	eventPipe := pipeline.New[any](t.Context())
	stateX, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexX]},
		store.MockingStore(ts), txPool, eventPipe)
	require.NoError(t, err)
	stateY, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexY]},
		store.MockingStore(ts), txPool, eventPipe)
	require.NoError(t, err)
	stateB, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexB]},
		store.MockingStore(ts), txPool, eventPipe)
	require.NoError(t, err)
	stateP, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexP]},
		store.MockingStore(ts), txPool, eventPipe)
	require.NoError(t, err)

	network := make([]consMessage, 0)
	td := &testData{
		TestSuite: ts,
		valKeys:   valKeys,
		txPool:    txPool,
		genDoc:    genDoc,
		network:   network,
	}
	broadcasterFunc := func(sender crypto.Address, msg message.Message) {
		for _, key := range valKeys {
			td.network = append(td.network,
				consMessage{
					sender:   sender,
					receiver: key.Address(),
					message:  msg,
				})
		}
	}
	td.consX = makeConsensus(t.Context(), testConfig(), stateX, valKeys[tIndexX],
		valKeys[tIndexX].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consY = makeConsensus(t.Context(), testConfig(), stateY, valKeys[tIndexY],
		valKeys[tIndexY].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consB = makeConsensus(t.Context(), testConfig(), stateB, valKeys[tIndexB],
		valKeys[tIndexB].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consP = makeConsensus(t.Context(), testConfig(), stateP, valKeys[tIndexP],
		valKeys[tIndexP].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())

	// -------------------------------
	// Better logging during testing
	overrideLogger := func(cons *consensusV2, name string) {
		cons.logger = logger.NewSubLogger("_consensus",
			testsuite.NewOverrideLogStringer(fmt.Sprintf("%s - %s: ", name, t.Name()), cons))
	}

	overrideLogger(td.consX, "consX")
	overrideLogger(td.consY, "consY")
	overrideLogger(td.consB, "consB")
	overrideLogger(td.consP, "consP")
	// -------------------------------

	logger.Info("setup finished, start running the test", "name", t.Name())

	return td
}

func (td *testData) shouldNotPublish(t *testing.T, cons *consensusV2, msgType message.Type) {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == msgType {
			require.Error(t, fmt.Errorf("should not publish %s", msgType))
		}
	}
}

func (td *testData) shouldPublishBlockAnnounce(t *testing.T, cons *consensusV2, hash hash.Hash) {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeBlockAnnounce {
			m := consMsg.message.(*message.BlockAnnounceMessage)
			assert.Equal(t, hash, m.Block.Hash())

			return
		}
	}
	require.NoError(t, errors.New("Block announce message not published"))
}

func (td *testData) shouldPublishProposal(t *testing.T, cons *consensusV2,
	height uint32, round int16,
) *proposal.Proposal {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeProposal {
			m := consMsg.message.(*message.ProposalMessage)
			require.Equal(t, height, m.Proposal.Height())
			require.Equal(t, round, m.Proposal.Round())

			return m.Proposal
		}
	}
	require.NoError(t, errors.New("Proposal message not published"))

	return nil
}

func (td *testData) shouldPublishQueryProposal(t *testing.T, cons *consensusV2, height uint32, round int16) {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender != cons.valKey.Address() ||
			consMsg.message.Type() != message.TypeQueryProposal {
			continue
		}

		m := consMsg.message.(*message.QueryProposalMessage)
		assert.Equal(t, m.Height, height)
		assert.Equal(t, m.Round, round)
		assert.Equal(t, m.Querier, cons.valKey.Address())

		return
	}
	require.NoError(t, errors.New("Query proposal message not published"))
}

func (td *testData) shouldPublishQueryVote(t *testing.T, cons *consensusV2, height uint32, round int16) {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender != cons.valKey.Address() ||
			consMsg.message.Type() != message.TypeQueryVote {
			continue
		}

		m := consMsg.message.(*message.QueryVoteMessage)
		assert.Equal(t, m.Height, height)
		assert.Equal(t, m.Round, round)
		assert.Equal(t, m.Querier, cons.valKey.Address())

		return
	}
	require.NoError(t, errors.New("Query proposal message not published"))
}

func (td *testData) shouldPublishVote(t *testing.T, cons *consensusV2, voteType vote.Type, hash hash.Hash) *vote.Vote {
	t.Helper()

	for i := len(td.network) - 1; i >= 0; i-- {
		consMsg := td.network[i]
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeVote {
			m := consMsg.message.(*message.VoteMessage)
			if m.Vote.Type() == voteType &&
				m.Vote.BlockHash() == hash {
				return m.Vote
			}
		}
	}
	require.NoError(t, errors.New("Vote message not published"))

	return nil
}

func (*testData) checkHeightRound(t *testing.T, cons *consensusV2, height uint32, round int16) {
	t.Helper()

	h, r := cons.HeightRound()
	assert.Equal(t, h, height)
	assert.Equal(t, r, round)
}

func (td *testData) addPrecommitVote(t *testing.T, cons *consensusV2, blockHash hash.Hash,
	height uint32, round int16, valID int,
) *vote.Vote {
	t.Helper()

	v := vote.NewPrecommitVote(blockHash, height, round, td.valKeys[valID].Address())

	return td.addVote(t, cons, v, valID)
}

func (td *testData) addCPPreVote(t *testing.T, cons *consensusV2, blockHash hash.Hash,
	height uint32, round int16, cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	t.Helper()

	v := vote.NewCPPreVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(t, cons, v, valID)
}

func (td *testData) addCPMainVote(t *testing.T, cons *consensusV2, blockHash hash.Hash,
	height uint32, round int16, cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	t.Helper()

	v := vote.NewCPMainVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(t, cons, v, valID)
}

func (td *testData) addCPDecidedVote(t *testing.T, cons *consensusV2, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	t.Helper()

	v := vote.NewCPDecidedVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(t, cons, v, valID)
}

func (td *testData) addVote(t *testing.T, cons *consensusV2, vote *vote.Vote, valID int) *vote.Vote {
	t.Helper()

	td.HelperSignVote(td.valKeys[valID], vote)
	cons.AddVote(vote)

	return vote
}

func (*testData) newHeightTimeout(cons *consensusV2) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

func (*testData) queryProposalTimeout(cons *consensusV2) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetQueryProposal})
	cons.lk.Unlock()
}

func (*testData) changeProposerTimeout(cons *consensusV2) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetChangeProposer})
	cons.lk.Unlock()
}

func (*testData) queryVoteTimeout(cons *consensusV2) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetQueryVote})
	cons.lk.Unlock()
}

// enterNewHeight helps tests to enter new height safely
// without scheduling new height. It boosts the test speed.
func (td *testData) enterNewHeight(cons *consensusV2) {
	cons.lk.Lock()
	cons.enterNewState(cons.newHeightState)
	cons.lk.Unlock()

	td.newHeightTimeout(cons)
}

// enterNextRound helps tests to enter next round safely.
func (*testData) enterNextRound(cons *consensusV2) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func (td *testData) commitBlockForAllStates(t *testing.T) (*block.Block, *certificate.Certificate) {
	t.Helper()

	height := td.consX.bcState.LastBlockHeight()
	var err error
	prop := td.makeProposal(t, height+1, 0)

	cert := certificate.NewCertificate(height+1, 0)
	sb := cert.SignBytesPrecommit(prop.Block().Hash())
	sigX := td.consX.valKey.Sign(sb)
	sigY := td.consY.valKey.Sign(sb)
	sigP := td.consP.valKey.Sign(sb)

	sig, _ := bls.SignatureAggregate(sigX, sigY, sigP)
	cert.SetSignature([]int32{tIndexX, tIndexY, tIndexB, tIndexP}, []int32{tIndexB}, sig)
	blk := prop.Block()

	err = td.consX.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consY.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consB.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consP.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)

	return blk, cert
}

// makeProposal generates a signed and valid proposal for the given height and round.
// If rewardAddr is provided, it will be used instead of the consensus instance's default reward address.
func (td *testData) makeProposal(t *testing.T, height uint32, round int16, rewardAddr ...crypto.Address,
) *proposal.Proposal {
	t.Helper()

	var cons *consensusV2
	switch (height % 4) + uint32(round%4) {
	case 1:
		cons = td.consX
	case 2:
		cons = td.consY
	case 3:
		cons = td.consB
	case 4, 0:
		cons = td.consP
	}

	// Use provided reward address or fall back to consensus instance's default
	addr := cons.rewardAddr
	if len(rewardAddr) > 0 {
		addr = rewardAddr[0]
	}

	blk, err := cons.bcState.ProposeBlock(cons.valKey, addr)
	require.NoError(t, err)
	p := proposal.NewProposal(height, round, blk)
	td.HelperSignProposal(cons.valKey, p)

	return p
}

// makeChangeProposerJusts generates justifications for changing the proposer at the specified height and round.
// If `proposal` is nil, it creates justifications for not changing the proposer;
// otherwise, it generates justifications to change the proposer.
// It returns three justifications:
//
//  1. `JustInitNo` if the proposal is set, or `JustInitYes` if not for the pre-vote step,
//  2. `JustMainVoteNoConflict` for the main-vote step,
//  3. `JustDecided` for the decided step.
func (td *testData) makeChangeProposerJusts(t *testing.T, propBlockHash hash.Hash,
	height uint32, round int16,
) (preVoteJust, mainVoteJust, decidedJust vote.Just) {
	t.Helper()

	cpRound := int16(0)

	// Create PreVote Justification
	var cpValue vote.CPValue

	if propBlockHash != hash.UndefHash {
		cpValue = vote.CPValueNo
		committers := make([]int32, 0, len(td.consP.validators))
		sigs := make([]*bls.Signature, 0, len(td.consP.validators))
		for i, val := range td.consP.validators {
			vote := vote.NewPrecommitVote(propBlockHash, height, round, val.Address())
			signBytes := vote.SignBytes()

			committers = append(committers, val.Number())
			sigs = append(sigs, td.valKeys[i].Sign(signBytes))
		}
		aggSig, _ := bls.SignatureAggregate(sigs...)
		cert := certificate.NewCertificate(height, round)
		cert.SetSignature(committers, []int32{}, aggSig)

		preVoteJust = &vote.JustInitNo{
			QCert: cert,
		}
	} else {
		cpValue = vote.CPValueYes
		preVoteJust = &vote.JustInitYes{}
	}

	// Create MainVote Justification
	preVoteCommitters := make([]int32, 0, len(td.consP.validators))
	preVoteSigs := make([]*bls.Signature, 0, len(td.consP.validators))
	for i, val := range td.consP.validators {
		preVote := vote.NewCPPreVote(propBlockHash, height, round,
			cpRound, cpValue, preVoteJust, val.Address())
		signBytes := preVote.SignBytes()

		preVoteCommitters = append(preVoteCommitters, val.Number())
		preVoteSigs = append(preVoteSigs, td.valKeys[i].Sign(signBytes))
	}
	preVoteAggSig, _ := bls.SignatureAggregate(preVoteSigs...)
	certPreVote := certificate.NewCertificate(height, round)
	certPreVote.SetSignature(preVoteCommitters, []int32{}, preVoteAggSig)
	mainVoteJust = &vote.JustMainVoteNoConflict{QCert: certPreVote}

	// Create Decided Justification
	mainVoteCommitters := make([]int32, 0, len(td.consP.validators))
	mainVoteSigs := make([]*bls.Signature, 0, len(td.consP.validators))
	for i, val := range td.consP.validators {
		mainVote := vote.NewCPMainVote(propBlockHash, height, round,
			cpRound, cpValue, mainVoteJust, val.Address())
		signBytes := mainVote.SignBytes()

		mainVoteCommitters = append(mainVoteCommitters, val.Number())
		mainVoteSigs = append(mainVoteSigs, td.valKeys[i].Sign(signBytes))
	}
	mainVoteAggSig, _ := bls.SignatureAggregate(mainVoteSigs...)
	certMainVote := certificate.NewCertificate(height, round)
	certMainVote.SetSignature(mainVoteCommitters, []int32{}, mainVoteAggSig)
	decidedJust = &vote.JustDecided{QCert: certMainVote}

	return preVoteJust, mainVoteJust, decidedJust
}

func TestStart(t *testing.T) {
	td := setup(t)

	td.consX.MoveToNewHeight()

	td.checkHeightRound(t, td.consX, 1, 0)
}

func TestNotInCommittee(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	str := store.MockingStore(td.TestSuite)

	state, _ := state.LoadOrNewState(td.genDoc, []*bls.ValidatorKey{valKey}, str, td.txPool, nil)
	pipe := pipeline.New[message.Message](t.Context())
	consInst := NewConsensus(t.Context(), testConfig(), state, valKey, valKey.Address(), pipe,
		newConcreteMediator())
	cons := consInst.(*consensusV2)

	td.enterNewHeight(cons)
	td.newHeightTimeout(cons)
	assert.Equal(t, cons.currentState.name(), "new-height")
}

func TestIsProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	assert.False(t, td.consX.IsProposer())
	assert.True(t, td.consY.IsProposer())
}

func TestVoteWithInvalidHeight(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1
	td.enterNewHeight(td.consP)

	vote1 := td.addPrecommitVote(t, td.consP, td.RandHash(), 1, 0, tIndexX)
	vote2 := td.addPrecommitVote(t, td.consP, td.RandHash(), 2, 0, tIndexX)
	vote3 := td.addPrecommitVote(t, td.consP, td.RandHash(), 2, 0, tIndexY)
	vote4 := td.addPrecommitVote(t, td.consP, td.RandHash(), 3, 0, tIndexX)

	require.False(t, td.consP.HasVote(vote1.Hash()))
	require.True(t, td.consP.HasVote(vote2.Hash()))
	require.True(t, td.consP.HasVote(vote3.Hash()))
	require.False(t, td.consP.HasVote(vote4.Hash()))
}

func TestConsensusAbsoluteCommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 2, 0)

	prop := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(prop)

	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, prop.Block().Hash())
	td.addPrecommitVote(t, td.consX, prop.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(t, td.consX, prop.Block().Hash(), 2, 0, tIndexB)
	td.addPrecommitVote(t, td.consX, prop.Block().Hash(), 2, 0, tIndexP)

	td.shouldPublishBlockAnnounce(t, td.consX, prop.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	vote1 := td.addPrecommitVote(t, td.consP, td.RandHash(), 1, 0, tIndexX)
	vote2 := td.addPrecommitVote(t, td.consP, td.RandHash(), 1, 2, tIndexX)
	vote3 := td.addPrecommitVote(t, td.consP, td.RandHash(), 1, 1, tIndexX)
	vote4 := td.addPrecommitVote(t, td.consP, td.RandHash(), 2, 0, tIndexX)
	vote5, _ := td.GenerateTestPrecommitVote(1, 0)
	td.consP.AddVote(vote5)

	assert.False(t, td.consP.HasVote(vote1.Hash())) // previous round
	assert.True(t, td.consP.HasVote(vote2.Hash()))  // next round
	assert.True(t, td.consP.HasVote(vote3.Hash()))
	assert.False(t, td.consP.HasVote(vote4.Hash())) // valid votes for the next height
	assert.False(t, td.consP.HasVote(vote5.Hash())) // invalid votes

	assert.Equal(t, []*vote.Vote{vote3}, td.consP.AllVotes())
	assert.NotContains(t, td.consP.AllVotes(), vote2)
}

// TestConsensusDelayedProposal tests the scenario where a node receives votes
// before receiving the proposal due to network delays.
func TestConsensusDelayedProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)
	blockHash := prop.Block().Hash()

	// consP receives other votes first
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexX)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexY)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexB)

	// consP receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, blockHash)
	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
}

// TestConsensusDelayedVote tests the scenario where a node receives votes
// after timing out due to network delays.
func TestConsensusDelayedVote(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)
	blockHash := prop.Block().Hash()

	td.consP.SetProposal(prop)
	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, blockHash)

	// consP moves to change proposer state
	td.changeProposerTimeout(td.consP)

	// consP receives other votes now
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexX)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexY)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexB)

	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
}

func TestHandleQueryVote(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	height := uint32(2)

	assert.Nil(t, td.consP.HandleQueryVote(height, 0))

	// round 0
	preVoteJust, mainVoteJust, decidedJust := td.makeChangeProposerJusts(t, hash.UndefHash, height, 0)
	r0Vote1 := td.addCPPreVote(t, td.consP, hash.UndefHash, height, 0, vote.CPValueYes, preVoteJust, tIndexY)
	r0Vote2 := td.addCPMainVote(t, td.consP, hash.UndefHash, height, 0, vote.CPValueYes, mainVoteJust, tIndexY)
	r0Vote3 := td.addCPDecidedVote(t, td.consP, hash.UndefHash, height, 0, vote.CPValueYes, decidedJust, tIndexY)

	// round 1
	td.enterNextRound(td.consP)

	hash := td.RandHash()
	preVoteJust, mainVoteJust, decidedJust = td.makeChangeProposerJusts(t, hash, height, 1)
	r1Vote1 := td.addPrecommitVote(t, td.consP, td.RandHash(), height, 1, tIndexY)
	r1Vote2 := td.addCPPreVote(t, td.consP, hash, height, 1, vote.CPValueNo, preVoteJust, tIndexY)
	r1Vote3 := td.addCPMainVote(t, td.consP, hash, height, 1, vote.CPValueNo, mainVoteJust, tIndexY)
	r1Vote4 := td.addCPDecidedVote(t, td.consP, hash, height, 1, vote.CPValueNo, decidedJust, tIndexY)

	// Round 2
	td.enterNextRound(td.consP)

	td.addPrecommitVote(t, td.consP, td.RandHash(), height, 2, tIndexY)

	require.True(t, td.consP.HasVote(r0Vote1.Hash()))
	require.True(t, td.consP.HasVote(r0Vote2.Hash()))
	require.True(t, td.consP.HasVote(r0Vote3.Hash()))
	require.True(t, td.consP.HasVote(r1Vote1.Hash()))
	require.True(t, td.consP.HasVote(r1Vote2.Hash()))
	require.True(t, td.consP.HasVote(r1Vote3.Hash()))
	require.True(t, td.consP.HasVote(r1Vote4.Hash()))

	rndVote0 := td.consP.HandleQueryVote(height, 0)
	assert.Equal(t, r0Vote3, rndVote0, "should send the decided vote for the round 0")

	rndVote1 := td.consP.HandleQueryVote(height, 1)
	assert.Equal(t, r1Vote4, rndVote1, "should send the decided vote for the round 1")

	rndVote2 := td.consP.HandleQueryVote(height, 2)
	assert.Equal(t, int16(2), rndVote2.Round(), "should send the precommit vote for the current round")

	rndVote3 := td.consP.HandleQueryVote(height, 3)
	assert.Nil(t, rndVote3, "should not send a vote for the next round")

	rndVote4 := td.consP.HandleQueryVote(height+1, 0)
	assert.Nil(t, rndVote4, "should not have a vote for the next height")
}

func TestHandleQueryProposal(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	// Round 1
	td.enterNextRound(td.consX)
	td.enterNextRound(td.consY) // consY is the proposer
	td.consX.SetProposal(td.consY.Proposal())

	height := uint32(1)
	round := int16(1)

	prop0 := td.consY.HandleQueryProposal(height, round-1)
	assert.Nil(t, prop0, "proposer should not send a proposal for the previous round")

	prop1 := td.consX.HandleQueryProposal(height, round)
	assert.Nil(t, prop1, "non-proposer should not send a proposal")

	prop2 := td.consY.HandleQueryProposal(height, round)
	assert.NotNil(t, prop2, "proposer should send a proposal")

	td.consX.cpDecidedCert = td.GenerateTestCertificate(1) // TODO: better way?
	prop3 := td.consX.HandleQueryProposal(height, round)
	assert.NotNil(t, prop3, "non-proposer should send a proposal on decided proposal")

	prop4 := td.consX.HandleQueryProposal(height+1, 0)
	assert.Nil(t, prop4, "should not have a proposal for the next height")
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// It should ignore proposals for previous rounds
	prop := td.makeProposal(t, 1, 0)
	td.consP.SetProposal(prop)

	assert.Nil(t, td.consP.Proposal())
	td.checkHeightRound(t, td.consP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	prop := td.makeProposal(t, 1, 0)
	td.consP.SetProposal(prop)

	assert.Nil(t, td.consP.Proposal())
	td.checkHeightRound(t, td.consP, 2, 0)
}

func TestDoubleProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	height := uint32(4)
	round := int16(0)
	prop1 := td.makeProposal(t, height, round, td.RandAccAddress())
	prop2 := td.makeProposal(t, height, round, td.RandAccAddress())
	assert.NotEqual(t, prop1.Hash(), prop2.Hash())

	td.consX.SetProposal(prop1)
	td.consX.SetProposal(prop2)

	assert.Equal(t, td.consX.Proposal().Hash(), prop1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	pipe := pipeline.New[message.Message](t.Context())
	consInst := NewConsensus(t.Context(), testConfig(), state.MockingState(td.TestSuite),
		valKey, valKey.Address(), pipe, newConcreteMediator())
	nonActiveCons := consInst.(*consensusV2)

	t.Run("non-active instances should be in new-height state", func(t *testing.T) {
		nonActiveCons.MoveToNewHeight()
		td.newHeightTimeout(nonActiveCons)
		td.checkHeightRound(t, nonActiveCons, 1, 0)

		// Double entry
		nonActiveCons.MoveToNewHeight()
		td.newHeightTimeout(nonActiveCons)
		td.checkHeightRound(t, nonActiveCons, 1, 0)

		assert.False(t, nonActiveCons.IsActive())
		assert.Equal(t, nonActiveCons.currentState.name(), "new-height")
	})

	t.Run("non-active instances should ignore proposals", func(t *testing.T) {
		prop := td.makeProposal(t, 1, 0)
		nonActiveCons.SetProposal(prop)

		assert.Nil(t, nonActiveCons.Proposal())
	})

	t.Run("non-active instances should ignore votes", func(t *testing.T) {
		v := td.addPrecommitVote(t, nonActiveCons, td.RandHash(), 1, 0, tIndexX)

		assert.False(t, nonActiveCons.HasVote(v.Hash()))
	})

	t.Run("non-active instances should move to new height", func(t *testing.T) {
		b1, cert1 := td.commitBlockForAllStates(t)

		nonActiveCons.MoveToNewHeight()
		td.checkHeightRound(t, nonActiveCons, 1, 0)

		assert.NoError(t, nonActiveCons.bcState.CommitBlock(b1, cert1))

		nonActiveCons.MoveToNewHeight()
		td.newHeightTimeout(nonActiveCons)
		td.checkHeightRound(t, nonActiveCons, 2, 0)
	})
}

func TestVoteWithBigRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)

	vote := td.addPrecommitVote(t, td.consX, td.RandHash(), 1, util.MaxInt16, tIndexB)
	assert.True(t, td.consX.HasVote(vote.Hash()))
}

func TestProposalWithBigRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	prop := td.makeProposal(t, 1, util.MaxInt16)
	td.consP.SetProposal(prop)
	assert.Nil(t, td.consP.Proposal())
}

func TestInvalidProposal(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	prop := td.makeProposal(t, 1, 0)
	prop.SetSignature(nil) // Make proposal invalid
	td.consP.SetProposal(prop)
	assert.Nil(t, td.consP.Proposal())
}

// TestCases runs some special cases to test the consensus algorithm.
func TestCasesNormal(t *testing.T) {
	tests := []struct {
		seed        int64
		certRound   int16
		description string
	}{
		{1758015756630707317, 1, "precommit: startChangingProposer on 1f+1 pre-votes"},
		{1758013793525473037, 1, "cp_prevote: has one main-vote for `yes` from previous round"},
		{1758014151948449992, 0, "cp_prevote: all main-votes are `abstain` from previous round"},
		{1758014780626270377, 2, "cp_mainvote: has 2f+1 pre-votes for `no`, decided on `no (biased)`"},
		{1758014728375405734, 1, "cp_mainvote: has 2f+1 pre-votes for `yes`"},
		{1758014840359554900, 1, "cp_mainvote: has no pre-votes quorum"},
		{1758015210671486317, 2, "cp_decide: decide on `yes`"},
		{1758015258023425309, 0, "cp_decide: conflicting main-votes"},
		{1758015607471650150, 0, "cons.cpRound = 1, decided vote for `yes` in cpRound 0"},
	}

	for _, tt := range tests {
		td := setupWithSeed(t, tt.seed)
		cert := td.executeConsensusNormal(t)

		require.Equal(t, tt.certRound, cert.Round(),
			"test '%s' failed. round not matched (expected %d, got %d)",
			tt.description, tt.certRound, cert.Round())
	}
}

// TestConsensusNormal runs multiple iterations of the consensus process to
// ensure stability and correctness across repeated executions.
// If the consensus process fails in any iteration, the test will fail.
// A random seed is generated for each iteration, allowing the test to be
// reproduced deterministically with same seed.
// This test doesn't involve any Byzantine behavior or network partitions.
func TestConsensusNormal(t *testing.T) {
	for i := 0; i < 10; i++ {
		td := setup(t)
		td.executeConsensusNormal(t)
	}
}

func (td *testData) executeConsensusNormal(t *testing.T) *certificate.Certificate {
	t.Helper()

	td.commitBlockForAllStates(t)
	height := uint32(2)

	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)
	td.enterNewHeight(td.consB)
	td.enterNewHeight(td.consP)

	cert, err := executeConsensus(td, td.RandBool())
	require.NoError(t, err)
	require.Equal(t, cert.Height(), height)

	return cert
}

func TestCasesByzantine(t *testing.T) {
	tests := []struct {
		seed        int64
		certRound   int16
		description string
	}{
		{1758019943838125552, 0, "double proposal detected"},
	}

	for _, tt := range tests {
		td := setupWithSeed(t, tt.seed)
		cert := td.executeConsensusByzantine(t)

		require.Equal(t, tt.certRound, cert.Round(),
			"test '%s' failed. round not matched (expected %d, got %d)",
			tt.description, tt.certRound, cert.Round())
	}
}

// TestConsensusByzantine tests a scenario involving a Byzantine node and a network partition.
//
// We have four nodes: X, Y, B, and P, which:
// - B is a Byzantine node
// - X, Y, and P are honest nodes
// - However, P is partitioned and perceives the network through B.
//
// Byzantine node B acts maliciously by double proposing:
// sending one proposal to X and Y, and another proposal to P.
// Consensus halts because X and Y cannot reach consensus without P,
// and P cannot reach consensus without X and Y.
//
// The network partition is healed after some time.
// Once the partition is healed, honest nodes can reach consensus.
//
// The Byzantine node B acts maliciously by double proposing and double voting in this test.
func TestConsensusByzantine(t *testing.T) {
	for i := 0; i < 10; i++ {
		td := setup(t)
		td.executeConsensusByzantine(t)
	}
}

func (td *testData) executeConsensusByzantine(t *testing.T) *certificate.Certificate {
	t.Helper()

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	height := uint32(3)
	round := int16(0)

	// =================================
	// Byzantine node B
	prop1 := td.makeProposal(t, height, round, td.RandAccAddress())
	prop2 := td.makeProposal(t, height, round, td.RandAccAddress())

	require.NotEqual(t, prop1.Block().Hash(), prop2.Block().Hash())
	require.Equal(t, td.consB.valKey.Address(), prop1.Block().Header().ProposerAddress())
	require.Equal(t, td.consB.valKey.Address(), prop2.Block().Header().ProposerAddress())

	// =================================
	// Honest nodes X, Y
	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	td.consX.SetProposal(prop1)
	td.consY.SetProposal(prop1)

	voteX := td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, prop1.Block().Hash())
	voteY := td.shouldPublishVote(t, td.consY, vote.VoteTypePrecommit, prop1.Block().Hash())

	// Byzantine node partitioned the network and blocked the Node P.
	// X and Y request to change proposer

	td.changeProposerTimeout(td.consX)
	td.changeProposerTimeout(td.consY)

	// X and Y are unable to progress...

	// =================================
	// Honest node P, partitioned node
	td.enterNewHeight(td.consP)

	// P receives the second proposal
	td.consP.SetProposal(prop2)

	voteP := td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop2.Block().Hash())

	// Request to change proposer
	td.changeProposerTimeout(td.consP)

	// P is unable to progress...

	// =================================
	// Byzantine node B votes on both proposals
	byzVote1 := vote.NewPrecommitVote(prop1.Block().Hash(), height, round, td.consB.valKey.Address())
	byzVote2 := vote.NewPrecommitVote(prop2.Block().Hash(), height, round, td.consB.valKey.Address())
	td.HelperSignVote(td.consB.valKey, byzVote1)
	td.HelperSignVote(td.consB.valKey, byzVote2)

	aggSig1, _ := bls.SignatureAggregate(voteX.Signature(), voteY.Signature())
	aggSig2, _ := bls.SignatureAggregate(voteP.Signature(), byzVote1.Signature())
	cert1 := certificate.NewCertificate(height, round)
	cert1.SetSignature([]int32{tIndexX, tIndexY, tIndexB, tIndexP}, []int32{tIndexB, tIndexP}, aggSig1)

	cert2 := certificate.NewCertificate(height, round)
	cert2.SetSignature([]int32{tIndexX, tIndexY, tIndexB, tIndexP}, []int32{tIndexX, tIndexY}, aggSig2)

	byzVote3 := vote.NewCPPreVote(prop1.Block().Hash(), height, round, 0, vote.CPValueNo,
		&vote.JustInitNo{QCert: cert1}, td.consB.valKey.Address())
	byzVote4 := vote.NewCPPreVote(prop2.Block().Hash(), height, round, 0, vote.CPValueNo,
		&vote.JustInitNo{QCert: cert2}, td.consB.valKey.Address())

	td.HelperSignVote(td.consB.valKey, byzVote2)
	td.HelperSignVote(td.consB.valKey, byzVote3)

	td.consB.broadcastVote(byzVote1)
	td.consB.broadcastVote(byzVote2)
	td.consB.broadcastVote(byzVote3)
	td.consB.broadcastVote(byzVote4)

	// =================================
	// Now, Partition heals
	fmt.Println("== Partition heals")
	td.checkHeightRound(t, td.consX, height, round)
	td.checkHeightRound(t, td.consY, height, round)
	td.checkHeightRound(t, td.consP, height, round)

	cert, err := executeConsensus(td, true)
	require.NoError(t, err)
	require.Equal(t, cert.Height(), height)

	return cert
}

// executeConsensus runs the consensus algorithm until it reaches consensus or the network is exhausted.
// If `withoutByzantineNode` is true, it simulates the scenario where the Byzantine node goes offline.
// It returns the certificate of the committed block if consensus is reached,
// or an error if consensus is violated or cannot be reached.
func executeConsensus(td *testData, withoutByzantineNode bool) (
	*certificate.Certificate, error,
) {
	instances := []*consensusV2{td.consX, td.consY, td.consB, td.consP}

	if withoutByzantineNode {
		// remove byzantine node (Byzantine node goes offline)
		instances = []*consensusV2{td.consX, td.consY, td.consP}
	}

	// 50% chance for the first proposal to be lost,
	// then decrease the chance by 5% for each iteration.
	changeProposerChance := 50

	blockAnnounces := map[crypto.Address]*message.BlockAnnounceMessage{}
	for len(td.network) > 0 {
		rndIndex := td.RandIntMax(len(td.network))
		rndMsg := td.network[rndIndex]
		td.network = slices.Delete(td.network, rndIndex, rndIndex+1)

		switch rndMsg.message.Type() {
		case message.TypeVote:
			m := rndMsg.message.(*message.VoteMessage)
			for _, cons := range instances {
				if cons.valKey.Address() == rndMsg.receiver {
					cons.AddVote(m.Vote)
				}
			}

		case message.TypeProposal:
			m := rndMsg.message.(*message.ProposalMessage)
			for _, cons := range instances {
				if cons.valKey.Address() == rndMsg.receiver {
					cons.SetProposal(m.Proposal)
				}
			}

		case message.TypeQueryProposal:
			m := rndMsg.message.(*message.QueryProposalMessage)
			for _, cons := range instances {
				p := cons.HandleQueryProposal(m.Height, m.Round)
				if p != nil {
					td.network = append(td.network, consMessage{
						sender:  cons.valKey.Address(),
						message: message.NewProposalMessage(p),
					})
				}
			}

		case message.TypeQueryVote:
			// To make the test reproducible, we ignore the QueryVote message.
			// This is because QueryVote returns a random vote that can make the test non-reproducible.

		case message.TypeBlockAnnounce:
			m := rndMsg.message.(*message.BlockAnnounceMessage)
			blockAnnounces[rndMsg.sender] = m

		case
			message.TypeHello,
			message.TypeHelloAck,
			message.TypeTransaction,
			message.TypeBlocksRequest,
			message.TypeBlocksResponse:
			//
		}

		for _, cons := range instances {
			rnd := td.RandIntMax(100)
			if rnd < changeProposerChance ||
				len(td.network) == 0 {
				td.changeProposerTimeout(cons)
			}
		}
		changeProposerChance -= 5
	}

	// Verify whether more than (1f+1) nodes have committed to the same block.
	if len(blockAnnounces) >= 2 {
		var firstAnnounce *message.BlockAnnounceMessage
		for _, msg := range blockAnnounces {
			if firstAnnounce == nil {
				firstAnnounce = msg
			} else if msg.Block.Hash() != firstAnnounce.Block.Hash() {
				return nil, fmt.Errorf("consensus violated, seed %v", td.TestSuite.Seed)
			}
		}

		// everything is ok
		return firstAnnounce.Certificate, nil
	}

	return nil, fmt.Errorf("unable to reach consensus, seed %v", td.TestSuite.Seed)
}
