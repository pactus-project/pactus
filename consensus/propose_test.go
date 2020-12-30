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
	setup(t)

	tConsX.enterNewHeight()
	tConsP.enterNewHeight()
	p := tConsX.LastProposal()

	commitFirstBlock(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.LastProposal())
}

func TestSecondProposalCommitted(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsY.enterNewHeight()
	tConsY.enterNewRound(1)

	p1 := tConsX.LastProposal()
	p2 := tConsY.LastProposal() // valid proposal for second round

	tConsX.SetProposal(p1)
	tConsX.SetProposal(p2)

	assert.NotNil(t, tConsX.pendingVotes.RoundProposal(0))
	assert.NotNil(t, tConsX.pendingVotes.RoundProposal(1))

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 1, p2.Block().Hash(), tIndexP, false)

	precommits0 := tConsX.pendingVotes.PrecommitVoteSet(0)
	precommits1 := tConsX.pendingVotes.PrecommitVoteSet(1)
	assert.Equal(t, precommits0.Len(), 0)
	require.NotNil(t, precommits1)
	assert.Equal(t, precommits1.Len(), 3)
	assert.Equal(t, precommits1.ToCommit().Round(), 1)
}

func TestNetworkLagging1(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p1 := tConsX.LastProposal()
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p1)

	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrepare)
	shouldPublishProposalReqquest(t, tConsP)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexB, false)

	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)

	// Proposal received now, set it
	tConsP.SetProposal(p1)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p1.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
}

func TestNetworkLagging2(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p1 := tConsX.LastProposal()
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p1)

	// Networks lags and we don't receive prepare from val_1 and pre-commit from val_4
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexY, false)

	checkHRS(t, tConsP, 1, 0, hrs.StepTypePropose)
	assert.Nil(t, tConsP.pendingVotes.roundVotes[0].Precommits.QuorumBlock())

	shouldPublishProposalReqquest(t, tConsP)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Proposal received now, set it
	tConsP.SetProposal(p1)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrepare)
	// We can't go to precommit stage, because we haven't prepared yet
	// But if we receive another vote we go to commit phase directly
	// Let's do it
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexB, false)
	checkHRSWait(t, tConsP, 2, 0, hrs.StepTypePropose)
}
