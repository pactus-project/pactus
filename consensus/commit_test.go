package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestCommitExecute(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	h := uint32(4)
	r := int16(0)
	p1 := td.makeProposal(t, h, r)
	trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, td.signers[0].Address(),
		td.signers[1].Address(), 1000, 1000, "proposal changer")
	td.signers[0].SignMsg(trx)
	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.enterNewHeight(td.consX)

	td.addVote(td.consX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexX)
	td.addVote(td.consX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexB)

	s := &commitState{td.consX}

	// No proposal
	td.consX.lk.Lock()
	s.decide()
	td.consX.lk.Unlock()
	td.checkHeightRound(t, td.consX, h, r)

	// Invalid proposal
	td.consX.SetProposal(p2)
	td.consX.lk.Lock()
	s.decide()
	td.consX.lk.Unlock()
	assert.Nil(t, td.consX.RoundProposal(0))

	td.consX.SetProposal(p1)
	txs := td.txPool.Txs
	td.txPool.Txs = []*tx.Tx{}
	td.consX.lk.Lock()
	s.decide()
	td.consX.lk.Unlock()
	assert.NotNil(t, td.consX.RoundProposal(0))
	td.checkHeightRound(t, td.consX, h, r)

	v := vote.NewVote(vote.VoteTypePrecommit, h, r, p1.Block().Hash(), td.signers[tIndexP].Address())
	td.signers[tIndexP].SignMsg(v)
	s.onAddVote(v)
	assert.Contains(t, td.consX.AllVotes(), v)

	td.txPool.Txs = txs
	td.consX.lk.Lock()
	s.decide()
	td.consX.lk.Unlock()

	td.shouldPublishBlockAnnounce(t, td.consX, p1.Block().Hash())
}
