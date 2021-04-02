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

	tConsX.MoveToNewHeight()

	checkHRSWait(t, tConsX, 1, 0, hrs.StepTypePrepare)

	// Calling MoveToNewHeight for the second time
	testEnterNewHeight(tConsX)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrepare)
}

func TestConsensusBehindState(t *testing.T) {
	setup(t)

	// Consensus starts here
	testEnterNewHeight(tConsP)

	p := makeProposal(t, 1, 0)
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commits a block
	commitBlockForAllStates(t)
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)

	precommits := tConsP.pendingVotes.PrecommitVoteSet(0)
	require.NotNil(t, precommits)
	require.NotNil(t, precommits.ToCertificate())

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))

	assert.NoError(t, tConsP.state.CommitBlock(1, p.Block(), *precommits.ToCertificate()))
	// We don't get any error here, but the block is not committed again. Check logs.
}

func TestConsensusBehindState2(t *testing.T) {
	setup(t)

	// Consensus starts here
	testEnterNewHeight(tConsP)

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commits a block and trig consensus
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexP)

	precommits := tConsP.pendingVotes.PrecommitVoteSet(r)
	require.NotNil(t, precommits)
	require.NotNil(t, precommits.ToCertificate())

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))
}
