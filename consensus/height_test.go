package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/vote"
)

func TestMoveToNewHeight(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsP.MoveToNewHeight()
	checkHRSWait(t, tConsP, 2, 0, hrs.StepTypePropose)

	// Calling MoveToNewHeight for the second time
	tConsP.enterNewHeight()
	checkHRSWait(t, tConsP, 2, 0, hrs.StepTypePropose)
}

func TestConsensusBehindState(t *testing.T) {
	setup(t)

	// Consensus starts here
	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p := tConsX.LastProposal()
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commits a block
	commitBlockForAllStates(t)
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP, false)

	precommits := tConsP.pendingVotes.PrecommitVoteSet(0)
	require.NotNil(t, precommits)
	require.NotNil(t, precommits.ToCommit())

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))

	assert.NoError(t, tConsP.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))
	// We don't get any error here, but the block is not committed again. Check logs.
}

func TestConsensusBehindState2(t *testing.T) {
	setup(t)

	// Consensus starts here
	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p := tConsX.LastProposal()
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commits a block and trig consensus
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP, false)

	precommits := tConsP.pendingVotes.PrecommitVoteSet(0)
	require.NotNil(t, precommits)
	require.NotNil(t, precommits.ToCommit())

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))

	assert.Error(t, tConsP.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))
}
