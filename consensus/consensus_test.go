package consensus

import (
	"errors"
	"fmt"
	"testing"
	"time"

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
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
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
	sender  crypto.Address
	message message.Message
}
type testData struct {
	*testsuite.TestSuite

	valKeys      []*bls.ValidatorKey
	txPool       *txpool.MockTxPool
	genDoc       *genesis.Genesis
	consX        *consensus // Good peer
	consY        *consensus // Good peer
	consB        *consensus // Byzantine or offline peer
	consP        *consensus // Partitioned peer
	consMessages []consMessage
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

	ts := testsuite.NewTestSuiteFromSeed(seed)

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
	eventPipe := pipeline.MockingPipeline[any]()
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

	consMessages := make([]consMessage, 0)
	td := &testData{
		TestSuite:    ts,
		valKeys:      valKeys,
		txPool:       txPool,
		genDoc:       genDoc,
		consMessages: consMessages,
	}
	broadcasterFunc := func(sender crypto.Address, msg message.Message) {
		fmt.Printf("received a message %s: %s\n", msg.Type(), msg.String())
		td.consMessages = append(td.consMessages, consMessage{
			sender:  sender,
			message: msg,
		})
	}
	td.consX = makeConsensus(testConfig(), stateX, valKeys[tIndexX],
		valKeys[tIndexX].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consY = makeConsensus(testConfig(), stateY, valKeys[tIndexY],
		valKeys[tIndexY].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consB = makeConsensus(testConfig(), stateB, valKeys[tIndexB],
		valKeys[tIndexB].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consP = makeConsensus(testConfig(), stateP, valKeys[tIndexP],
		valKeys[tIndexP].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())

	// -------------------------------
	// Better logging during testing
	overrideLogger := func(cons *consensus, name string) {
		cons.logger = logger.NewSubLogger("_consensus",
			testsuite.NewOverrideStringer(fmt.Sprintf("%s - %s: ", name, t.Name()), cons))
	}

	overrideLogger(td.consX, "consX")
	overrideLogger(td.consY, "consY")
	overrideLogger(td.consB, "consB")
	overrideLogger(td.consP, "consP")
	// -------------------------------

	logger.Info("setup finished, start running the test", "name", t.Name())

	return td
}

func (td *testData) shouldPublishBlockAnnounce(t *testing.T, cons *consensus, hash hash.Hash) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeBlockAnnounce {
			m := consMsg.message.(*message.BlockAnnounceMessage)
			assert.Equal(t, hash, m.Block.Hash())

			return
		}
	}
	require.NoError(t, errors.New("Not found"))
}

func (td *testData) shouldPublishProposal(t *testing.T, cons *consensus,
	height uint32, round int16,
) *proposal.Proposal {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeProposal {
			m := consMsg.message.(*message.ProposalMessage)
			require.Equal(t, height, m.Proposal.Height())
			require.Equal(t, round, m.Proposal.Round())

			return m.Proposal
		}
	}

	return nil
}

func (td *testData) shouldNotPublish(t *testing.T, cons *consensus, msgType message.Type) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == msgType {
			require.Error(t, fmt.Errorf("should not public %s", msgType))
		}
	}
}

func (td *testData) shouldPublishQueryProposal(t *testing.T, cons *consensus, height uint32) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender != cons.valKey.Address() ||
			consMsg.message.Type() != message.TypeQueryProposal {
			continue
		}

		m := consMsg.message.(*message.QueryProposalMessage)
		assert.Equal(t, m.Height, height)
		assert.Equal(t, m.Querier, cons.valKey.Address())

		return
	}
}

func (td *testData) shouldPublishQueryVote(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	for _, consMsg := range td.consMessages {
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
}

func (td *testData) shouldPublishVote(t *testing.T, cons *consensus, voteType vote.Type, hash hash.Hash) *vote.Vote {
	t.Helper()

	for i := len(td.consMessages) - 1; i >= 0; i-- {
		consMsg := td.consMessages[i]
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeVote {
			m := consMsg.message.(*message.VoteMessage)
			if m.Vote.Type() == voteType &&
				m.Vote.BlockHash() == hash {
				return m.Vote
			}
		}
	}

	return nil
}

func (*testData) checkHeightRound(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	h, r := cons.HeightRound()
	assert.Equal(t, h, height)
	assert.Equal(t, r, round)
}

func (td *testData) addPrepareVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	valID int,
) *vote.Vote {
	v := vote.NewPrepareVote(blockHash, height, round, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addPrecommitVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	valID int,
) *vote.Vote {
	v := vote.NewPrecommitVote(blockHash, height, round, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPPreVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPPreVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addCPMainVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPMainVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addCPDecidedVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPDecidedVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addVote(cons *consensus, v *vote.Vote, valID int) *vote.Vote {
	td.HelperSignVote(td.valKeys[valID], v)
	cons.AddVote(v)

	return v
}

func (*testData) newHeightTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

func (*testData) queryProposalTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetQueryProposal})
	cons.lk.Unlock()
}

func (*testData) changeProposerTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{time.Hour, cons.height, cons.round, tickerTargetChangeProposer})
	cons.lk.Unlock()
}

// enterNewHeight helps tests to enter new height safely
// without scheduling new height. It boosts the test speed.
func (td *testData) enterNewHeight(cons *consensus) {
	cons.lk.Lock()
	cons.enterNewState(cons.newHeightState)
	cons.lk.Unlock()

	td.newHeightTimeout(cons)
}

// enterNextRound helps tests to enter next round safely.
func (*testData) enterNextRound(cons *consensus) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func (td *testData) commitBlockForAllStates(t *testing.T) (*block.Block, *certificate.BlockCertificate) {
	t.Helper()

	height := td.consX.bcState.LastBlockHeight()
	var err error
	prop := td.makeProposal(t, height+1, 0)

	cert := certificate.NewBlockCertificate(height+1, 0)
	signBytes := cert.SignBytes(prop.Block().Hash())
	sig1 := td.consX.valKey.Sign(signBytes)
	sig2 := td.consY.valKey.Sign(signBytes)
	sig3 := td.consB.valKey.Sign(signBytes)
	sig4 := td.consP.valKey.Sign(signBytes)

	sig := bls.SignatureAggregate(sig1, sig2, sig3, sig4)
	cert.SetSignature([]int32{tIndexX, tIndexY, tIndexB, tIndexP}, []int32{}, sig)
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

func (td *testData) makeProposal(t *testing.T, height uint32, round int16) *proposal.Proposal {
	t.Helper()

	var prop *proposal.Proposal
	switch (height % 4) + uint32(round%4) {
	case 1:
		blk, err := td.consX.bcState.ProposeBlock(td.consX.valKey, td.consX.rewardAddr)
		require.NoError(t, err)
		prop = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consX.valKey, prop)
	case 2:
		blk, err := td.consY.bcState.ProposeBlock(td.consY.valKey, td.consY.rewardAddr)
		require.NoError(t, err)
		prop = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consY.valKey, prop)
	case 3:
		blk, err := td.consB.bcState.ProposeBlock(td.consB.valKey, td.consB.rewardAddr)
		require.NoError(t, err)
		prop = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consB.valKey, prop)
	case 0, 4:
		blk, err := td.consP.bcState.ProposeBlock(td.consP.valKey, td.consP.rewardAddr)
		require.NoError(t, err)
		prop = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consP.valKey, prop)
	}

	return prop
}

func (td *testData) makeMainVoteCertificate(t *testing.T,
	height uint32, round, cpRound int16,
) *certificate.VoteCertificate {
	t.Helper()

	// === make valid certificate
	preVoteCommitters := []int32{}
	preVoteSigs := []*bls.Signature{}
	for i, val := range td.consP.validators {
		preVoteJust := &vote.JustInitYes{}
		preVote := vote.NewCPPreVote(hash.UndefHash, height, round,
			cpRound, vote.CPValueYes, preVoteJust, val.Address())
		sbPreVote := preVote.SignBytes()

		preVoteCommitters = append(preVoteCommitters, val.Number())
		preVoteSigs = append(preVoteSigs, td.valKeys[i].Sign(sbPreVote))
	}
	preVoteAggSig := bls.SignatureAggregate(preVoteSigs...)
	certPreVote := certificate.NewVoteCertificate(height, round)
	certPreVote.SetSignature(preVoteCommitters, []int32{}, preVoteAggSig)

	mainVoteCommitters := []int32{}
	mainVoteSigs := []*bls.Signature{}
	for i, val := range td.consP.validators {
		mainVoteJust := &vote.JustMainVoteNoConflict{
			QCert: certPreVote,
		}
		mainVote := vote.NewCPMainVote(hash.UndefHash, height, round, cpRound, vote.CPValueYes, mainVoteJust, val.Address())
		sbMainVote := mainVote.SignBytes()

		mainVoteCommitters = append(mainVoteCommitters, val.Number())
		mainVoteSigs = append(mainVoteSigs, td.valKeys[i].Sign(sbMainVote))
	}
	mainVoteAggSig := bls.SignatureAggregate(mainVoteSigs...)
	certMainVote := certificate.NewVoteCertificate(height, round)
	certMainVote.SetSignature(mainVoteCommitters, []int32{}, mainVoteAggSig)

	return certMainVote
}

func TestStart(t *testing.T) {
	td := setup(t)

	td.consX.MoveToNewHeight()
	td.checkHeightRound(t, td.consX, 1, 0)
}

func TestNotInCommittee(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()

	state := state.MockingState(td.TestSuite)
	pipe := pipeline.MockingPipeline[message.Message]()
	consInt := NewConsensus(testConfig(), state, valKey,
		valKey.Address(), pipe, newConcreteMediator())
	cons := consInt.(*consensus)

	td.enterNewHeight(cons)
	td.newHeightTimeout(cons)
	assert.Equal(t, "new-height", cons.currentState.name())
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

	v1 := td.addPrepareVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	v2 := td.addPrepareVote(td.consP, td.RandHash(), 2, 0, tIndexX)
	v3 := td.addPrepareVote(td.consP, td.RandHash(), 2, 0, tIndexY)
	v4 := td.addPrepareVote(td.consP, td.RandHash(), 3, 0, tIndexX)

	require.False(t, td.consP.HasVote(v1.Hash()))
	require.True(t, td.consP.HasVote(v2.Hash()))
	require.True(t, td.consP.HasVote(v3.Hash()))
	require.False(t, td.consP.HasVote(v4.Hash()))
}

func TestConsensusNormalCase(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 2, 0)

	prop := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(prop)

	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, prop.Block().Hash())

	td.addPrecommitVote(td.consX, prop.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(td.consX, prop.Block().Hash(), 2, 0, tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, prop.Block().Hash())

	td.shouldPublishBlockAnnounce(t, td.consX, prop.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	vote1 := td.addPrepareVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	vote2 := td.addPrepareVote(td.consP, td.RandHash(), 1, 2, tIndexX)
	vote3 := td.addPrepareVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	vote4 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	vote5 := td.addPrepareVote(td.consP, td.RandHash(), 2, 0, tIndexX)
	vote6, _ := td.GenerateTestPrepareVote(1, 0)
	td.consP.AddVote(vote6)

	assert.False(t, td.consP.HasVote(vote1.Hash())) // previous round
	assert.True(t, td.consP.HasVote(vote2.Hash()))  // next round
	assert.True(t, td.consP.HasVote(vote3.Hash()))
	assert.True(t, td.consP.HasVote(vote4.Hash()))
	assert.False(t, td.consP.HasVote(vote5.Hash())) // valid votes for the next height
	assert.False(t, td.consP.HasVote(vote6.Hash())) // invalid votes

	assert.Equal(t, []*vote.Vote{vote3, vote4}, td.consP.AllVotes())
	assert.NotContains(t, td.consP.AllVotes(), vote2)
}

// TestConsensusLateProposal tests the scenario where a slow node receives a proposal
// after committing the block.
func TestConsensusLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)
	require.NotNil(t, prop)

	td.commitBlockForAllStates(t) // height 2

	// Partitioned node receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, prop.Block().Hash())
}

// TestSetProposalOnPrepare tests the scenario where a slow node receives a proposal
// in prepare phase.
func TestSetProposalOnPrepare(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)
	require.NotNil(t, prop)

	// The partitioned node receives all the votes first
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
}

// TestSetProposalOnPrecommit tests the scenario where a slow node receives a proposal
// in precommit phase.
func TestSetProposalOnPrecommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)
	require.NotNil(t, prop)

	// The partitioned node receives all the votes first
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consP, prop.Block().Hash())
}

// update me from TestHandleQueryVote: consensus:v1.
func TestHandleQueryVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	cpRound := int16(0)
	height := uint32(1)
	assert.Nil(t, td.consP.HandleQueryVote(height, 0))

	// Add some votes for Round 0
	td.addCPDecidedVote(td.consP, hash.UndefHash, height, 0, vote.CPValueYes,
		&vote.JustDecided{QCert: td.makeMainVoteCertificate(t, height, 0, cpRound)}, tIndexY)

	// Add some votes for Round 1
	td.enterNextRound(td.consP)
	td.addPrepareVote(td.consP, td.RandHash(), height, 1, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, height, 1, vote.CPValueYes,
		&vote.JustInitYes{}, tIndexY)
	td.addCPDecidedVote(td.consP, hash.UndefHash, height, 1, vote.CPValueYes,
		&vote.JustDecided{QCert: td.makeMainVoteCertificate(t, height, 1, cpRound)}, tIndexY)

	// Add some votes for Round 2
	td.enterNextRound(td.consP)
	td.addPrepareVote(td.consP, td.RandHash(), height, 2, tIndexY)

	t.Run("Query vote for round 0: should send the decided vote for the round 0", func(t *testing.T) {
		rndVote := td.consP.HandleQueryVote(height, 0)
		assert.Equal(t, vote.VoteTypeCPDecided, rndVote.Type())
		assert.Equal(t, height, rndVote.Height())
		assert.Equal(t, int16(0), rndVote.Round())
	})

	t.Run("Query vote for round 1: should send the decided vote for the round 1", func(t *testing.T) {
		rndVote := td.consP.HandleQueryVote(height, 1)
		assert.Equal(t, vote.VoteTypeCPDecided, rndVote.Type())
		assert.Equal(t, height, rndVote.Height())
		assert.Equal(t, int16(1), rndVote.Round())
	})

	t.Run("Query vote for round 2: should send the prepare vote for the round 2", func(t *testing.T) {
		rndVote := td.consP.HandleQueryVote(height, 2)
		assert.Equal(t, vote.VoteTypePrepare, rndVote.Type())
		assert.Equal(t, height, rndVote.Height())
		assert.Equal(t, int16(2), rndVote.Round())
	})

	t.Run("Query vote for round 3: should not send a vote for the next round", func(t *testing.T) {
		rndVote := td.consP.HandleQueryVote(height, 3)
		assert.Nil(t, rndVote)
	})

	t.Run("Query vote for height 2: should not send a vote for the next height", func(t *testing.T) {
		rndVote := td.consP.HandleQueryVote(height+1, 0)
		assert.Nil(t, rndVote)
	})
}

func TestHandleQueryProposal(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	// Round 1
	td.enterNextRound(td.consX)
	td.enterNextRound(td.consY) // consY is the proposer

	prop0 := td.consY.HandleQueryProposal(1, 0)
	assert.Nil(t, prop0, "proposer should not send a proposal for the previous round")

	prop1 := td.consX.HandleQueryProposal(1, 1)
	assert.Nil(t, prop1, "non-proposer should not send a proposal")

	prop2 := td.consY.HandleQueryProposal(1, 1)
	assert.NotNil(t, prop2, "proposer should send a proposal")

	td.consX.cpDecided = 0
	td.consX.SetProposal(td.consY.Proposal())
	prop3 := td.consX.HandleQueryProposal(1, 1)
	assert.NotNil(t, prop3, "non-proposer should send a proposal on decided proposal")

	prop4 := td.consX.HandleQueryProposal(2, 0)
	assert.Nil(t, prop4, "should not have a proposal for the next height")
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	td := setup(t)

	p := td.makeProposal(t, 1, 0)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// It should ignore proposal for previous rounds
	td.consP.SetProposal(p)

	assert.Nil(t, td.consP.Proposal())
	td.checkHeightRound(t, td.consP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	td := setup(t)

	p := td.makeProposal(t, 1, 0)
	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	td.consP.SetProposal(p)
	assert.Nil(t, td.consP.Proposal())
	td.checkHeightRound(t, td.consP, 2, 0)
}

func TestDuplicateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	height := uint32(4)
	round := int16(0)
	prop1 := td.makeProposal(t, height, round)
	trx := tx.NewTransferTx(height, td.consX.rewardAddr, td.RandAccAddress(), 1000,
		1000, tx.WithMemo("proposal changer"))
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)

	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, height, round)
	assert.NotEqual(t, prop1.Hash(), p2.Hash())

	td.consX.SetProposal(prop1)
	td.consX.SetProposal(p2)

	assert.Equal(t, prop1.Hash(), td.consX.Proposal().Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	pipe := pipeline.MockingPipeline[message.Message]()
	consInt := NewConsensus(testConfig(), state.MockingState(td.TestSuite),
		valKey, valKey.Address(), pipe, newConcreteMediator())
	nonActiveCons := consInt.(*consensus)

	t.Run("non-active instances should be in new-height state", func(t *testing.T) {
		nonActiveCons.MoveToNewHeight()
		td.newHeightTimeout(nonActiveCons)
		td.checkHeightRound(t, nonActiveCons, 1, 0)

		// Double entry
		nonActiveCons.MoveToNewHeight()
		td.newHeightTimeout(nonActiveCons)
		td.checkHeightRound(t, nonActiveCons, 1, 0)

		assert.False(t, nonActiveCons.IsActive())
		assert.Equal(t, "new-height", nonActiveCons.currentState.name())
	})

	t.Run("non-active instances should ignore proposals", func(t *testing.T) {
		p := td.makeProposal(t, 1, 0)
		nonActiveCons.SetProposal(p)

		assert.Nil(t, nonActiveCons.Proposal())
	})

	t.Run("non-active instances should ignore votes", func(t *testing.T) {
		v := td.addPrepareVote(nonActiveCons, td.RandHash(), 1, 0, tIndexX)

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

	v := td.addPrepareVote(td.consX, td.RandHash(), 1, util.MaxInt16, tIndexB)
	assert.True(t, td.consX.HasVote(v.Hash()))
}

func TestProposalWithBigRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, 1, util.MaxInt16)
	td.consP.SetProposal(p)
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

func TestCases(t *testing.T) {
	tests := []struct {
		seed        int64
		round       int16
		description string
	}{
		{1697898884837384019, 2, "1/3+ cp:PRE-VOTE in Prepare step"},
		{1734526933123806220, 1, "1/3+ cp:PRE-VOTE in Precommit step"},
		{1734526832618973590, 1, "Conflicting cp:PRE-VOTE in cp_round=0"},
		{1734527064850322674, 2, "Conflicting cp:PRE-VOTE in cp_round=1"},
		{1734526579569939721, 1, "consP & consB: Change Proposer, consX & consY: Commit (2 block announces)"},
	}

	for no, tt := range tests {
		td := setupWithSeed(t, tt.seed)
		td.commitBlockForAllStates(t)

		td.enterNewHeight(td.consX)
		td.enterNewHeight(td.consY)
		td.enterNewHeight(td.consB)
		td.enterNewHeight(td.consP)

		cert, err := checkConsensus(td, 2, nil)
		require.NoError(t, err,
			"test %v failed: %s", no+1, err)
		require.Equal(t, tt.round, cert.Round(),
			"test %v failed. round not matched (expected %d, got %d)",
			no+1, tt.round, cert.Round())
	}
}

func TestFaulty(t *testing.T) {
	for i := 0; i < 10; i++ {
		td := setup(t)
		td.commitBlockForAllStates(t)

		td.enterNewHeight(td.consX)
		td.enterNewHeight(td.consY)
		td.enterNewHeight(td.consB)
		td.enterNewHeight(td.consP)

		_, err := checkConsensus(td, 2, nil)
		require.NoError(t, err)
	}
}

// We have four nodes: X, Y, B, and P, which:
// - B is a Byzantine node
// - X, Y, and P are honest nodes
// - However, P is partitioned and perceives the network through B.
//
// At height H, B acts maliciously by double proposing:
// sending one proposal to X and Y, and another proposal to P.
//
// Once the partition is healed, honest nodes should either reach consensus
// on the first proposal or change the proposer.
// This is due to the randomness of the binary agreement.
func TestByzantine(t *testing.T) {
	td := setup(t)

	for i := 0; i < 6; i++ {
		td.commitBlockForAllStates(t)
	}

	height := uint32(7)
	round := int16(0)
	prop1 := td.makeProposal(t, height, round)

	// =================================
	// X, Y votes
	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	td.consX.SetProposal(prop1)
	td.consY.SetProposal(prop1)

	voteX := td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, prop1.Block().Hash())
	voteY := td.shouldPublishVote(t, td.consY, vote.VoteTypePrepare, prop1.Block().Hash())

	// Byzantine node doesn't broadcast the prepare vote
	// X and Y request to change proposer

	td.changeProposerTimeout(td.consX)
	td.changeProposerTimeout(td.consY)

	td.shouldPublishVote(t, td.consX, vote.VoteTypeCPPreVote, hash.UndefHash)
	td.shouldPublishVote(t, td.consY, vote.VoteTypeCPPreVote, hash.UndefHash)

	// X and Y are unable to progress

	// =================================
	// B votes
	td.enterNewHeight(td.consB)

	td.consB.SetProposal(prop1)

	td.consB.AddVote(voteX)
	td.consB.AddVote(voteY)
	td.shouldPublishVote(t, td.consB, vote.VoteTypePrepare, prop1.Block().Hash())
	td.shouldPublishVote(t, td.consB, vote.VoteTypePrecommit, prop1.Block().Hash())

	td.changeProposerTimeout(td.consB)

	// B requests to NOT change the proposer
	byzVote1 := td.shouldPublishVote(t, td.consB, vote.VoteTypeCPPreVote, prop1.Block().Hash())

	// =================================
	// P votes
	// Byzantine node create the second proposal and send it to the partitioned node P
	byzTrx := tx.NewTransferTx(height,
		td.consB.rewardAddr, td.RandAccAddress(), 1000, 1000)
	td.HelperSignTransaction(td.consB.valKey.PrivateKey(), byzTrx)
	assert.NoError(t, td.txPool.AppendTx(byzTrx))
	prop2 := td.makeProposal(t, height, round)

	require.NotEqual(t, prop1.Block().Hash(), prop2.Block().Hash())
	require.Equal(t, td.consB.valKey.Address(), prop1.Block().Header().ProposerAddress())
	require.Equal(t, td.consB.valKey.Address(), prop2.Block().Header().ProposerAddress())

	td.enterNewHeight(td.consP)

	// P receives the Seconds proposal
	td.consP.SetProposal(prop2)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, prop2.Block().Hash())
	byzVote2 := td.addPrepareVote(td.consP, prop2.Block().Hash(), height, round, tIndexB)

	// Request to change proposer
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)

	// P is unable to progress

	// =================================

	td.checkHeightRound(t, td.consX, height, round)
	td.checkHeightRound(t, td.consY, height, round)
	td.checkHeightRound(t, td.consP, height, round)

	// Let's make Byzantine node happy by removing his votes from the log
	for j := len(td.consMessages) - 1; j >= 0; j-- {
		if td.consMessages[j].sender == td.consB.valKey.Address() {
			td.consMessages = slices.Delete(td.consMessages, j, j+1)
		}
	}

	// =================================
	// Now, Partition heals
	fmt.Println("== Partition heals")
	cert, err := checkConsensus(td, height, []*vote.Vote{byzVote1, byzVote2})

	require.NoError(t, err)
	require.Equal(t, height, cert.Height())
	require.Contains(t, cert.Absentees(), int32(tIndexB))
}

func checkConsensus(td *testData, height uint32, byzVotes []*vote.Vote) (
	*certificate.BlockCertificate, error,
) {
	instances := []*consensus{td.consX, td.consY, td.consB, td.consP}

	if len(byzVotes) > 0 {
		for _, v := range byzVotes {
			td.consB.broadcastVote(v)
		}

		// remove byzantine node (Byzantine node goes offline)
		instances = []*consensus{td.consX, td.consY, td.consP}
	}

	// 70% chance for the first block to be lost
	changeProposerChance := 70

	blockAnnounces := map[crypto.Address]*message.BlockAnnounceMessage{}
	for len(td.consMessages) > 0 {
		rndIndex := td.RandInt(len(td.consMessages))
		rndMsg := td.consMessages[rndIndex]
		td.consMessages = slices.Delete(td.consMessages, rndIndex, rndIndex+1)

		switch rndMsg.message.Type() {
		case message.TypeVote:
			m := rndMsg.message.(*message.VoteMessage)
			if m.Vote.Height() == height {
				for _, cons := range instances {
					cons.AddVote(m.Vote)
				}
			}

		case message.TypeProposal:
			m := rndMsg.message.(*message.ProposalMessage)
			if m.Proposal.Height() == height {
				for _, cons := range instances {
					cons.SetProposal(m.Proposal)
				}
			}

		case message.TypeQueryProposal:
			m := rndMsg.message.(*message.QueryProposalMessage)
			if m.Height == height {
				for _, cons := range instances {
					p := cons.Proposal()
					if p != nil {
						td.consMessages = append(td.consMessages, consMessage{
							sender:  cons.valKey.Address(),
							message: message.NewProposalMessage(p),
						})
					}
				}
			}

		case message.TypeBlockAnnounce:
			m := rndMsg.message.(*message.BlockAnnounceMessage)
			blockAnnounces[rndMsg.sender] = m

		case
			message.TypeHello,
			message.TypeHelloAck,
			message.TypeTransaction,
			message.TypeQueryVote,
			message.TypeBlocksRequest,
			message.TypeBlocksResponse:
			//
		}

		for _, cons := range instances {
			rnd := td.RandInt(100)
			if rnd < changeProposerChance ||
				len(td.consMessages) == 0 {
				td.changeProposerTimeout(cons)
			}
		}
		changeProposerChance -= 5
	}

	// Check if more than 1/3 of nodes has committed the same block
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
