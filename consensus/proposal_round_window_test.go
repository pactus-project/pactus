package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A proposal for a round too far ahead is not kept: the log would keep every round it receives.
func TestProposalForFarRound(t *testing.T) {
	td := setup(t)
	td.enterNewHeight(td.consP)

	round := td.consP.round + maxFutureRounds + 1
	prop := td.makeProposal(t, td.consP.height, round)
	td.consP.SetProposal(prop)

	td.consP.lk.Lock()
	defer td.consP.lk.Unlock()
	assert.Nil(t, td.consP.log.RoundProposal(round))
}
