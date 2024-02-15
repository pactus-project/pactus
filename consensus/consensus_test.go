package consensus

import (
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
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/tx"
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
	}
}

func setup(t *testing.T) *testData {
	t.Helper()
	queryVoteInitialTimeout = 2 * time.Hour

	return setupWithSeed(t, testsuite.GenerateSeed())
}

func setupWithSeed(t *testing.T, seed int64) *testData {
	t.Helper()

	fmt.Printf("=== test %s, seed: %d\n", t.Name(), seed)

	ts := testsuite.NewTestSuiteForSeed(seed)

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
	params := param.DefaultParams()
	params.CommitteeSize = 4

	// to prevent triggering timers before starting the tests to avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockIntervalInSecond).Add(time.Duration(params.BlockIntervalInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)
	stX, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexX]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stY, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexY]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stB, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexB]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stP, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKeys[tIndexP]},
		store.MockingStore(ts), txPool, nil)
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
	td.consX = newConsensus(testConfig(), stX, valKeys[tIndexX],
		valKeys[tIndexX].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consY = newConsensus(testConfig(), stY, valKeys[tIndexY],
		valKeys[tIndexY].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consB = newConsensus(testConfig(), stB, valKeys[tIndexB],
		valKeys[tIndexB].PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	td.consP = newConsensus(testConfig(), stP, valKeys[tIndexP],
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

func (td *testData) shouldPublishBlockAnnounce(t *testing.T, cons *consensus, h hash.Hash) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeBlockAnnounce {
			m := consMsg.message.(*message.BlockAnnounceMessage)
			assert.Equal(t, m.Block.Hash(), h)

			return
		}
	}
	require.NoError(t, fmt.Errorf("Not found"))
}

func (td *testData) shouldPublishProposal(t *testing.T, cons *consensus,
	height uint32, round int16,
) *proposal.Proposal {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeProposal {
			m := consMsg.message.(*message.ProposalMessage)
			require.Equal(t, m.Proposal.Height(), height)
			require.Equal(t, m.Proposal.Round(), round)

			return m.Proposal
		}
	}
	require.NoError(t, fmt.Errorf("Not found"))

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
	require.NoError(t, fmt.Errorf("Not found"))
}

func (td *testData) shouldPublishQueryVote(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender != cons.valKey.Address() ||
			consMsg.message.Type() != message.TypeQueryVotes {
			continue
		}

		m := consMsg.message.(*message.QueryVotesMessage)
		assert.Equal(t, m.Height, height)
		assert.Equal(t, m.Round, round)
		assert.Equal(t, m.Querier, cons.valKey.Address())

		return
	}
	require.NoError(t, fmt.Errorf("Not found"))
}

func (td *testData) shouldPublishVote(t *testing.T, cons *consensus, voteType vote.Type, h hash.Hash) *vote.Vote {
	t.Helper()

	for i := len(td.consMessages) - 1; i >= 0; i-- {
		consMsg := td.consMessages[i]
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeVote {
			m := consMsg.message.(*message.VoteMessage)
			if m.Vote.Type() == voteType &&
				m.Vote.BlockHash() == h {
				return m.Vote
			}
		}
	}
	require.NoError(t, fmt.Errorf("Not found"))

	return nil
}

func checkHeightRound(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	h, r := cons.HeightRound()
	assert.Equal(t, h, height)
	assert.Equal(t, r, round)
}

func (td *testData) checkHeightRound(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	checkHeightRound(t, cons, height, round)
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
	cpRound int16, cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPPreVote(blockHash, height, round, cpRound, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addCPMainVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPMainVote(blockHash, height, round, cpRound, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addCPDecidedVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpRound int16, cpVal vote.CPValue, just vote.Just, valID int,
) {
	v := vote.NewCPDecidedVote(blockHash, height, round, cpRound, cpVal, just, td.valKeys[valID].Address())
	td.addVote(cons, v, valID)
}

func (td *testData) addVote(cons *consensus, v *vote.Vote, valID int) *vote.Vote {
	td.HelperSignVote(td.valKeys[valID], v)
	cons.AddVote(v)

	return v
}

func newHeightTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

func (td *testData) newHeightTimeout(cons *consensus) {
	newHeightTimeout(cons)
}

func (td *testData) queryProposalTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetQueryProposal})
	cons.lk.Unlock()
}

func (td *testData) changeProposerTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetChangeProposer})
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
func (td *testData) enterNextRound(cons *consensus) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func (td *testData) commitBlockForAllStates(t *testing.T) (*block.Block, *certificate.Certificate) {
	t.Helper()

	height := td.consX.bcState.LastBlockHeight()
	var err error
	p := td.makeProposal(t, height+1, 0)

	sb := certificate.BlockCertificateSignBytes(p.Block().Hash(), height+1, 0)
	sig1 := td.consX.valKey.Sign(sb)
	sig2 := td.consY.valKey.Sign(sb)
	sig3 := td.consB.valKey.Sign(sb)
	sig4 := td.consP.valKey.Sign(sb)

	sig := bls.SignatureAggregate(sig1, sig2, sig3, sig4)
	cert := certificate.NewCertificate(height+1, 0,
		[]int32{tIndexX, tIndexY, tIndexB, tIndexP}, []int32{}, sig)
	blk := p.Block()

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

	var p *proposal.Proposal
	switch (height % 4) + uint32(round%4) {
	case 1:
		blk, err := td.consX.bcState.ProposeBlock(td.consX.valKey, td.consX.rewardAddr)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consX.valKey, p)
	case 2:
		blk, err := td.consY.bcState.ProposeBlock(td.consY.valKey, td.consY.rewardAddr)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consY.valKey, p)
	case 3:
		blk, err := td.consB.bcState.ProposeBlock(td.consB.valKey, td.consB.rewardAddr)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consB.valKey, p)
	case 0, 4:
		blk, err := td.consP.bcState.ProposeBlock(td.consP.valKey, td.consP.rewardAddr)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.HelperSignProposal(td.consP.valKey, p)
	}

	return p
}

func TestStart(t *testing.T) {
	td := setup(t)

	td.consX.Start()
	td.shouldPublishQueryProposal(t, td.consX, 1)
	td.shouldPublishQueryVote(t, td.consX, 1, 0)
}

func TestNotInCommittee(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	str := store.MockingStore(td.TestSuite)

	st, _ := state.LoadOrNewState(td.genDoc, []*bls.ValidatorKey{valKey}, str, td.txPool, nil)
	Cons := NewConsensus(testConfig(), st, valKey, valKey.Address(), make(chan message.Message, 100),
		newConcreteMediator())
	cons := Cons.(*consensus)

	td.enterNewHeight(cons)
	td.newHeightTimeout(cons)
	assert.Equal(t, cons.currentState.name(), "new-height")
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

	p := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(p)

	td.addPrepareVote(td.consX, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consX, p.Block().Hash(), 2, 0, tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p.Block().Hash())

	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, p.Block().Hash())

	td.shouldPublishBlockAnnounce(t, td.consX, p.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	v1 := td.addPrepareVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	v2 := td.addPrepareVote(td.consP, td.RandHash(), 1, 2, tIndexX)
	v3 := td.addPrepareVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	v4 := td.addPrecommitVote(td.consP, td.RandHash(), 1, 1, tIndexX)
	v5 := td.addPrepareVote(td.consP, td.RandHash(), 2, 0, tIndexX)
	v6, _ := td.GenerateTestPrepareVote(1, 0)
	td.consP.AddVote(v6)

	assert.True(t, td.consP.HasVote(v1.Hash())) // previous round
	assert.True(t, td.consP.HasVote(v2.Hash())) // next round
	assert.True(t, td.consP.HasVote(v3.Hash()))
	assert.True(t, td.consP.HasVote(v4.Hash()))
	assert.False(t, td.consP.HasVote(v5.Hash())) // valid votes for the next height
	assert.False(t, td.consP.HasVote(v6.Hash())) // invalid votes

	assert.Equal(t, td.consP.AllVotes(), []*vote.Vote{v1, v3, v4})
	assert.NotContains(t, td.consP.AllVotes(), v2)
}

// TestConsensusLateProposal tests the scenario where a slow node receives a proposal
// in precommit phase.
func TestConsensusLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	p := td.makeProposal(t, h, r)
	require.NotNil(t, p)

	td.commitBlockForAllStates(t) // height 2

	// The partitioned node receives all the votes first
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(p)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
}

// TestConsensusVeryLateProposal tests the scenario where a slow node receives a proposal
// in prepare phase.
func TestConsensusVeryLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	p := td.makeProposal(t, h, r)
	require.NotNil(t, p)

	td.commitBlockForAllStates(t) // height 2

	// The partitioned node receives all the votes first
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(p)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consP, p.Block().Hash())
}

func TestPickRandomVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	assert.Nil(t, td.consP.PickRandomVote(0))
	cpRound := int16(1)

	// === make valid certificate
	sbPreVote := certificate.BlockCertificateSignBytes(hash.UndefHash, 1, 0)
	sbPreVote = append(sbPreVote, util.StringToBytes(vote.VoteTypeCPPreVote.String())...)
	sbPreVote = append(sbPreVote, util.Int16ToSlice(cpRound)...)
	sbPreVote = append(sbPreVote, byte(vote.CPValueOne))

	sbMainVote := certificate.BlockCertificateSignBytes(hash.UndefHash, 1, 0)
	sbMainVote = append(sbMainVote, util.StringToBytes(vote.VoteTypeCPMainVote.String())...)
	sbMainVote = append(sbMainVote, util.Int16ToSlice(cpRound)...)
	sbMainVote = append(sbMainVote, byte(vote.CPValueOne))

	committers := []int32{}
	preVoteSigs := []*bls.Signature{}
	mainVoteSigs := []*bls.Signature{}
	for i, val := range td.consP.validators {
		committers = append(committers, val.Number())
		preVoteSigs = append(preVoteSigs, td.valKeys[i].Sign(sbPreVote))
		mainVoteSigs = append(mainVoteSigs, td.valKeys[i].Sign(sbMainVote))
	}

	preVoteAggSig := bls.SignatureAggregate(preVoteSigs...)
	mainVoteAggSig := bls.SignatureAggregate(mainVoteSigs...)

	certPreVote := certificate.NewCertificate(1, 0, committers, []int32{}, preVoteAggSig)
	certMainVote := certificate.NewCertificate(1, 0, committers, []int32{}, mainVoteAggSig)
	// ====

	// round 0
	td.addPrepareVote(td.consP, td.RandHash(), 1, 0, tIndexX)
	td.addPrepareVote(td.consP, td.RandHash(), 1, 0, tIndexY)
	td.addCPPreVote(td.consP, hash.UndefHash, 1, 0, cpRound+1, vote.CPValueOne,
		&vote.JustPreVoteHard{QCert: certPreVote}, tIndexY)
	td.addCPMainVote(td.consP, hash.UndefHash, 1, 0, cpRound, vote.CPValueOne,
		&vote.JustMainVoteNoConflict{QCert: certPreVote}, tIndexY)
	td.addCPDecidedVote(td.consP, hash.UndefHash, 1, 0, cpRound, vote.CPValueOne,
		&vote.JustDecided{QCert: certMainVote}, tIndexY)

	assert.NotNil(t, td.consP.PickRandomVote(0))

	// Round 1
	td.enterNextRound(td.consP)
	td.addPrepareVote(td.consP, td.RandHash(), 1, 1, tIndexY)

	rndVote0 := td.consP.PickRandomVote(0)
	assert.NotEqual(t, rndVote0.Type(), vote.VoteTypePrepare, "Should not pick prepare votes")

	rndVote1 := td.consP.PickRandomVote(1)
	assert.Equal(t, rndVote1.Type(), vote.VoteTypePrepare)
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

	h := uint32(4)
	r := int16(0)
	p1 := td.makeProposal(t, h, r)
	trx := tx.NewTransferTx(h, td.consX.rewardAddr,
		td.RandAccAddress(), 1000, 1000, "proposal changer")
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)

	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.consX.SetProposal(p1)
	td.consX.SetProposal(p2)

	assert.Equal(t, td.consX.Proposal().Hash(), p1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	Cons := NewConsensus(testConfig(), state.MockingState(td.TestSuite),
		valKey, valKey.Address(), make(chan message.Message, 100), newConcreteMediator())
	nonActiveCons := Cons.(*consensus)

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
	assert.Equal(t, td.consP.log.RoundProposal(util.MaxInt16), p)
	assert.Nil(t, td.consP.Proposal())
}

func TestCases(t *testing.T) {
	tests := []struct {
		seed        int64
		round       int16
		description string
	}{
		{1697898884837384019, 2, "1/3+ cp:PRE-VOTE in prepare step"},
		{1694848907840926239, 0, "1/3+ cp:PRE-VOTE in precommit step"},
		{1694849103290580532, 1, "Conflicting votes, cp-round=0"},
		{1697900665869342730, 1, "Conflicting votes, cp-round=1"},
		{1697887970998950590, 1, "consP & consB: Change Proposer, consX & consY: Commit (2 block announces)"},
	}

	for i, test := range tests {
		td := setupWithSeed(t, test.seed)
		td.commitBlockForAllStates(t)

		td.enterNewHeight(td.consX)
		td.enterNewHeight(td.consY)
		td.enterNewHeight(td.consB)
		td.enterNewHeight(td.consP)

		cert, err := checkConsensus(td, 2, nil)
		require.NoError(t, err,
			"test %v failed: %s", i+1, err)
		require.Equal(t, cert.Round(), test.round,
			"test %v failed. round not matched (expected %d, got %d)",
			i+1, test.round, cert.Round())
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

	h := uint32(7)
	r := int16(0)
	p1 := td.makeProposal(t, h, r)

	// =================================
	// X, Y votes
	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	td.consX.SetProposal(p1)
	td.consY.SetProposal(p1)

	voteX := td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p1.Block().Hash())
	voteY := td.shouldPublishVote(t, td.consY, vote.VoteTypePrepare, p1.Block().Hash())

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

	td.consB.SetProposal(p1)

	td.consB.AddVote(voteX)
	td.consB.AddVote(voteY)
	td.shouldPublishVote(t, td.consB, vote.VoteTypePrepare, p1.Block().Hash())
	td.shouldPublishVote(t, td.consB, vote.VoteTypePrecommit, p1.Block().Hash())

	td.changeProposerTimeout(td.consB)

	// B requests to NOT change the proposer
	byzVote1 := td.shouldPublishVote(t, td.consB, vote.VoteTypeCPPreVote, p1.Block().Hash())

	// =================================
	// P votes
	// Byzantine node create the second proposal and send it to the partitioned node P
	byzTrx := tx.NewTransferTx(h,
		td.consB.rewardAddr, td.RandAccAddress(), 1000, 1000, "")
	td.HelperSignTransaction(td.consB.valKey.PrivateKey(), byzTrx)
	assert.NoError(t, td.txPool.AppendTx(byzTrx))
	p2 := td.makeProposal(t, h, r)

	require.NotEqual(t, p1.Block().Hash(), p2.Block().Hash())
	require.Equal(t, p1.Block().Header().ProposerAddress(), td.consB.valKey.Address())
	require.Equal(t, p2.Block().Header().ProposerAddress(), td.consB.valKey.Address())

	td.enterNewHeight(td.consP)

	// P receives the Seconds proposal
	td.consP.SetProposal(p2)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, p2.Block().Hash())
	byzVote2 := td.addPrepareVote(td.consP, p2.Block().Hash(), h, r, tIndexB)

	// Request to change proposer
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)

	// P is unable to progress

	// =================================

	td.checkHeightRound(t, td.consX, h, r)
	td.checkHeightRound(t, td.consY, h, r)
	td.checkHeightRound(t, td.consP, h, r)

	// Let's make Byzantine node happy by removing his votes from the log
	for j := len(td.consMessages) - 1; j >= 0; j-- {
		if td.consMessages[j].sender == td.consB.valKey.Address() {
			td.consMessages = slices.Delete(td.consMessages, j, j+1)
		}
	}

	// =================================
	// Now, Partition heals
	fmt.Println("== Partition heals")
	cert, err := checkConsensus(td, h, []*vote.Vote{byzVote1, byzVote2})

	require.NoError(t, err)
	require.Equal(t, cert.Height(), h)
	require.Contains(t, cert.Absentees(), int32(tIndexB))
}

func checkConsensus(td *testData, height uint32, byzVotes []*vote.Vote) (
	*certificate.Certificate, error,
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
			for _, cons := range instances {
				p := cons.Proposal()
				if p != nil {
					td.consMessages = append(td.consMessages, consMessage{
						sender:  cons.valKey.Address(),
						message: message.NewProposalMessage(p),
					})
				}
			}

		case message.TypeBlockAnnounce:
			m := rndMsg.message.(*message.BlockAnnounceMessage)
			blockAnnounces[rndMsg.sender] = m

		case
			message.TypeHello,
			message.TypeHelloAck,
			message.TypeTransactions,
			message.TypeQueryVotes,
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
