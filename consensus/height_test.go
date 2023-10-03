package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeightTimeout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consY)
	td.commitBlockForAllStates(t)

	s := &newHeightState{td.consY}
	s.enter()

	// Invalid target
	s.onTimeout(&ticker{Height: 2, Target: -1})
	td.checkHeightRound(t, td.consY, 2, 0)

	s.onTimeout(&ticker{Height: 2, Target: tickerTargetNewHeight})
	td.checkHeightRound(t, td.consY, 2, 0)
	td.shouldPublishProposal(t, td.consY, 2, 0)
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
