package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeightTimeout(t *testing.T) {
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
	td.newHeightTimeout(td.consX)

	// double entry and timeout
	td.consX.MoveToNewHeight()
	td.newHeightTimeout(td.consX)

	td.checkHeightRound(t, td.consX, 2, 0)
	assert.True(t, td.consX.active)
	assert.NotEqual(t, td.consX.currentState.name(), "new-height")
}

func TestUpdateCertificate(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	p := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(p)

	td.addPrepareVote(td.consX, p.Block().Hash(), 2, 0, tIndexX)
	td.addPrepareVote(td.consX, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consX, p.Block().Hash(), 2, 0, tIndexB)

	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexX)
	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexB)

	cert1 := td.consX.state.LastCertificate()
	assert.Contains(t, cert1.Committers(), int32(tIndexP))
	assert.Contains(t, cert1.Absentees(), int32(tIndexP))

	td.addPrecommitVote(td.consX, p.Block().Hash(), 2, 0, tIndexP)

	td.newHeightTimeout(td.consX)

	// This certificate has all signers' vote
	cert2 := td.consX.state.LastCertificate()
	assert.Empty(t, cert2.Absentees())
}
