package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestCommitExecute(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	p1 := makeProposal(t, 2, 0)
	trx := tx.NewSendTx(crypto.UndefHash, 1, tSigners[0].Address(), tSigners[1].Address(), 1000, 1000, "proposal changer")
	tSigners[0].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx))
	p2 := makeProposal(t, 2, 0)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	testEnterNewHeight(tConsX)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p1.Block().Hash(), tIndexX)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p1.Block().Hash(), tIndexP)

	s := &commitState{tConsX}

	// No proposal
	tConsX.lk.Lock()
	s.decide()
	tConsX.lk.Unlock()
	checkHeightRound(t, tConsX, 2, 0)

	// Invalid proposal
	tConsX.SetProposal(p2)
	tConsX.lk.Lock()
	s.decide()
	tConsX.lk.Unlock()
	assert.Nil(t, tConsX.RoundProposal(0))

	tConsX.SetProposal(p1)
	txs := tTxPool.Txs
	tTxPool.Txs = []*tx.Tx{}
	tConsX.lk.Lock()
	s.decide()
	tConsX.lk.Unlock()
	assert.NotNil(t, tConsX.RoundProposal(0))
	checkHeightRound(t, tConsX, 2, 0)

	tTxPool.Txs = txs
	tConsX.lk.Lock()
	s.decide()
	tConsX.lk.Unlock()
	shouldPublishBlockAnnounce(t, tConsX, p1.Block().Hash())
}
