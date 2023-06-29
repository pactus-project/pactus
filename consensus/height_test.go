package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestNewHeightTimedout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	td.commitBlockForAllStates(t)

	s := &newHeightState{td.consX}

	// Invalid target
	s.onTimeout(&ticker{Height: 2, Target: 3})
	td.checkHeightRound(t, td.consX, 1, 0)

	s.onTimeout(&ticker{Height: 2, Target: tickerTargetNewHeight})
	td.checkHeightRound(t, td.consX, 2, 0)
}

func TestNewHeightEntry(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.consX.MoveToNewHeight()
	td.consX.MoveToNewHeight()

	td.checkHeightRoundWait(t, td.consX, 2, 0)
	assert.True(t, td.consX.active)
	assert.NotEqual(t, td.consX.currentState.name(), "new-height")
}
func TestUpdateCertificate(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	p := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(p)

	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexB)

	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexX)
	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexB)

	assert.Equal(t, td.consX.state.LastBlockHeight(), uint32(2))

	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)
	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexP)

	td.enterNewHeight(td.consX)

	// This certificate has all signers' vote
	cert := td.consX.state.LastCertificate()

	assert.Empty(t, cert.Absentees())
}

func TestConsensusHeightIsShorterThanState(t *testing.T) {
	td := setup(t)

	// Consensus starts here
	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, 1, 0)
	assert.NoError(t, td.consP.state.ValidateBlock(p.Block()))
	td.consP.SetProposal(p)

	// --------------------------------
	// Commit a block
	td.commitBlockForAllStates(t)
	// --------------------------------

	// Consensus tries to add more votes and commit the block which is committed before.
	td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)

	assert.Error(t, td.consP.state.ValidateBlock(p.Block()))
}
