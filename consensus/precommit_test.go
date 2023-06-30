package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, 2, 0)

	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	td.shouldPublishQueryProposal(t, td.consP, 2, 0)
}

func TestPrecommitInvalidProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	p1 := td.makeProposal(t, 2, 0)
	trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, td.signers[0].Address(),
		td.signers[1].Address(), 1000, 1000, "invalid proposal")
	td.signers[0].SignMsg(trx)
	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, 2, 0)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.enterNewHeight(td.consP)

	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexB)

	td.consP.SetProposal(p2)
	assert.Nil(t, td.consP.RoundProposal(0))

	td.consP.SetProposal(p1)
	assert.NotNil(t, td.consP.RoundProposal(0))
}

// Np is partitioned by Nb and goes into the change-proposer state.
// Nx receives prepare votes from Ny and Nb and moves to the precommit state.
// However, Nb doesn't broadcast its precommit vote.
// Once the partition heals, Nx should move to the next round.
func TestPrecommitTimeout(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	p := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(p)

	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p.Block().Hash())
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, p.Block().Hash())

	// Nx and Ny timeout and broadcast change-proposer.
	td.addVote(td.consX, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)
	td.shouldPublishVote(t, td.consX, vote.VoteTypeChangeProposer, hash.UndefHash)

	// partition heals.
	td.addVote(td.consX, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexP)

	// Nx moves to the next round.
	td.checkHeightRoundWait(t, td.consX, 2, 1)
}
