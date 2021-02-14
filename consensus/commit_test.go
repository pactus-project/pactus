package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestResetBooleanOnNewHeight(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	assert.Equal(t, tConsX.hrs.Height(), 2)
	assert.Equal(t, tConsX.hrs.Round(), 0)
	assert.False(t, tConsX.isPrepared)
	assert.False(t, tConsX.isPreCommitted)
	assert.False(t, tConsX.isCommitted)
}

func TestResetBooleanOnNewRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsX.enterNewHeight()
	tConsX.isProposed = true
	tConsX.isPrepared = true
	tConsX.isPreCommitted = true
	tConsX.isCommitted = true

	tConsX.enterNewRound(1)

	assert.Equal(t, tConsX.hrs.Height(), 2)
	assert.Equal(t, tConsX.hrs.Round(), 1)
	assert.False(t, tConsX.isPrepared)
	assert.False(t, tConsX.isPreCommitted)
	assert.False(t, tConsX.isCommitted)
}

func TestEnterCommitAllFailed(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 1
	p1 := makeProposal(t, h, r)

	tConsY.enterNewHeight()
	tConsY.enterNewRound(1)

	// Invalid round
	tConsY.enterCommit(0)
	assert.False(t, tConsY.isCommitted)

	// No quorum
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexP, false)

	v3 := vote.NewPrecommit(h, r, crypto.UndefHash, tSigners[tIndexB].Address())
	tSigners[tIndexB].SignMsg(v3)
	ok, _ := tConsY.pendingVotes.AddVote(v3)
	assert.True(t, ok)

	// Still no quorum
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	testAddVote(t, tConsY, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexB, false)

	// No proposal
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)
	shouldPublishQueryProposal(t, tConsY, h, r)

	// Invalid proposal
	trx := tx.NewSendTx(crypto.UndefHash, 1, tSigners[tIndexX].Address(), tSigners[tIndexY].Address(), 1000, 1000, "")
	tSigners[tIndexX].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx)) // This will change block
	b2, err := tConsY.state.ProposeBlock(0)  // Propose again
	require.NoError(t, err)
	assert.NotEqual(t, b2.Hash(), p1.Block().Hash())
	p2 := vote.NewProposal(h, r, *b2)
	tSigners[tIndexX].SignMsg(p2)
	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p2)

	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)

	// Valid proposal but committing block will fail (no transaction)
	tConsY.pendingVotes.SetRoundProposal(p2.Round(), p1)
	tTxPool.Txs = make([]*tx.Tx, 0)
	tConsY.enterCommit(1)
	assert.False(t, tConsY.isCommitted)
}

func TestSetStaleProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	p := makeProposal(t, 2, 0)

	commitBlockForAllStates(t)

	tConsX.SetProposal(p)
}
