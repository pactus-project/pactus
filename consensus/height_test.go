package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/vote"
)

func TestNewHeightTimedout(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)
	commitBlockForAllStates(t)

	s := &newHeightState{tConsX}

	// Invalid target
	s.timedout(&ticker{Height: 2, Target: 3})
	checkHeightRound(t, tConsX, 1, 0)

	s.timedout(&ticker{Height: 2, Target: tickerTargetNewHeight})
	checkHeightRound(t, tConsX, 2, 0)
}

func TestNewHeightDuplicateEntry(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)
	testEnterNewRound(tConsX)

	s := &newHeightState{tConsX}

	s.timedout(&ticker{Height: 1, Target: tickerTargetNewHeight})
	checkHeightRound(t, tConsX, 1, 1)
}

func TestUpdateCertificate(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)

	commitBlockForAllStates(t)

	s := &newHeightState{tConsX}

	h := tConsX.state.LastBlockHash()
	cert1 := tConsX.state.LastCertificate()
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, h, tIndexX)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, h, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, h, tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, h, tIndexP)

	s.execute()

	// This certificate has all signers' vote
	cert2 := tConsX.state.LastCertificate()

	assert.NotEqual(t, cert1.Hash(), cert2.Hash())
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
