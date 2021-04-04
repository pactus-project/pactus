package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestCommitExecute(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p1 := makeProposal(t, h, r)
	p2, _ := proposal.GenerateTestProposal(2, 0)

	testEnterNewHeight(tConsX)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexX)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexP)

	s := &commitState{tConsX}

	// No proposal
	s.execute()
	checkHeightRound(t, tConsX, 2, 0)

	// Invalid proposal
	tConsX.SetProposal(p2)
	s.execute()
	assert.Nil(t, tConsX.RoundProposal(0))

	tConsX.SetProposal(p1)
	txs := tTxPool.Txs
	tTxPool.Txs = []*tx.Tx{}
	s.execute()
	assert.NotNil(t, tConsX.RoundProposal(0))
	checkHeightRound(t, tConsX, 2, 0)

	tTxPool.Txs = txs
	s.execute()
	shouldPublishBlockAnnounce(t, tConsX, p1.Block().Hash())
}
