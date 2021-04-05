package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	shouldPublishQueryProposal(t, tConsP, 2, 0) // prepare stage, ignore it

	p := makeProposal(t, 2, 0)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	s := &precommitState{tConsP, false}
	tConsX.lk.Lock()
	s.vote()
	tConsX.lk.Unlock()
	shouldPublishQueryProposal(t, tConsP, 2, 0)
}

func TestPrecommitInvalidProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	p1 := makeProposal(t, 2, 0)
	trx := tx.NewSendTx(crypto.UndefHash, 1, tSigners[0].Address(), tSigners[1].Address(), 1000, 1000, "proposal changer")
	tSigners[0].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx))
	p2 := makeProposal(t, 2, 0)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	testEnterNewHeight(tConsP)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, p1.Block().Hash(), tIndexB)

	tConsP.SetProposal(p2)
	assert.Nil(t, tConsP.RoundProposal(0))

	tConsP.SetProposal(p1)
	assert.NotNil(t, tConsP.RoundProposal(0))
}
