package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/vote"
)

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	b1 := cons.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(p1)

	// Proposal is valid here
	assert.NoError(t, p1.Verify(cons.proposer(0).PublicKey()))

	commitFirstBlock(t, cons.state)

	// Same proposal is invalid here. because proposer is moved to next validator,
	assert.Error(t, p1.Verify(cons.proposer(0).PublicKey()))
}

func TestInvalidProposer(t *testing.T) {
	cons2 := newTestConsensus(t, VAL2)
	cons3 := newTestConsensus(t, VAL3)

	cons3.enterNewHeight(1)

	b1 := cons2.state.ProposeBlock()
	invalidProposal := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(invalidProposal)

	validProposal := vote.NewProposal(1, 1, b1)

	signers[1].SignMsg(validProposal)

	cons3.SetProposal(invalidProposal)
	cons3.SetProposal(validProposal)

	assert.Nil(t, cons3.votes.RoundProposal(0))
	assert.NotNil(t, cons3.votes.RoundProposal(1))
}

func TestSecondProposalCommitted(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.enterNewHeight(1)

	b1 := cons1.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(p1)

	b2 := cons2.state.ProposeBlock()
	p2 := vote.NewProposal(1, 1, b2) // valid proposal for second round
	signers[1].SignMsg(p2)

	cons1.SetProposal(p1)
	cons1.SetProposal(p2)

	assert.NotNil(t, cons1.votes.RoundProposal(0))
	assert.NotNil(t, cons1.votes.RoundProposal(1))

	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, b2.Hash(), VAL2, false)
	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, b2.Hash(), VAL3, false)
	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, b2.Hash(), VAL4, false)

	precommits0 := cons1.votes.Precommits(0)
	precommits1 := cons1.votes.Precommits(1)
	assert.Equal(t, precommits0.Len(), 0)
	require.NotNil(t, precommits1)
	assert.Equal(t, precommits1.Len(), 3)
	assert.Equal(t, precommits1.ToCommit().Round(), 1)
}

func TestNetworkLagging1(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	// In some slow containers, it goes to prevote stage before setting proposal here
	// Let set the height manually ;)
	//cons1.enterNewHeight(1)
	cons1.updateHeight(1)
	cons2.enterNewHeight(1)

	b1 := cons1.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(p1)

	cons1.SetProposal(p1)
	// We don't set proposal for second validator here
	// cons2.SetProposal(p1)

	assert.NotNil(t, cons1.votes.RoundProposal(0))
	assert.Nil(t, cons2.votes.RoundProposal(0))

	v1 := testAddVote(t, cons1, vote.VoteTypePrevote, 1, 0, b1.Hash(), VAL1, false)
	v3 := testAddVote(t, cons1, vote.VoteTypePrevote, 1, 0, b1.Hash(), VAL3, false)

	checkHRSWait(t, cons2, 1, 0, hrs.StepTypePrevote)
	cons2.enterPrecommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrevote)
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrevote)
	assert.Equal(t, len(cons2.votes.votes), 1) // UndefHash vote

	assert.NoError(t, cons2.AddVote(v1))
	assert.NoError(t, cons2.AddVote(v3))
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrevoteWait)
	assert.Equal(t, len(cons2.votes.votes), 3)
	assert.Nil(t, cons2.votes.roundVoteSets[0].Prevotes.QuorumBlock())

	// Proposal received now, set it
	cons2.SetProposal(p1)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrecommit)
	assert.True(t, cons2.votes.roundVoteSets[0].Prevotes.QuorumBlock().EqualsTo(b1.Hash()))
}

func TestNetworkLagging2(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.enterNewHeight(1)
	cons2.enterNewHeight(1)

	b1 := cons1.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(p1)

	cons1.SetProposal(p1)
	// We don't set proposal for second validator here
	// cons2.SetProposal(p1)

	assert.NotNil(t, cons1.votes.RoundProposal(0))
	assert.Nil(t, cons2.votes.RoundProposal(0))

	prevote3 := testAddVote(t, cons1, vote.VoteTypePrevote, 1, 0, b1.Hash(), VAL3, false)
	prevote4 := testAddVote(t, cons1, vote.VoteTypePrevote, 1, 0, b1.Hash(), VAL4, false)
	precommit1 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, b1.Hash(), VAL1, false)
	precommit3 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, b1.Hash(), VAL3, false)

	checkHRSWait(t, cons2, 1, 0, hrs.StepTypePrevote)
	cons2.enterPrecommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrevote)
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrevote)
	assert.Equal(t, len(cons2.votes.votes), 1) // UndefHash vote

	// Networks lags and we don't receive pre-vote from val_1 and pre-commit from val_4
	assert.NoError(t, cons2.AddVote(precommit1))
	assert.NoError(t, cons2.AddVote(precommit3))
	assert.NoError(t, cons2.AddVote(prevote4))
	assert.NoError(t, cons2.AddVote(prevote3))
	assert.Equal(t, len(cons2.votes.votes), 5)
	assert.Nil(t, cons2.votes.roundVoteSets[0].Precommits.QuorumBlock())

	// Proposal received now, set it
	cons2.SetProposal(p1)

	// Cons3 has enough votes to go to next height
	checkHRSWait(t, cons2, 2, 0, hrs.StepTypePrevote)
}

func TestNetworkLagging3(t *testing.T) {
	// Cons2 goes to next height without receiving any prevotes
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.enterNewHeight(1)
	cons2.enterNewHeight(1)

	b1 := cons1.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	signers[0].SignMsg(p1)

	cons1.SetProposal(p1)
	// We don't set proposal for second validator here
	// cons2.SetProposal(p1)

	assert.NotNil(t, cons1.votes.RoundProposal(0))
	assert.Nil(t, cons2.votes.RoundProposal(0))

	precommit1 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, b1.Hash(), VAL1, false)
	precommit3 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, b1.Hash(), VAL3, false)
	precommit4 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, b1.Hash(), VAL4, false)

	// Networks lags and we don't receive pre-vote from val_1 and pre-commit from val_4
	assert.NoError(t, cons2.AddVote(precommit1))
	assert.NoError(t, cons2.AddVote(precommit3))
	assert.NoError(t, cons2.AddVote(precommit4))
	assert.Equal(t, len(cons2.votes.votes), 3)
	assert.True(t, cons2.votes.roundVoteSets[0].Precommits.QuorumBlock().EqualsTo(b1.Hash()))

	// Proposal received now, set it
	cons2.SetProposal(p1)

	// Cons3 has enough votes to go to next height
	checkHRSWait(t, cons2, 2, 0, hrs.StepTypePrevote)
}
