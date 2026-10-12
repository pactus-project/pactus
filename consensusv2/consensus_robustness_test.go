package consensusv2

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// proposeFailingState is a state that is not able to propose a block.
type proposeFailingState struct {
	state.State

	err error
}

func (s *proposeFailingState) ProposeBlock(*bls.ValidatorKey, crypto.Address) (*block.Block, error) {
	return nil, s.err
}

// Based on PIP-19, a proposer with a low availability score is skipped
// without waiting for its proposal.
func TestProposeLowAvailabilityScore(t *testing.T) {
	t.Run("the score of the proposer is above the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.5
		td.stateX.FakeScore = 0.9

		td.enterNewHeight(td.consX)

		td.shouldPublishAProposal(t, td.consX)
	})

	t.Run("the score of the proposer is below the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.9
		td.stateX.FakeScore = 0.5

		// ConsX is the proposer of this round.
		td.enterNewHeight(td.consX)

		td.shouldNotPublish(t, td.consX, message.TypeProposal)
		td.shouldPublishVote(t, td.consX, vote.VoteTypeCPPreVote, hash.UndefHash)
	})

	t.Run("the score is equal to the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.75
		td.stateX.FakeScore = 0.75

		td.enterNewHeight(td.consX)

		td.shouldPublishAProposal(t, td.consX)
	})
}

// shouldPublishAProposal checks that the consensus instance published a proposal.
func (td *testData) shouldPublishAProposal(t *testing.T, cons *consensusV2) {
	t.Helper()

	for _, consMsg := range td.network {
		if consMsg.sender == cons.valKey.Address() &&
			consMsg.message.Type() == message.TypeProposal {
			return
		}
	}

	require.Fail(t, "should publish a proposal")
}

func TestProposeBlockFailure(t *testing.T) {
	td := setup(t)
	td.consX.bcState = &proposeFailingState{State: td.stateX, err: errors.New("unable to propose a block")}

	// ConsX is the proposer, but it can't make a block.
	td.enterNewHeight(td.consX)

	td.shouldNotPublish(t, td.consX, message.TypeProposal)
	assert.Nil(t, td.consX.Proposal())
	assert.Equal(t, "precommit", td.consX.currentState.name(), "it should wait for the proposal of others")
}

func TestCommitBlockFailure(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := types.Height(2)
	round := types.Round(0)

	td.stateP.ErrCommit = errors.New("unable to commit the block")

	td.enterNewHeight(td.consP)
	prop := td.makeProposal(t, height, round)
	blockHash := prop.Block().Hash()

	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexX)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexY)
	td.addPrecommitVote(t, td.consP, blockHash, height, round, tIndexB)
	td.consP.SetProposal(prop)

	// The block is not announced, if it is not committed.
	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, blockHash)
	td.shouldNotPublish(t, td.consP, message.TypeBlockAnnounce)
	assert.Equal(t, "new-height", td.consP.currentState.name())
}

func TestConsensusKey(t *testing.T) {
	td := setup(t)

	assert.Equal(t, td.valKeys[tIndexX].PublicKey(), td.consX.ConsensusKey())
	assert.Equal(t, td.valKeys[tIndexP].PublicKey(), td.consP.ConsensusKey())
	assert.False(t, td.consX.IsDeprecated())
}

func TestTickerString(t *testing.T) {
	assert.Equal(t, "new-height", tickerTargetNewHeight.String())
	assert.Equal(t, "change-proposer", tickerTargetChangeProposer.String())
	assert.Equal(t, "query-proposal", tickerTargetQueryProposal.String())
	assert.Equal(t, "query-vote", tickerTargetQueryVote.String())
	assert.Equal(t, "Unknown", tickerTarget(100).String())

	ti := ticker{Duration: 5 * time.Second, Height: 10, Round: 2, Target: tickerTargetQueryVote}
	assert.Equal(t, "5s@ 10/2/query-vote", ti.LogString())
}

// The justification checkers don't trust the type of the justification they receive.
func TestCPCheckJustWithUnexpectedType(t *testing.T) {
	td := setup(t)
	changeProp := td.consX.changeProposer
	h := td.RandHash()

	checks := map[string]func(just vote.Just) error{
		"init-no": func(just vote.Just) error {
			return changeProp.cpCheckJustInitNo(just, h, 0, vote.CPValueNo)
		},
		"init-yes": func(just vote.Just) error {
			return changeProp.cpCheckJustInitYes(just, hash.UndefHash, 0, vote.CPValueYes)
		},
		"pre-vote-hard": func(just vote.Just) error {
			return changeProp.cpCheckJustPreVoteHard(just, h, 1, vote.CPValueNo)
		},
		"pre-vote-soft": func(just vote.Just) error {
			return changeProp.cpCheckJustPreVoteSoft(just, h, 1, vote.CPValueNo)
		},
		"main-vote-no-conflict": func(just vote.Just) error {
			return changeProp.cpCheckJustMainVoteNoConflict(just, h, 0, vote.CPValueNo)
		},
		"main-vote-conflict": func(just vote.Just) error {
			return changeProp.cpCheckJustMainVoteConflict(just, h, 0, vote.CPValueAbstain)
		},
		"decided": func(just vote.Just) error {
			return changeProp.cpCheckJustDecide(just, h, 0, vote.CPValueNo)
		},
	}

	// Each checker receives a justification that belongs to another checker.
	wrongJusts := map[string]vote.Just{
		"init-no":               &vote.JustInitYes{},
		"init-yes":              &vote.JustInitNo{},
		"pre-vote-hard":         &vote.JustPreVoteSoft{},
		"pre-vote-soft":         &vote.JustPreVoteHard{},
		"main-vote-no-conflict": &vote.JustDecided{},
		"main-vote-conflict":    &vote.JustMainVoteNoConflict{},
		"decided":               &vote.JustInitNo{},
	}

	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			var err error
			require.NotPanics(t, func() { err = check(wrongJusts[name]) })
			require.ErrorIs(t, err, InvalidJustificationError{Reason: "invalid just data"})
		})
	}
}

func TestMediator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	sender := consensus.NewMockConsensus(ts.MockController())
	other1 := consensus.NewMockConsensus(ts.MockController())
	other2 := consensus.NewMockConsensus(ts.MockController())

	testMediator := newConcreteMediator()
	testMediator.Register(sender)
	testMediator.Register(other1)
	testMediator.Register(other2)

	t.Run("proposal is set to the other instances", func(_ *testing.T) {
		prop := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
		other1.EXPECT().SetProposal(prop).Times(1)
		other2.EXPECT().SetProposal(prop).Times(1)

		testMediator.OnPublishProposal(sender, prop)
	})

	t.Run("vote is added to the other instances", func(_ *testing.T) {
		vte, _ := ts.GenerateTestPrepareVote(ts.RandHeight(), ts.RandRound())
		other1.EXPECT().AddVote(vte).Times(1)
		other2.EXPECT().AddVote(vte).Times(1)

		testMediator.OnPublishVote(sender, vte)
	})

	t.Run("block announcement moves the other instances to the new height", func(_ *testing.T) {
		other1.EXPECT().MoveToNewHeight().Times(1)
		other2.EXPECT().MoveToNewHeight().Times(1)

		testMediator.OnBlockAnnounce(sender)
	})
}
