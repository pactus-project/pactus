package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/vote"
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

	tConsX.MoveToNewHeight()
	tConsX.MoveToNewHeight()

	checkHeightRoundWait(t, tConsX, 1, 0)
}

func TestUpdateCertificate(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	p := makeProposal(t, 2, 0)
	tConsX.SetProposal(p)

	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexB)

	assert.Equal(t, tConsX.state.LastBlockHeight(), int32(2))

	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)
	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexP)

	testEnterNewHeight(tConsX)

	// This certificate has all signers' vote
	cert := tConsX.state.LastCertificate()

	assert.Empty(t, cert.Absentees())
}

func TestConsensusHeightIsShorterThanState(t *testing.T) {
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
	testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)

	assert.Error(t, tConsP.state.ValidateBlock(p.Block()))
}
