package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestConsensusBehindState(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	b := st.ProposeBlock()

	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[0].Address())
	pvals[0].SignMsg(v1)

	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[1].Address())
	pvals[1].SignMsg(v2)

	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[2].Address())
	pvals[2].SignMsg(v3)

	c := block.NewCommit(0,
		[]crypto.Address{pvals[0].Address(), pvals[1].Address(), pvals[2].Address()},
		[]crypto.Signature{*v1.Signature(), *v2.Signature(), *v3.Signature()})

	require.NotNil(t, c)
	err := st.ApplyBlock(b, *c)
	assert.NoError(t, err)
	assert.Equal(t, cons.hrs, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight))
	cons.ScheduleNewHeight()
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))

	// Calling ScheduleNewHeight for the second time
	cons.ScheduleNewHeight()
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))
}

func TestConsensusBehindState2(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	cons.enterNewHeight(1)
	p := cons.LastProposal()
	b := p.Block()

	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[0].Address())
	pvals[0].SignMsg(v1)

	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[1].Address())
	pvals[1].SignMsg(v2)

	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b.Hash(), pvals[2].Address())
	pvals[2].SignMsg(v3)

	cons.AddVote(v1)

	c := block.NewCommit(0,
		[]crypto.Address{pvals[0].Address(), pvals[1].Address(), pvals[2].Address()},
		[]crypto.Signature{*v1.Signature(), *v2.Signature(), *v3.Signature()})

	require.NotNil(t, c)
	assert.Equal(t, len(cons.votes.votes), 2)
	err := st.ApplyBlock(b, *c)
	assert.NoError(t, err)
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypePrevote))
	cons.ScheduleNewHeight()
	assert.Equal(t, len(cons.votes.votes), 2)
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))

	cons.ScheduleNewHeight()
	assert.Equal(t, len(cons.votes.votes), 2)
	assert.Equal(t, cons.hrs, hrs.NewHRS(1, 0, hrs.StepTypeCommit))
}
