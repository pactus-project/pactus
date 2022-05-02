package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/vote"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/tx"
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
	trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, tSigners[0].Address(), tSigners[1].Address(), 1000, 1000, "proposal changer")
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

func TestGoToChangeProposerFromPrecommit(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	p := makeProposal(t, 2, 0)

	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
}
