package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/vote"
)

func waitForNewHeight(t *testing.T, cons *consensus, expectedHeight int) {
	for i := 0; i < 10; i++ {
		if expectedHeight == cons.HRS().Height() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.Equal(t, expectedHeight, cons.hrs.Height())
}

func TestMoveToNewHeight(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.MoveToNewHeight()
	waitForNewHeight(t, cons, 1)

	commitFirstBlock(t, cons.state)

	cons.MoveToNewHeight()
	waitForNewHeight(t, cons, 2)
}

func TestConsensusBehindState(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	commitFirstBlock(t, cons.state)
	assert.Equal(t, cons.hrs, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight))

	cons.MoveToNewHeight()
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))

	// Calling MoveToNewHeight for the second time
	cons.MoveToNewHeight()
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))
}

func TestConsensusBehindState3(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	// Consensus starts here
	cons.enterNewHeight(1)
	p := cons.LastProposal()
	b := p.Block()
	assert.NoError(t, cons.state.ValidateBlock(b))

	// --------------------------------
	// Syncer commit a block and trig consensus
	commitFirstBlock(t, cons.state)
	cons.MoveToNewHeight()

	assert.Equal(t, len(cons.votes.votes), 1)
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))

	cons.MoveToNewHeight()
	assert.Equal(t, len(cons.votes.votes), 1)
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL2, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL4, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)
	precommits := cons.votes.PrecommitVoteSet(0)

	assert.Error(t, cons.state.ValidateBlock(b))

	assert.NoError(t, cons.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))
	// We don't get any error here, but the block is not committed again, Check the log
}
