package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/vote"
)

func commitFirstBlock(t *testing.T, st state.State) (b block.Block, votes [3]*vote.Vote) {
	b = st.ProposeBlock()

	votes[0] = vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), tSigners[0].Address())
	tSigners[0].SignMsg(votes[0])

	votes[1] = vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), tSigners[1].Address())
	tSigners[1].SignMsg(votes[1])

	votes[2] = vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), tSigners[2].Address())
	tSigners[2].SignMsg(votes[2])

	sig := crypto.Aggregate([]*crypto.Signature{votes[0].Signature(), votes[1].Signature(), votes[2].Signature()})
	c := block.NewCommit(0,
		[]block.Committer{
			{Status: 1, Address: tSigners[0].Address()},
			{Status: 1, Address: tSigners[1].Address()},
			{Status: 1, Address: tSigners[2].Address()},
			{Status: 0, Address: tSigners[3].Address()},
		},
		sig)

	require.NotNil(t, c)
	err := st.ApplyBlock(1, b, *c)
	assert.NoError(t, err)

	return
}

func TestInvalidStepAfterBlockCommit(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	commitFirstBlock(t, cons.state)

	cons.MoveToNewHeight()

	assert.True(t, cons.invalidHeight(1))
	assert.True(t, cons.invalidHeightRound(1, 0))
	assert.True(t, cons.invalidHeightRoundStep(1, 0, hrs.StepTypeCommit))

	// manually move to next height
	cons.enterNewHeight(2)

	assert.False(t, cons.invalidHeight(2))
	assert.False(t, cons.invalidHeightRound(2, 0))
	assert.False(t, cons.invalidHeightRoundStep(2, 0, hrs.StepTypeCommit))
}

func TestEnterCommit(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons2 := newTestConsensus(t, VAL2)

	cons1.MoveToNewHeight()
	cons2.MoveToNewHeight()
	checkHRSWait(t, cons1, 1, 0, hrs.StepTypePrepare)
	checkHRSWait(t, cons2, 1, 0, hrs.StepTypePrepare)
	p1 := cons1.LastProposal()

	// Invalid height
	cons2.enterCommit(2, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrepare)

	// No quorum
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, cons2, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), VAL1, false)
	testAddVote(t, cons2, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), VAL2, false)

	v3 := vote.NewPrecommit(1, 0, crypto.UndefHash, tSigners[VAL3].Address())
	tSigners[VAL3].SignMsg(v3)
	ok, _ := cons2.votes.AddVote(v3)
	assert.True(t, ok)

	// Undef quorum
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrepare)

	v4 := vote.NewPrecommit(1, 0, p1.Block().Hash(), tSigners[VAL4].Address())
	tSigners[VAL4].SignMsg(v4)
	ok, _ = cons2.votes.AddVote(v4)
	assert.True(t, ok)

	// No proposal
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrepare)
	shouldPublishProposalReqquest(t, cons2)

	time.Sleep(1 * time.Second) // This will change block timestamp
	b2 := cons1.state.ProposeBlock()
	p2 := vote.NewProposal(1, 0, b2)
	tSigners[VAL1].SignMsg(p2)
	cons2.votes.SetRoundProposal(p2.Round(), p2)

	// Invalid proposal
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypePrepare)

	cons2.votes.SetRoundProposal(p2.Round(), p1)

	// Everything is good
	cons2.enterCommit(1, 0)
	checkHRS(t, cons2, 1, 0, hrs.StepTypeCommit)
}
