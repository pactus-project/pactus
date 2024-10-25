package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestNewHeightTimeout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consY)
	td.commitBlockForAllStates(t)

	consState := &newHeightState{td.consY}
	consState.enter()

	// Invalid target
	consState.onTimeout(&ticker{Height: 2, Target: -1})
	td.checkHeightRound(t, td.consY, 2, 0)

	consState.onTimeout(&ticker{Height: 2, Target: tickerTargetNewHeight})
	td.checkHeightRound(t, td.consY, 2, 0)
	td.shouldPublishProposal(t, td.consY, 2, 0)
}

func TestNewHeightDoubleEntry(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.consX.MoveToNewHeight()
	td.newHeightTimeout(td.consX)

	// double entry and timeout
	td.consX.MoveToNewHeight()

	td.checkHeightRound(t, td.consX, 2, 0)
	assert.True(t, td.consX.active)
	assert.NotEqual(t, "new-height", td.consX.currentState.name())
}

func TestNewHeightTimeBehindNetwork(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.consP.MoveToNewHeight()

	height := uint32(2)
	round := int16(0)
	prop := td.makeProposal(t, height, round)

	td.consP.SetProposal(prop)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, prop.Block().Hash())
	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
}
