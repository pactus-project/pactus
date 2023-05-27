package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	p := makeProposal(t, 2, 0)

	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	shouldPublishQueryProposal(t, tConsP, 2, 0)
}

func TestPrecommitInvalidProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	p1 := makeProposal(t, 2, 0)
	trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, tSigners[0].Address(),
		tSigners[1].Address(), 1000, 1000, "invalid proposal")
	tSigners[0].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx))
	p2 := makeProposal(t, 2, 0)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	testEnterNewHeight(tConsP)

	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexB)

	tConsP.SetProposal(p2)
	assert.Nil(t, tConsP.RoundProposal(0))

	tConsP.SetProposal(p1)
	assert.NotNil(t, tConsP.RoundProposal(0))
}

// Np is partitioned by Nb and goes into the change-proposer state.
// Nx receives prepare votes from Ny and Nb and moves to the precommit state.
// However, Nb doesn't broadcast its precommit vote.
// Once the partition heals, Nx should move to the next round.
func TestPrecommitTimeout(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	p := makeProposal(t, 2, 0)
	tConsX.SetProposal(p)

	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())
	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p.Block().Hash())

	// Nx and Ny timeout and broadcast change-proposer.
	testAddVote(tConsX, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)
	shouldPublishVote(t, tConsX, vote.VoteTypeChangeProposer, hash.UndefHash)

	// partition heals.
	testAddVote(tConsX, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexP)

	// Nx moves to the next round.
	checkHeightRoundWait(t, tConsX, 2, 1)
}
