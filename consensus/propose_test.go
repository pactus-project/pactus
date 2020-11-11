package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/vote"
)

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	b1 := cons.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	pvals[0].SignMsg(p1)

	// Proposal is valid here
	assert.NoError(t, p1.Verify(cons.proposer(0).PublicKey()))

	commitFirstBlock(t, cons.state)

	// Same proposal is invalid here. because proposer is moved to next validator,
	assert.Error(t, p1.Verify(cons.proposer(0).PublicKey()))
}

func TestInvalidProposer(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.updateHeight(1)

	b1 := cons2.state.ProposeBlock()
	invalidProposal := vote.NewProposal(1, 0, b1)
	pvals[0].SignMsg(invalidProposal)

	validProposal := vote.NewProposal(1, 1, b1)
	pvals[1].SignMsg(validProposal)

	cons1.SetProposal(invalidProposal)
	cons1.SetProposal(validProposal)

	assert.Nil(t, cons1.votes.RoundProposal(0))
	assert.NotNil(t, cons1.votes.RoundProposal(1))
}

func TestSecondProposalCommitted(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.updateHeight(1)

	b1 := cons1.state.ProposeBlock()
	p1 := vote.NewProposal(1, 0, b1)
	pvals[0].SignMsg(p1)

	b2 := cons2.state.ProposeBlock()
	p2 := vote.NewProposal(1, 1, b2) // valid proposal for second round
	pvals[1].SignMsg(p2)

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
