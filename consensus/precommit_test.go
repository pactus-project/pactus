package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitTimedout(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)

	s := &precommitState{tConsY, false}

	// Invalid target
	tConsX.lk.Lock()
	s.timedout(&ticker{Height: 1, Target: tickerTargetPrepare})
	tConsX.lk.Unlock()
	assert.False(t, s.hasTimedout)

	tConsX.lk.Lock()
	s.timedout(&ticker{Height: 1, Target: tickerTargetPrecommit})
	tConsX.lk.Unlock()
	assert.True(t, s.hasTimedout)

	// Add votes calls execute
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexY)
	v := testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexP)

	tConsX.lk.Lock()
	s.voteAdded(v)
	tConsX.lk.Unlock()
	checkHeightRound(t, tConsY, 1, 1)
}

func TestPrecommitGotoNewRound(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)

	s := &precommitState{tConsY, false}

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexP)

	tConsX.lk.Lock()
	s.execute()
	tConsX.lk.Unlock()
	checkHeightRound(t, tConsY, 1, 1)
}

func TestPrecommitGotoNewHeight(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)
	testEnterNewHeight(tConsY)
	tConsY.SetProposal(p)

	s := &precommitState{tConsY, false}

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)

	tConsX.lk.Lock()
	s.execute()
	tConsX.lk.Unlock()
	shouldPublishBlockAnnounce(t, tConsY, p.Block().Hash())
}

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

func TestPrecommitNullVote1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexB)

	s := &precommitState{tConsP, false}
	tConsX.lk.Lock()
	s.vote()
	tConsX.lk.Unlock()
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestPrecommitNullVote2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.UndefHash, tIndexB)

	s := &precommitState{tConsP, false}
	tConsX.lk.Lock()
	s.vote()
	tConsX.lk.Unlock()
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
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

	s := &precommitState{tConsP, false}
	tConsP.SetProposal(p2)

	assert.NotNil(t, tConsP.RoundProposal(0))
	tConsX.lk.Lock()
	s.vote()
	tConsX.lk.Unlock()
	assert.Nil(t, tConsP.RoundProposal(0))
}
