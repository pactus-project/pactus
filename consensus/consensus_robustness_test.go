package consensus

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
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

// Based on PIP-19, the Change-Proposer phase starts if the proposer has a low availability score.
func TestProposeLowAvailabilityScore(t *testing.T) {
	t.Run("the score of the proposer is above the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.5
		td.stateX.FakeScore = 0.9

		td.enterNewHeight(td.consX)

		assert.Equal(t, "prepare", td.consX.currentState.name())
	})

	t.Run("the score of the proposer is below the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.9
		td.stateX.FakeScore = 0.5

		// ConsX is the proposer of this round.
		td.enterNewHeight(td.consX)

		// It doesn't wait for the others: the Change-Proposer phase is started.
		assert.NotEqual(t, "prepare", td.consX.currentState.name())
		assert.NotNil(t, td.consX.Proposal(), "the proposer already made its proposal in this version")
	})

	t.Run("the score is equal to the minimum", func(t *testing.T) {
		td := setup(t)
		td.consX.config.MinimumAvailabilityScore = 0.75
		td.stateX.FakeScore = 0.75

		td.enterNewHeight(td.consX)

		assert.Equal(t, "prepare", td.consX.currentState.name())
	})
}

func TestProposeBlockFailure(t *testing.T) {
	td := setup(t)
	td.consX.bcState = &proposeFailingState{State: td.stateX, err: errors.New("unable to propose a block")}

	// ConsX is the proposer, but it can't make a block.
	td.enterNewHeight(td.consX)

	td.shouldNotPublish(t, td.consX, message.TypeProposal)
	assert.Nil(t, td.consX.Proposal())
	assert.Equal(t, "prepare", td.consX.currentState.name(), "it should wait for the proposal of others")
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

	td.addPrepareVote(td.consP, blockHash, height, round, tIndexX)
	td.addPrepareVote(td.consP, blockHash, height, round, tIndexY)
	td.addPrepareVote(td.consP, blockHash, height, round, tIndexB)
	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexX)
	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexY)
	td.addPrecommitVote(td.consP, blockHash, height, round, tIndexB)
	td.consP.SetProposal(prop)

	// The consensus goes on even if the state is not able to commit the block.
	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, blockHash)
	assert.Equal(t, "new-height", td.consP.currentState.name())
}

func TestConsensusKey(t *testing.T) {
	td := setup(t)

	assert.Equal(t, td.valKeys[tIndexX].PublicKey(), td.consX.ConsensusKey())
	assert.Equal(t, td.valKeys[tIndexP].PublicKey(), td.consP.ConsensusKey())
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
			return changeProp.checkJustInitNo(just, h)
		},
		"init-yes": changeProp.checkJustInitYes,
		"pre-vote-hard": func(just vote.Just) error {
			return changeProp.checkJustPreVoteHard(just, h, 1, vote.CPValueNo)
		},
		"pre-vote-soft": func(just vote.Just) error {
			return changeProp.checkJustPreVoteSoft(just, h, 1)
		},
		"main-vote-no-conflict": func(just vote.Just) error {
			return changeProp.checkJustMainVoteNoConflict(just, h, 0, vote.CPValueNo)
		},
		"main-vote-conflict": func(just vote.Just) error {
			return changeProp.checkJustMainVoteConflict(just, h, 0)
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
	}

	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			wrongJust := wrongJusts[name]

			var err error
			require.NotPanics(t, func() { err = check(wrongJust) })
			require.ErrorIs(t, err, InvalidJustificationError{
				JustType: wrongJust.Type(),
				Reason:   "invalid just data",
			})
		})
	}
}

func TestMediator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	sender := NewMockConsensus(ts.MockController())
	other1 := NewMockConsensus(ts.MockController())
	other2 := NewMockConsensus(ts.MockController())

	testMediator := NewConcreteMediator()
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
