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

	sb := vote.CommitSignBytes(b.Hash(), 0)
	sig1 := tSigners[0].Sign(sb)
	sig2 := tSigners[1].Sign(sb)
	sig3 := tSigners[2].Sign(sb)

	sig := crypto.Aggregate([]*crypto.Signature{sig1, sig2, sig3})
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
	setup(t)

	commitFirstBlock(t, tConsY.state)

	tConsY.enterNewHeight()

	assert.Equal(t, tConsY.hrs.Height(), 2)
	assert.Equal(t, tConsY.hrs.Round(), 0)
	assert.False(t, tConsX.isProposed)
	assert.False(t, tConsX.isPrepared)
	assert.False(t, tConsX.isPreCommitted)
	assert.False(t, tConsX.isCommitted)
}

func TestEnterCommit(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsY.enterNewHeight()
	p1 := tConsX.LastProposal()

	// Invalid round
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	// No quorum
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexP, false)

	v3 := vote.NewPrecommit(1, 0, crypto.UndefHash, tSigners[tIndexB].Address())
	tSigners[tIndexB].SignMsg(v3)
	ok, _ := tConsY.pendingVotes.AddVote(v3)
	assert.True(t, ok)

	// Undef quorum
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexB, false)

	// No proposal
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)
	shouldPublishProposalReqquest(t, tConsY)

	time.Sleep(2 * time.Second) // This will change block timestamp
	b2 := tConsX.state.ProposeBlock()
	assert.NotEqual(t, b2.Hash(), p1.Block().Hash())
	p2 := vote.NewProposal(1, 0, b2)
	tSigners[tIndexX].SignMsg(p2)
	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p2)

	// Invalid proposal
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)

	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p1)

	// Everything is good
	tConsY.enterCommit(0)
	assert.True(t, tConsY.isCommitted)

	checkHRSWait(t, tConsY, 2, 0, hrs.StepTypePrepare)
}
