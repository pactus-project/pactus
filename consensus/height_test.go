package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

func TestNewHeightTimedout(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)
	commitBlockForAllStates(t)

	s := &newHeightState{tConsX}

	// Invalid target
	s.onTimedout(&ticker{Height: 2, Target: 3})
	checkHeightRound(t, tConsX, 1, 0)

	s.onTimedout(&ticker{Height: 2, Target: tickerTargetNewHeight})
	checkHeightRound(t, tConsX, 2, 0)
}

func TestNewHeightDuplicateEntry(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)
	testEnterPropose(tConsX)

	s := &newHeightState{tConsX}

	s.onTimedout(&ticker{Height: 1, Target: tickerTargetNewHeight})
	checkHeightRound(t, tConsX, 1, 1)
}

func TestUpdateCertificate(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	p := makeProposal(t, 2, 0)
	tConsX.SetProposal(p)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexB)

	assert.Equal(t, tConsX.state.LastBlockHeight(), 2)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexP)

	testEnterNewHeight(tConsX)

	// This certificate has all signers' vote
	cert := tConsX.state.LastCertificate()

	assert.Empty(t, cert.Absences())
}

func TestConsensusBehindState(t *testing.T) {
	setup(t)

	// Consensus starts here
	testEnterNewHeight(tConsP)

	p := makeProposal(t, 1, 0)
	assert.NoError(t, tConsP.state.ValidateBlock(p.Block()))
	tConsP.SetProposal(p)

	// --------------------------------
	// Syncer commits a block
	commitBlockForAllStates(t)
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed by syncer before.
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))
}
