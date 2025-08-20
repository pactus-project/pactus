package consensusv2

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
	consX        *consensusV2 // Good peer
	consY        *consensusV2 // Good peer
	consB        *consensusV2 // Byzantine or offline peer
	consP        *consensusV2 // Partitioned peer
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
	overrideLogger := func(cons *consensusV2, name string) {
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

func (td *testData) shouldNotPublish(t *testing.T, cons *consensusV2, msgType message.Type) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == msgType {
			require.Error(t, fmt.Errorf("should not publish %s", msgType))
		}
	}
}

func (td *testData) shouldPublishBlockAnnounce(t *testing.T, cons *consensusV2, hash hash.Hash) {
	t.Helper()

	for _, consMsg := range td.consMessages {
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

	for _, consMsg := range td.consMessages {
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

	for _, consMsg := range td.consMessages {
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
	require.NoError(t, errors.New("Query proposal message not published"))
}

func (td *testData) shouldPublishVote(t *testing.T, cons *consensusV2, voteType vote.Type, hash hash.Hash) *vote.Vote {
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
	require.NoError(t, errors.New("Vote message not published"))

	return nil
}

func (*testData) checkHeightRound(t *testing.T, cons *consensusV2, height uint32, round int16) {
	t.Helper()

	h, r := cons.HeightRound()
	assert.Equal(t, h, height)
	assert.Equal(t, r, round)
}

func (td *testData) addPrecommitVote(cons *consensusV2, blockHash hash.Hash, height uint32, round int16,
	valID int,
) *vote.Vote {
	v := vote.NewPrecommitVote(blockHash, height, round, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPPreVote(cons *consensusV2, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	v := vote.NewCPPreVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPMainVote(cons *consensusV2, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	v := vote.NewCPMainVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPDecidedVote(cons *consensusV2, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	v := vote.NewCPDecidedVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addVote(cons *consensusV2, v *vote.Vote, valID int) *vote.Vote {
	td.HelperSignVote(td.valKeys[valID], v)
	cons.AddVote(v)

	return v
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
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetQueryVote})
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

func (td *testData) commitBlockForAllStates(t *testing.T) (*block.Block, *certificate.BlockCertificate) {
	t.Helper()

	height := td.consX.bcState.LastBlockHeight()
	var err error
	prop := td.makeProposal(t, height+1, 0)

	cert := certificate.NewBlockCertificate(height+1, 0)
	sb := cert.SignBytes(prop.Block().Hash())
	sigX := td.consX.valKey.Sign(sb)
	sigY := td.consY.valKey.Sign(sb)
	sigP := td.consP.valKey.Sign(sb)

	sig := bls.SignatureAggregate(sigX, sigY, sigP)
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
func (td *testData) makeProposal(t *testing.T, height uint32, round int16) *proposal.Proposal {
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

	blk, err := cons.bcState.ProposeBlock(cons.valKey, cons.rewardAddr)
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
) (vote.Just, vote.Just, vote.Just) {
	t.Helper()

	cpRound := int16(0)

	// Create PreVote Justification
	var preVoteJust vote.Just
	var cpValue vote.CPValue

	if propBlockHash != hash.UndefHash {
		cpValue = vote.CPValueNo
		committers := []int32{}
		sigs := []*bls.Signature{}
		for i, val := range td.consP.validators {
			vote := vote.NewPrecommitVote(propBlockHash, height, round, val.Address())
			signBytes := vote.SignBytes()

			committers = append(committers, val.Number())
			sigs = append(sigs, td.valKeys[i].Sign(signBytes))
		}
		aggSig := bls.SignatureAggregate(sigs...)
		cert := certificate.NewVoteCertificate(height, round)
		cert.SetSignature(committers, []int32{}, aggSig)

		preVoteJust = &vote.JustInitNo{
			QCert: cert,
		}
	} else {
		cpValue = vote.CPValueYes
		preVoteJust = &vote.JustInitYes{}
	}

	// Create MainVote Justification
	preVoteCommitters := []int32{}
	preVoteSigs := []*bls.Signature{}
	for i, val := range td.consP.validators {
		preVote := vote.NewCPPreVote(propBlockHash, height, round,
			cpRound, cpValue, preVoteJust, val.Address())
		signBytes := preVote.SignBytes()

		preVoteCommitters = append(preVoteCommitters, val.Number())
		preVoteSigs = append(preVoteSigs, td.valKeys[i].Sign(signBytes))
	}
	preVoteAggSig := bls.SignatureAggregate(preVoteSigs...)
	certPreVote := certificate.NewVoteCertificate(height, round)
	certPreVote.SetSignature(preVoteCommitters, []int32{}, preVoteAggSig)
	mainVoteJust := &vote.JustMainVoteNoConflict{QCert: certPreVote}

	// Create Decided Justification
	mainVoteCommitters := []int32{}
	mainVoteSigs := []*bls.Signature{}
	for i, val := range td.consP.validators {
		mainVote := vote.NewCPMainVote(propBlockHash, height, round,
			cpRound, cpValue, mainVoteJust, val.Address())
		signBytes := mainVote.SignBytes()

		mainVoteCommitters = append(mainVoteCommitters, val.Number())
		mainVoteSigs = append(mainVoteSigs, td.valKeys[i].Sign(signBytes))
	}
	mainVoteAggSig := bls.SignatureAggregate(mainVoteSigs...)
	certMainVote := certificate.NewVoteCertificate(height, round)
	certMainVote.SetSignature(mainVoteCommitters, []int32{}, mainVoteAggSig)
	decidedJust := &vote.JustDecided{QCert: certMainVote}

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
	pipe := pipeline.MockingPipeline[message.Message]()
	consInst := NewConsensus(testConfig(), state, valKey, valKey.Address(), pipe,
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

	vote1 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	vote2 := td.addPrecommitVote(td.consP, td.RandHash(), 2, 0, tIndexX)
	vote3 := td.addPrecommitVote(td.consP, td.RandHash(), 2, 0, tIndexY)
	vote4 := td.addPrecommitVote(td.consP, td.RandHash(), 3, 0, tIndexX)

	require.False(t, td.consP.HasVote(vote1.Hash()))
	require.True(t, td.consP.HasVote(vote2.Hash()))
	require.True(t, td.consP.HasVote(vote3.Hash()))
	require.False(t, td.consP.HasVote(vote4.Hash()))
}

func TestConsensusFastPath(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 2, 0)

	prop := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(prop)

	td.addPrecommitVote(td.consX, prop.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(td.consX, prop.Block().Hash(), 2, 0, tIndexB)
	td.addPrecommitVote(td.consX, prop.Block().Hash(), 2, 0, tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, prop.Block().Hash())

	td.shouldPublishBlockAnnounce(t, td.consX, prop.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	vote1 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	vote2 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 2, tIndexX)
	vote3 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	vote4 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	vote5 := td.addPrecommitVote(td.consP, td.RandHash(), 2, 0, tIndexX)
	vote6, _ := td.GenerateTestPrepareVote(1, 0)
	td.consP.AddVote(vote6)

	assert.False(t, td.consP.HasVote(vote1.Hash())) // previous round
	assert.True(t, td.consP.HasVote(vote2.Hash()))  // next round
	assert.True(t, td.consP.HasVote(vote3.Hash()))
	assert.True(t, td.consP.HasVote(vote4.Hash()))
	assert.False(t, td.consP.HasVote(vote5.Hash())) // valid votes for the next height
	assert.False(t, td.consP.HasVote(vote6.Hash())) // invalid votes

	assert.Equal(t, td.consP.AllVotes(), []*vote.Vote{vote3, vote4})
	assert.NotContains(t, td.consP.AllVotes(), vote2)
}

// TestConsensusLateProposal tests the scenario where a slow node doesn't have the proposal
// in prepare phase.
// func TestConsensusLateProposal(t *testing.T) {
// 	td := setup(t)

// 	td.commitBlockForAllStates(t) // height 1

// 	td.enterNewHeight(td.consP)

// 	height := uint32(2)
// 	round := int16(0)
// 	prop := td.makeProposal(t, height, round)
// 	blockHash := prop.Block().Hash()

// 	td.commitBlockForAllStates(t) // height 2

// 	// consP receives all the votes first
// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexX)
// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexY)
// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexB)

// 	td.shouldPublishQueryProposal(t, td.consP, height, round)

// 	// consP receives proposal now
// 	td.consP.SetProposal(prop)

// 	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, blockHash)
// 	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
// }

// TestConsensusVeryLateProposal tests the scenario where a slow node doesn't have the proposal
// in precommit phase.
// func TestConsensusVeryLateProposal(t *testing.T) {
// 	td := setup(t)

// 	td.commitBlockForAllStates(t) // height 1

// 	td.enterNewHeight(td.consP)

// 	height := uint32(2)
// 	round := int16(0)
// 	prop := td.makeProposal(t, height, round)
// 	blockHash := prop.Block().Hash()

// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexX)
// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexY)

// 	// consP timed out
// 	td.changeProposerTimeout(td.consP)

// 	_, _, decidedJust := td.makeChangeProposerJusts(t, prop.Block().Hash(), height, round)
// 	td.addCPDecidedVote(td.consP, prop.Block().Hash(), height, round, vote.CPValueNo, decidedJust, tIndexX)

// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexX)
// 	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexY)

// 	td.shouldPublishQueryProposal(t, td.consP, height, round)

// 	// consP receives proposal now
// 	td.consP.SetProposal(prop)

// 	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
// 	td.shouldPublishBlockAnnounce(t, td.consP, prop.Block().Hash())
// }

func TestHandleQueryVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	height := uint32(1)
	round := int16(0)

	assert.Nil(t, td.consP.HandleQueryVote(height, round))

	// round 0
	preVoteJust, _, decidedJust := td.makeChangeProposerJusts(t, hash.UndefHash, height, round)

	vote1 := td.addCPDecidedVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, decidedJust, tIndexY)

	// round 1
	td.enterNextRound(td.consP)
	round++

	preVoteJust, _, decidedJust = td.makeChangeProposerJusts(t, hash.UndefHash, height, round)
	vote2 := td.addPrecommitVote(td.consP, td.RandHash(), height, round, tIndexX)
	vote3 := td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, preVoteJust, tIndexY)
	vote4 := td.addCPDecidedVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, decidedJust, tIndexY)

	// Round 2
	td.enterNextRound(td.consP)
	round++

	vote5 := td.addPrecommitVote(td.consP, td.RandHash(), height, round, tIndexY)

	require.True(t, td.consP.HasVote(vote1.Hash()))
	require.True(t, td.consP.HasVote(vote2.Hash()))
	require.True(t, td.consP.HasVote(vote3.Hash()))
	require.True(t, td.consP.HasVote(vote4.Hash()))
	require.True(t, td.consP.HasVote(vote5.Hash()))

	rndVote0 := td.consP.HandleQueryVote(height, 0)
	assert.Equal(t, rndVote0, vote4, "should send the decided vote for the previous round")

	rndVote1 := td.consP.HandleQueryVote(height, 1)
	assert.Equal(t, rndVote1, vote4, "should send the decided vote for the previous round")

	rndVote2 := td.consP.HandleQueryVote(height, 2)
	assert.Equal(t, rndVote2, vote5, "should send the prepare vote for the current round")

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

	td.consX.cpDecidedCert = td.GenerateTestVoteCertificate(1) // TODO: better way?
	prop3 := td.consX.HandleQueryProposal(height, round)
	assert.NotNil(t, prop3, "non-proposer should send a proposal on decided proposal")

	prop4 := td.consX.HandleQueryProposal(height+1, 0)
	assert.Nil(t, prop4, "should not have a proposal for the next height")
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	td := setup(t)

	prop := td.makeProposal(t, 1, 0)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// It should ignore proposal for previous rounds
	td.consP.SetProposal(prop)

	assert.Nil(t, td.consP.Proposal())
	td.checkHeightRound(t, td.consP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	td := setup(t)

	prop := td.makeProposal(t, 1, 0)
	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	td.consP.SetProposal(prop)
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
	trx := tx.NewTransferTx(height, td.consX.rewardAddr, td.RandAccAddress(), 1000, 1000)
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)

	assert.NoError(t, td.txPool.AppendTx(trx))
	prop2 := td.makeProposal(t, height, round)
	assert.NotEqual(t, prop1.Hash(), prop2.Hash())

	td.consX.SetProposal(prop1)
	td.consX.SetProposal(prop2)

	assert.Equal(t, td.consX.Proposal().Hash(), prop1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	pipe := pipeline.MockingPipeline[message.Message]()
	consInst := NewConsensus(testConfig(), state.MockingState(td.TestSuite),
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
		v := td.addPrecommitVote(nonActiveCons, td.RandHash(), 1, 0, tIndexX)

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

	vote := td.addPrecommitVote(td.consX, td.RandHash(), 1, util.MaxInt16, tIndexB)
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

func TestCases(t *testing.T) {
	tests := []struct {
		seed        int64
		round       int16
		description string
	}{
		// {1697898884837384019, 2, "1/3+ cp:PRE-VOTE in prepare step"},
		// {1694848907840926239, 0, "1/3+ cp:PRE-VOTE in precommit step"},
		// {1694849103290580532, 1, "Conflicting votes, cp-round=0"},
		// {1697900665869342730, 1, "Conflicting votes, cp-round=1"},
		// {1697887970998950590, 1, "consP & consB: Change Proposer, consX & consY: Commit (2 block announces)"},
		{1717870730391411396, 2, "move to the next round on decided vote"},
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
		require.Equal(t, cert.Round(), tt.round,
			"test %v failed. round not matched (expected %d, got %d)",
			no+1, tt.round, cert.Round())
	}
}

// func TestFaulty(t *testing.T) {
// 	for i := 0; i < 100000; i++ {
// 		td := setup(t)
// 		td.commitBlockForAllStates(t)

// 		td.enterNewHeight(td.consX)
// 		td.enterNewHeight(td.consY)
// 		td.enterNewHeight(td.consB)
// 		td.enterNewHeight(td.consP)

// 		_, err := checkConsensus(td, 2, nil)
// 		require.NoError(t, err)
// 	}
// }

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
// func TestByzantine1(t *testing.T) {
// 	td := setup(t)

// 	for i := 0; i < 6; i++ {
// 		td.commitBlockForAllStates(t)
// 	}

// 	height := uint32(7)
// 	round := int16(0)
// 	prop1 := td.makeProposal(t, height, round)

// 	// =================================
// 	// X, Y votes
// 	td.enterNewHeight(td.consX)
// 	td.enterNewHeight(td.consY)

// 	td.consX.SetProposal(prop1)
// 	td.consY.SetProposal(prop1)

// 	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, prop1.Block().Hash())
// 	td.shouldPublishVote(t, td.consY, vote.VoteTypePrecommit, prop1.Block().Hash())

// 	// Byzantine node doesn't broadcast the prepare vote
// 	// X and Y request to change proposer

// 	td.changeProposerTimeout(td.consX)
// 	td.changeProposerTimeout(td.consY)

// 	voteX := td.shouldPublishVote(t, td.consX, vote.VoteTypeCPPreVote, hash.UndefHash)
// 	voteY := td.shouldPublishVote(t, td.consY, vote.VoteTypeCPPreVote, hash.UndefHash)

// 	// X and Y are unable to progress

// 	// =================================
// 	// B votes
// 	td.enterNewHeight(td.consB)

// 	td.consB.SetProposal(prop1)

// 	td.consB.AddVote(voteX)
// 	td.consB.AddVote(voteY)
// 	td.shouldPublishVote(t, td.consB, vote.VoteTypePrecommit, prop1.Block().Hash())

// 	td.changeProposerTimeout(td.consB)

// 	// B requests to NOT change the proposer
// 	// byzVote1 := td.shouldPublishVote(t, td.consB, vote.VoteTypeCPPreVote, p1.Block().Hash())

// 	// =================================
// 	// P votes
// 	// Byzantine node create the second proposal and send it to the partitioned node P
// 	byzTrx := tx.NewTransferTx(height,
// 		td.consB.rewardAddr, td.RandAccAddress(), 1000, 1000)
// 	td.HelperSignTransaction(td.consB.valKey.PrivateKey(), byzTrx)
// 	assert.NoError(t, td.txPool.AppendTx(byzTrx))
// 	prop2 := td.makeProposal(t, height, round)

// 	require.NotEqual(t, prop1.Block().Hash(), prop2.Block().Hash())
// 	require.Equal(t, prop1.Block().Header().ProposerAddress(), td.consB.valKey.Address())
// 	require.Equal(t, prop2.Block().Header().ProposerAddress(), td.consB.valKey.Address())

// 	td.enterNewHeight(td.consP)

// 	// P receives the Seconds proposal
// 	td.consP.SetProposal(prop2)

// 	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop2.Block().Hash())
// 	byzVote2 := td.addPrecommitVote(td.consP, prop2.Block().Hash(), height, round, tIndexB)

// 	// Request to change proposer
// 	td.changeProposerTimeout(td.consP)

// 	byzVote1 := td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)

// 	// P is unable to progress

// 	// =================================

// 	td.checkHeightRound(t, td.consX, height, round)
// 	td.checkHeightRound(t, td.consY, height, round)
// 	td.checkHeightRound(t, td.consP, height, round)

// 	// Let's make Byzantine node happy by removing his votes from the log
// 	for j := len(td.consMessages) - 1; j >= 0; j-- {
// 		if td.consMessages[j].sender == td.consB.valKey.Address() {
// 			td.consMessages = slices.Delete(td.consMessages, j, j+1)
// 		}
// 	}

// 	// =================================
// 	// Now, Partition heals
// 	fmt.Println("== Partition heals")
// 	cert, err := checkConsensus(td, height, []*vote.Vote{byzVote1, byzVote2})

// 	require.NoError(t, err)
// 	require.Equal(t, cert.Height(), height)
// 	require.Contains(t, cert.Absentees(), int32(tIndexB))
// }

// In this test, B is a Byzantine node and the network is partitioned.
// B acts maliciously by double proposing:
// sending one proposal to X and Y, and another proposal to P, M and N.
//
// Once the partition is healed, honest nodes should either reach consensus
// on one proposal or change the proposer.
// This is due to the randomness of the binary agreement.
// func TestByzantine(t *testing.T) {
// 	td := setup(t)

// 	for i := 0; i < 6; i++ {
// 		td.commitBlockForAllStates(t)
// 	}

// 	height := uint32(7)
// 	round := int16(0)

// 	// =================================
// 	// X, Y votes
// 	td.enterNewHeight(td.consX)
// 	td.enterNewHeight(td.consY)

// 	prop := td.makeProposal(t, height, round)
// 	require.Equal(t, prop.Block().Header().ProposerAddress(), td.consB.valKey.Address())

// 	// X and Y receive the first proposal
// 	td.consX.SetProposal(prop)
// 	td.consY.SetProposal(prop)

// 	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, prop.Block().Hash())
// 	td.shouldPublishVote(t, td.consY, vote.VoteTypePrecommit, prop.Block().Hash())

// 	// X and Y don't have enough votes, so they request to change the proposer
// 	td.changeProposerTimeout(td.consX)
// 	td.changeProposerTimeout(td.consY)

// 	// X and Y are unable to progress

// 	// =================================
// 	// P, M and N votes
// 	// Byzantine node create the second proposal and send it to the partitioned nodes
// 	byzTrx := tx.NewTransferTx(height,
// 		td.consB.rewardAddr, td.RandAccAddress(), 1000, 1000)
// 	td.HelperSignTransaction(td.consB.valKey.PrivateKey(), byzTrx)
// 	assert.NoError(t, td.txPool.AppendTx(byzTrx))
// 	byzProp := td.makeProposal(t, height, round)

// 	require.NotEqual(t, prop.Block().Hash(), byzProp.Block().Hash())
// 	require.Equal(t, byzProp.Block().Header().ProposerAddress(), td.consB.valKey.Address())

// 	td.enterNewHeight(td.consP)

// 	// P, M and N receive the Seconds proposal
// 	td.consP.SetProposal(byzProp)

// 	voteP := td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, byzProp.Block().Hash())

// 	// P, M and N don't have enough votes, so they request to change the proposer
// 	td.changeProposerTimeout(td.consP)

// 	// P, M and N are unable to progress

// 	// =================================
// 	// B votes
// 	// B requests to NOT change the proposer

// 	td.enterNewHeight(td.consB)

// 	voteB := vote.NewPrepareVote(byzProp.Block().Hash(), height, round, td.consB.valKey.Address())
// 	td.HelperSignVote(td.consB.valKey, voteB)
// 	byzJust0Block := &vote.JustInitNo{
// 		QCert: td.consB.makeVoteCertificate(
// 			map[crypto.Address]*vote.Vote{
// 				voteP.Signer(): voteP,
// 				voteB.Signer(): voteB,
// 			}),
// 	}
// 	byzVote := vote.NewCPPreVote(byzProp.Block().Hash(), height, round, 0, vote.CPValueNo, byzJust0Block, td.consB.valKey.Address())
// 	td.HelperSignVote(td.consB.valKey, byzVote)

// 	// =================================

// 	td.checkHeightRound(t, td.consX, height, round)
// 	td.checkHeightRound(t, td.consY, height, round)
// 	td.checkHeightRound(t, td.consP, height, round)

// 	// =================================
// 	// Now, Partition heals
// 	fmt.Println("== Partition heals")
// 	cert, err := checkConsensus(td, height, []*vote.Vote{byzVote})

// 	require.NoError(t, err)
// 	require.Equal(t, cert.Height(), height)
// 	require.Contains(t, cert.Absentees(), int32(tIndexB))
// }

func checkConsensus(td *testData, height uint32, byzVotes []*vote.Vote) (
	*certificate.BlockCertificate, error,
) {
	instances := []*consensusV2{td.consX, td.consY, td.consB, td.consP}

	if len(byzVotes) > 0 {
		for _, v := range byzVotes {
			td.consB.broadcastVote(v)
		}

		// remove byzantine node (Byzantine node goes offline)
		instances = []*consensusV2{td.consX, td.consY, td.consP}
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
					p := cons.HandleQueryProposal(m.Height, m.Round)
					if p != nil {
						td.consMessages = append(td.consMessages, consMessage{
							sender:  cons.valKey.Address(),
							message: message.NewProposalMessage(p),
						})
					}
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
			rnd := td.RandInt(100)
			if rnd < changeProposerChance ||
				len(td.consMessages) == 0 {
				td.changeProposerTimeout(cons)
			}
		}
		changeProposerChance -= 5
	}

	// Verify whether more than (2f+1) nodes have committed to the same block.
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
