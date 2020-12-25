package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	cons := newTestConsensus(t, VAL1)
	p := makeTestProposal(t, VAL1, 1, 0)

	// Proposal is valid here
	assert.NoError(t, p.Verify(cons.proposer(0).PublicKey()))

	commitFirstBlock(t, cons.state)

	// Same proposal is invalid here. because proposer is moved to next validator,
	assert.Error(t, p.Verify(cons.proposer(0).PublicKey()))
}

func TestInvalidProposer(t *testing.T) {
	cons2 := newTestConsensus(t, VAL2)
	cons3 := newTestConsensus(t, VAL3)

	cons3.enterNewHeight(1)

	b1 := cons2.state.ProposeBlock()
	invalidProposal := vote.NewProposal(1, 0, b1)
	tSigners[VAL1].SignMsg(invalidProposal)

	validProposal := makeTestProposal(t, VAL2, 1, 1)

	cons3.SetProposal(invalidProposal)
	cons3.SetProposal(validProposal)

	assert.Nil(t, cons3.votes.RoundProposal(0))
	assert.NotNil(t, cons3.votes.RoundProposal(1))
}

func TestSecondProposalCommitted(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)

	cons1.enterNewHeight(1)

	p1 := makeTestProposal(t, VAL1, 1, 0)
	p2 := makeTestProposal(t, VAL2, 1, 1) // valid proposal for second round

	cons1.SetProposal(p1)
	cons1.SetProposal(p2)

	assert.NotNil(t, cons1.votes.RoundProposal(0))
	assert.NotNil(t, cons1.votes.RoundProposal(1))

	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), VAL2, false)
	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), VAL3, false)
	testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), VAL4, false)

	precommits0 := cons1.votes.PrecommitVoteSet(0)
	precommits1 := cons1.votes.PrecommitVoteSet(1)
	assert.Equal(t, precommits0.Len(), 0)
	require.NotNil(t, precommits1)
	assert.Equal(t, precommits1.Len(), 3)
	assert.Equal(t, precommits1.ToCommit().Round(), 1)
}

func TestNetworkLagging1(t *testing.T) {
	cons2 := newTestConsensus(t, VAL2)

	cons2.enterNewHeight(1)

	p1 := makeTestProposal(t, VAL1, 1, 0)
	// We don't set proposal for second validator here
	// cons2.SetProposal(p1)

	assert.Nil(t, cons2.votes.RoundProposal(0))

	v1 := testAddVote(t, cons2, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), VAL1, false)
	v3 := testAddVote(t, cons2, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), VAL3, false)

	checkHRSWait(t, cons2, 1, 0, hrs.StepTypePropose)
	cons2.enterPrecommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)

	assert.NoError(t, cons2.addVote(v1))
	assert.NoError(t, cons2.addVote(v3))
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)
	assert.Nil(t, cons2.votes.roundVoteSets[0].Prepares.QuorumBlock())

	shouldPublishProposalReqquest(t, cons2)
	shouldPublishVote(t, cons2, vote.VoteTypePrepare, crypto.UndefHash)

	// Proposal received now, set it
	cons2.SetProposal(p1)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrecommit)

	shouldPublishVote(t, cons2, vote.VoteTypePrecommit, p1.Block().Hash())
	shouldPublishVote(t, cons2, vote.VoteTypePrepare, p1.Block().Hash())
}

func TestNetworkLagging2(t *testing.T) {
	cons2 := newTestConsensus(t, VAL2)

	cons2.enterNewHeight(1)

	p1 := makeTestProposal(t, VAL1, 1, 0)

	// We don't set proposal for second validator here
	// cons2.SetProposal(p1)

	assert.Nil(t, cons2.votes.RoundProposal(0))

	precommit1 := testAddVote(t, cons2, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), VAL1, false)
	precommit3 := testAddVote(t, cons2, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), VAL3, false)

	checkHRSWait(t, cons2, 1, 0, hrs.StepTypePropose)
	cons2.enterPrecommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)

	// Networks lags and we don't receive prepare from val_1 and pre-commit from val_4
	assert.NoError(t, cons2.addVote(precommit1))
	assert.NoError(t, cons2.addVote(precommit3))
	checkHRS(t, cons2, 1, 0, hrs.StepTypePropose)
	assert.Nil(t, cons2.votes.roundVoteSets[0].Precommits.QuorumBlock())

	shouldPublishProposalReqquest(t, cons2)
	shouldPublishVote(t, cons2, vote.VoteTypePrepare, crypto.UndefHash)

	// Proposal received now, set it
	cons2.SetProposal(p1)

	shouldPublishVote(t, cons2, vote.VoteTypePrepare, p1.Block().Hash())
}
