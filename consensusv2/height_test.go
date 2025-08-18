package consensusv2

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
	assert.NotEqual(t, td.consX.currentState.name(), "new-height")
}

func TestNewHeightTimeBehindNetwork(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.consP.MoveToNewHeight()

	h := uint32(2)
	r := int16(0)
	p := td.makeProposal(t, h, r)
	blockHash := p.Block().Hash()

	td.consP.SetProposal(p)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrecommitVote(td.consP, blockHash, h, r, tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, blockHash)
	td.shouldPublishBlockAnnounce(t, td.consP, blockHash)
}
