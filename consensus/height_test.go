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

	tConsP.MoveToNewHeight()

	commitBlockForAllStates(t)

	tConsP.MoveToNewHeight()
	checkHRSWait(t, tConsP, 2, 0, hrs.StepTypePropose)

	// Calling MoveToNewHeight for the second time
	tConsX.MoveToNewHeight()
	checkHRSWait(t, tConsP, 2, 0, hrs.StepTypePropose)
}

func TestConsensusBehindState3(t *testing.T) {
	setup(t)

	// Consensus starts here
	tConsX.enterNewHeight()
	checkHRSWait(t, tConsX, 1, 0, hrs.StepTypePrepare)

	p := tConsX.LastProposal()
	b := p.Block()
	assert.NoError(t, tConsX.state.ValidateBlock(b))

	// --------------------------------
	// Syncer commit a block and trig consensus
	commitBlockForAllStates(t)
	tConsX.MoveToNewHeight()

	assert.Equal(t, len(tConsX.RoundVotes(0)), 1)
	assert.Equal(t, tConsX.hrs, hrs.NewHRS(1, 0, hrs.StepTypePrepare))
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexB, false)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP, false)

	precommits := tConsX.pendingVotes.PrecommitVoteSet(0)
	require.NotNil(t, precommits)
	require.NotNil(t, precommits.ToCommit())

	assert.Error(t, tConsX.state.ValidateBlock(b))

	assert.NoError(t, tConsX.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))
	// We don't get any error here, but the block is not committed again. Check logs.
}
