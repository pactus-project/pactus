package fastconsensus

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
	tIndexM = 4
	tIndexN = 5
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
	consM        *consensus // Witness Peer
	consN        *consensus // Witness Peer
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

	ts := testsuite.NewTestSuiteForSeed(seed)

	_, valKeys := ts.GenerateTestCommittee(6)
	txPool := txpool.MockingTxPool()

	vals := make([]*validator.Validator, 6)
	for i, key := range valKeys {
		val := validator.NewValidator(key.PublicKey(), int32(i))
		vals[i] = val
	}

	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	params := param.DefaultParams()
	params.CommitteeSize = 6

	// To prevent triggering timers before starting the tests and
	// avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockIntervalInSecond).
		Add(time.Duration(params.BlockIntervalInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)

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

	instances := make([]*consensus, len(valKeys))
	for i, valKey := range valKeys {
		bcState, err := state.LoadOrNewState(genDoc, []*bls.ValidatorKey{valKey},
			store.MockingStore(ts), txPool, nil)
		require.NoError(t, err)

		instances[i] = makeConsensus(testConfig(), bcState, valKey,
			valKey.PublicKey().AccountAddress(), broadcasterFunc, newConcreteMediator())
	}

	td.consX = instances[tIndexX]
	td.consY = instances[tIndexY]
	td.consB = instances[tIndexB]
	td.consP = instances[tIndexP]
	td.consM = instances[tIndexM]
	td.consN = instances[tIndexN]

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
	overrideLogger(td.consM, "consM")
	overrideLogger(td.consN, "consN")
	// -------------------------------

	logger.Info("setup finished, start running the test", "name", t.Name())

	return td
}

func (td *testData) shouldNotPublish(t *testing.T, cons *consensus, msgType message.Type) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == msgType {
			require.Error(t, fmt.Errorf("should not publish %s", msgType))
		}
	}
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
	require.NoError(t, fmt.Errorf("Block announce message not published"))
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
	require.NoError(t, fmt.Errorf("Proposal message not published"))

	return nil
}

func (td *testData) shouldPublishQueryProposal(t *testing.T, cons *consensus, height uint32, round int16) {
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
	require.NoError(t, fmt.Errorf("Query proposal message not published"))
}

func (td *testData) shouldPublishQueryVote(t *testing.T, cons *consensus, height uint32, round int16) {
	t.Helper()

	for _, consMsg := range td.consMessages {
		if consMsg.sender != cons.valKey.Address() ||
			consMsg.message.Type() != message.TypeQueryVote {
			continue
		}

		m := consMsg.message.(*message.QueryVotesMessage)
		assert.Equal(t, m.Height, height)
		assert.Equal(t, m.Round, round)
		assert.Equal(t, m.Querier, cons.valKey.Address())

		return
	}
	require.NoError(t, fmt.Errorf("Query proposal message not published"))
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
	require.NoError(t, fmt.Errorf("Vote message not published"))

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
) *vote.Vote {
	v := vote.NewCPPreVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPMainVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	v := vote.NewCPMainVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addCPDecidedVote(cons *consensus, blockHash hash.Hash, height uint32, round int16,
	cpVal vote.CPValue, just vote.Just, valID int,
) *vote.Vote {
	v := vote.NewCPDecidedVote(blockHash, height, round, 0, cpVal, just, td.valKeys[valID].Address())

	return td.addVote(cons, v, valID)
}

func (td *testData) addVote(cons *consensus, v *vote.Vote, valID int) *vote.Vote {
	td.HelperSignVote(td.valKeys[valID], v)
	cons.AddVote(v)

	return v
}

func (*testData) newHeightTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

func (*testData) queryProposalTimeout(cons *consensus) {
	cons.lk.Lock()
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetQueryProposal})
	cons.lk.Unlock()
}

func (*testData) changeProposerTimeout(cons *consensus) {
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

	cert := certificate.NewBlockCertificate(height+1, 0, true)
	sb := cert.SignBytes(prop.Block().Hash())
	sig0 := td.consX.valKey.Sign(sb)
	sig1 := td.consY.valKey.Sign(sb)
	sig2 := td.consB.valKey.Sign(sb)
	sig3 := td.consP.valKey.Sign(sb)
	sig4 := td.consM.valKey.Sign(sb)

	sig := bls.SignatureAggregate(sig0, sig1, sig2, sig3, sig4)
	cert.SetSignature([]int32{0, 1, 2, 3, 4, 5}, []int32{5}, sig)
	blk := prop.Block()

	err = td.consX.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consY.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consB.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consP.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consM.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)
	err = td.consN.bcState.CommitBlock(blk, cert)
	assert.NoError(t, err)

	return blk, cert
}

// makeProposal generates a signed and valid proposal for the given height and round.
func (td *testData) makeProposal(t *testing.T, height uint32, round int16) *proposal.Proposal {
	t.Helper()

	var cons *consensus
	switch (height % 6) + uint32(round%6) {
	case 1:
		cons = td.consX
	case 2:
		cons = td.consY
	case 3:
		cons = td.consB
	case 4:
		cons = td.consP
	case 5:
		cons = td.consM
	case 0, 6:
		cons = td.consN
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
		prepareCommitters := []int32{}
		prepareSigs := []*bls.Signature{}
		for i, val := range td.consP.validators {
			prepareVote := vote.NewPrepareVote(propBlockHash, height, round, val.Address())
			signBytes := prepareVote.SignBytes()

			prepareCommitters = append(prepareCommitters, val.Number())
			prepareSigs = append(prepareSigs, td.valKeys[i].Sign(signBytes))
		}
		prepareAggSig := bls.SignatureAggregate(prepareSigs...)
		certPrepare := certificate.NewVoteCertificate(height, round)
		certPrepare.SetSignature(prepareCommitters, []int32{}, prepareAggSig)

		preVoteJust = &vote.JustInitNo{
			QCert: certPrepare,
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

	td.consX.Start()
	td.checkHeightRound(t, td.consX, 1, 0)
}

func TestNotInCommittee(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	str := store.MockingStore(td.TestSuite)

	st, _ := state.LoadOrNewState(td.genDoc, []*bls.ValidatorKey{valKey}, str, td.txPool, nil)
	consInst := NewConsensus(testConfig(), st, valKey, valKey.Address(), make(chan message.Message, 100),
		newConcreteMediator())
	cons := consInst.(*consensus)

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

func TestConsensusFastPath(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 2, 0)

	prop := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(prop)

	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexB)
	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexP)
	td.addPrepareVote(td.consX, prop.Block().Hash(), 2, 0, tIndexM)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, prop.Block().Hash())

	td.shouldPublishBlockAnnounce(t, td.consX, prop.Block().Hash())
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

	assert.False(t, td.consP.HasVote(v1.Hash())) // previous round
	assert.True(t, td.consP.HasVote(v2.Hash()))  // next round
	assert.True(t, td.consP.HasVote(v3.Hash()))
	assert.True(t, td.consP.HasVote(v4.Hash()))
	assert.False(t, td.consP.HasVote(v5.Hash())) // valid votes for the next height
	assert.False(t, td.consP.HasVote(v6.Hash())) // invalid votes

	assert.Equal(t, td.consP.AllVotes(), []*vote.Vote{v3, v4})
	assert.NotContains(t, td.consP.AllVotes(), v2)
}

// TestConsensusLateProposal tests the scenario where a slow node doesn't have the proposal
// in prepare phase.
func TestConsensusLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	prop := td.makeProposal(t, h, r)
	blockHash := prop.Block().Hash()

	td.commitBlockForAllStates(t) // height 2

	// consP receives all the votes first
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexY)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexB)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexM)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexN)

	td.shouldPublishQueryProposal(t, td.consP, h, r)

	// consP receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, blockHash)
	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
}

// TestConsensusVeryLateProposal tests the scenario where a slow node doesn't have the proposal
// in precommit phase.
func TestConsensusVeryLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	prop := td.makeProposal(t, h, r)
	blockHash := prop.Block().Hash()

	td.addPrepareVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexY)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexM)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexN)

	// consP timed out
	td.changeProposerTimeout(td.consP)

	_, _, decidedJust := td.makeChangeProposerJusts(t, prop.Block().Hash(), h, r)
	td.addCPDecidedVote(td.consP, prop.Block().Hash(), h, r, vote.CPValueNo, decidedJust, tIndexX)

	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexY)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexM)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexN)

	td.shouldPublishQueryProposal(t, td.consP, h, r)

	// consP receives proposal now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consP, prop.Block().Hash())
}

func TestPickRandomVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	assert.Nil(t, td.consP.PickRandomVote(0))

	h := uint32(1)
	r := int16(0)

	preVoteJust, mainVoteJust, decidedJust := td.makeChangeProposerJusts(t, hash.UndefHash, h, r)

	// round 0
	v1 := td.addPrepareVote(td.consP, td.RandHash(), h, r, tIndexX)
	v2 := td.addPrepareVote(td.consP, td.RandHash(), h, r, tIndexY)
	v3 := td.addCPPreVote(td.consP, hash.UndefHash, h, r, vote.CPValueYes, preVoteJust, tIndexY)
	v4 := td.addCPMainVote(td.consP, hash.UndefHash, h, r, vote.CPValueYes, mainVoteJust, tIndexY)
	v5 := td.addCPDecidedVote(td.consP, hash.UndefHash, h, r, vote.CPValueYes, decidedJust, tIndexY)

	// Round 1
	td.enterNextRound(td.consP)
	v6 := td.addPrepareVote(td.consP, td.RandHash(), h, r+1, tIndexY)

	require.True(t, td.consP.HasVote(v1.Hash()))
	require.True(t, td.consP.HasVote(v2.Hash()))
	require.True(t, td.consP.HasVote(v3.Hash()))
	require.True(t, td.consP.HasVote(v4.Hash()))
	require.True(t, td.consP.HasVote(v5.Hash()))
	require.True(t, td.consP.HasVote(v6.Hash()))

	rndVote0 := td.consP.PickRandomVote(r)
	assert.Equal(t, rndVote0, v5, "for past round should pick Decided votes only")

	rndVote1 := td.consP.PickRandomVote(r + 1)
	assert.Equal(t, rndVote1, v6)

	rndVote2 := td.consP.PickRandomVote(r + 2)
	assert.Nil(t, rndVote2)
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

	h := uint32(4)
	r := int16(0)
	prop1 := td.makeProposal(t, h, r)
	trx := tx.NewTransferTx(h, td.consX.rewardAddr,
		td.RandAccAddress(), 1000, 1000, "proposal changer")
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)

	assert.NoError(t, td.txPool.AppendTx(trx))
	prop2 := td.makeProposal(t, h, r)
	assert.NotEqual(t, prop1.Hash(), prop2.Hash())

	td.consX.SetProposal(prop1)
	td.consX.SetProposal(prop2)

	assert.Equal(t, td.consX.Proposal().Hash(), prop1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	valKey := td.RandValKey()
	consInst := NewConsensus(testConfig(), state.MockingState(td.TestSuite),
		valKey, valKey.Address(), make(chan message.Message, 100), newConcreteMediator())
	nonActiveCons := consInst.(*consensus)

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

	prop := td.makeProposal(t, 1, util.MaxInt16)
	td.consP.SetProposal(prop)
	assert.Nil(t, td.consP.Proposal())
}

func TestInvalidProposal(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, 1, 0)
	p.SetSignature(nil) // Make proposal invalid
	td.consP.SetProposal(p)
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

	for i, test := range tests {
		td := setupWithSeed(t, test.seed)
		td.commitBlockForAllStates(t)

		td.enterNewHeight(td.consX)
		td.enterNewHeight(td.consY)
		td.enterNewHeight(td.consB)
		td.enterNewHeight(td.consP)
		td.enterNewHeight(td.consM)
		td.enterNewHeight(td.consN)

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
		td.enterNewHeight(td.consM)
		td.enterNewHeight(td.consN)

		_, err := checkConsensus(td, 2, nil)
		require.NoError(t, err)
	}
}

// In this test, B is a Byzantine node and the network is partitioned.
// B acts maliciously by double proposing:
// sending one proposal to X and Y, and another proposal to P, M and N.
//
// Once the partition is healed, honest nodes should either reach consensus
// on one proposal or change the proposer.
// This is due to the randomness of the binary agreement.
func TestByzantine(t *testing.T) {
	td := setup(t)

	for i := 0; i < 8; i++ {
		td.commitBlockForAllStates(t)
	}

	h := uint32(9)
	r := int16(0)

	// =================================
	// X, Y votes
	td.enterNewHeight(td.consX)
	td.enterNewHeight(td.consY)

	prop := td.makeProposal(t, h, r)
	require.Equal(t, prop.Block().Header().ProposerAddress(), td.consB.valKey.Address())

	// X and Y receive the Seconds proposal
	td.consX.SetProposal(prop)
	td.consY.SetProposal(prop)

	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, prop.Block().Hash())
	td.shouldPublishVote(t, td.consY, vote.VoteTypePrepare, prop.Block().Hash())

	// X and Y don't have enough votes, so they request to change the proposer
	td.changeProposerTimeout(td.consX)
	td.changeProposerTimeout(td.consY)

	// X and Y are unable to progress

	// =================================
	// P, M and N votes
	// Byzantine node create the second proposal and send it to the partitioned nodes
	byzTrx := tx.NewTransferTx(h,
		td.consB.rewardAddr, td.RandAccAddress(), 1000, 1000, "")
	td.HelperSignTransaction(td.consB.valKey.PrivateKey(), byzTrx)
	assert.NoError(t, td.txPool.AppendTx(byzTrx))
	byzProp := td.makeProposal(t, h, r)

	require.NotEqual(t, prop.Block().Hash(), byzProp.Block().Hash())
	require.Equal(t, byzProp.Block().Header().ProposerAddress(), td.consB.valKey.Address())

	td.enterNewHeight(td.consP)
	td.enterNewHeight(td.consM)
	td.enterNewHeight(td.consN)

	// P, M and N receive the Seconds proposal
	td.consP.SetProposal(byzProp)
	td.consM.SetProposal(byzProp)
	td.consN.SetProposal(byzProp)

	voteP := td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, byzProp.Block().Hash())
	voteM := td.shouldPublishVote(t, td.consM, vote.VoteTypePrepare, byzProp.Block().Hash())
	voteN := td.shouldPublishVote(t, td.consN, vote.VoteTypePrepare, byzProp.Block().Hash())

	// P, M and N don't have enough votes, so they request to change the proposer
	td.changeProposerTimeout(td.consP)
	td.changeProposerTimeout(td.consM)
	td.changeProposerTimeout(td.consN)

	// P, M and N are unable to progress

	// =================================
	// B votes
	// B requests to NOT change the proposer

	td.enterNewHeight(td.consB)

	voteB := vote.NewPrepareVote(byzProp.Block().Hash(), h, r, td.consB.valKey.Address())
	td.HelperSignVote(td.consB.valKey, voteB)
	byzJust0Block := &vote.JustInitNo{
		QCert: td.consB.makeVoteCertificate(
			map[crypto.Address]*vote.Vote{
				voteP.Signer(): voteP,
				voteM.Signer(): voteM,
				voteN.Signer(): voteN,
				voteB.Signer(): voteB,
			}),
	}
	byzVote := vote.NewCPPreVote(byzProp.Block().Hash(), h, r, 0, vote.CPValueNo, byzJust0Block, td.consB.valKey.Address())
	td.HelperSignVote(td.consB.valKey, byzVote)

	// =================================

	td.checkHeightRound(t, td.consX, h, r)
	td.checkHeightRound(t, td.consY, h, r)
	td.checkHeightRound(t, td.consP, h, r)
	td.checkHeightRound(t, td.consM, h, r)
	td.checkHeightRound(t, td.consN, h, r)

	// =================================
	// Now, Partition heals
	fmt.Println("== Partition heals")
	cert, err := checkConsensus(td, h, []*vote.Vote{byzVote})

	require.NoError(t, err)
	require.Equal(t, cert.Height(), h)
	require.Contains(t, cert.Absentees(), int32(tIndexB))
}

func checkConsensus(td *testData, height uint32, byzVotes []*vote.Vote) (
	*certificate.BlockCertificate, error,
) {
	instances := []*consensus{td.consX, td.consY, td.consB, td.consP, td.consM, td.consN}

	if len(byzVotes) > 0 {
		for _, v := range byzVotes {
			td.consB.broadcastVote(v)
		}

		// remove byzantine node (Byzantine node goes offline)
		instances = []*consensus{td.consX, td.consY, td.consP, td.consM, td.consN}
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

	// Verify whether more than (3t+1) nodes have committed to the same block.
	if len(blockAnnounces) >= 4 {
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
