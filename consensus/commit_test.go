package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func commitFirstBlock(t *testing.T) (b block.Block, votes [3]*vote.Vote) {
	pb, err := tConsX.state.ProposeBlock(0)
	require.NoError(t, err)
	b = *pb

	sb := block.CommitSignBytes(b.Hash(), 0)
	sig1 := tSigners[0].Sign(sb)
	sig2 := tSigners[1].Sign(sb)
	sig3 := tSigners[2].Sign(sb)

	sig := crypto.Aggregate([]*crypto.Signature{sig1, sig2, sig3})
	c := block.NewCommit(b.Hash(), 0, []int{0, 1, 2}, []int{3}, sig)

	require.NotNil(t, c)
	err = tConsX.state.ApplyBlock(1, b, *c)
	assert.NoError(t, err)
	err = tConsY.state.ApplyBlock(1, b, *c)
	assert.NoError(t, err)
	err = tConsB.state.ApplyBlock(1, b, *c)
	assert.NoError(t, err)
	err = tConsP.state.ApplyBlock(1, b, *c)
	assert.NoError(t, err)

	return
}

func TestInvalidStepAfterBlockCommit(t *testing.T) {
	setup(t)

	commitFirstBlock(t)

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

	commitFirstBlock(t)

	tConsY.enterNewHeight()
	tConsY.enterNewRound(1)
	tConsB.enterNewHeight()
	tConsB.enterNewRound(1)
	p1 := tConsB.LastProposal()

	// Invalid round
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)

	// No quorum
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 2, 1, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 2, 1, p1.Block().Hash(), tIndexP, false)

	v3 := vote.NewPrecommit(2, 1, crypto.UndefHash, tSigners[tIndexB].Address())
	tSigners[tIndexB].SignMsg(v3)
	ok, _ := tConsY.pendingVotes.AddVote(v3)
	assert.True(t, ok)

	// Undef quorum
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 2, 1, p1.Block().Hash(), tIndexB, false)

	// No proposal
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)
	shouldPublishProposalReqquest(t, tConsY)

	pub := tSigners[tIndexX].PublicKey()
	trx := tx.NewSendTx(crypto.UndefHash, 1, tSigners[tIndexX].Address(), tSigners[tIndexY].Address(), 1000, 1000, "", &pub, nil)
	tSigners[tIndexX].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx)) // This will change block
	b2, err := tConsY.state.ProposeBlock(0)
	require.NoError(t, err)
	assert.NotEqual(t, b2.Hash(), p1.Block().Hash())
	p2 := vote.NewProposal(2, 1, *b2)
	tSigners[tIndexX].SignMsg(p2)
	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p2)

	// Invalid proposal
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p1)

	// Everything is good
	tConsY.enterCommit(1)
	shouldPublishProposalBlock(t, tConsY)
}
