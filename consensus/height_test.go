package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
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
	// Syncer commit a block and trig consensus
	commitBlockForAllStates(t)

	assert.Equal(t, len(tConsP.RoundVotes(0)), 1)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypePrepare)
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

	tConsP.enterNewHeight()
}

// This situation is very vert unlikely to happend, but it is good to have a test
func TestConsensusBehindState2(t *testing.T) {
	setup(t)

	// Consensus starts here
	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p := tConsX.LastProposal()
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commit a block and trig consensus
	// Add some transaction to change the validator set

	commitBlockForAllStates(t)

	committerHash1 := tConsX.state.ValidatorSet().CommittersHash()

	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	pub := tSigners[tIndexX].PublicKey()

	stamp := tConsX.state.LastBlockHash()
	trx1 := tx.NewBondTx(stamp, 1, tSigners[tIndexX].Address(), pub1, 0, "bond-tx", &pub, nil)
	tSigners[tIndexX].SignMsg(trx1)

	proof1 := priv1.Sign(stamp.RawBytes()).RawBytes()
	trx3 := tx.NewSortitionTx(stamp, 1, addr1, proof1, "sortition-tx", &pub1, nil)
	trx3.SetSignature(priv1.Sign(trx3.SignBytes()))

	assert.NoError(t, tTxPool.AppendTx(trx1))
	assert.NoError(t, tTxPool.AppendTx(trx3))

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	committerHash2 := tConsX.state.ValidatorSet().CommittersHash()

	assert.NotEqual(t, committerHash1, committerHash2)
	// Now validator set has changed
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
	// We don't get any error here, but the block is not committed again. Check logs.

	tConsP.enterNewHeight()
	assert.Equal(t, tConsP.hrs.Height(), 4)
}
