package consensus

import (
	"testing"

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
