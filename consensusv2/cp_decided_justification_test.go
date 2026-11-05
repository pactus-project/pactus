package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

// review: can a single Byzantine validator (1 of 4) end the current round with a
// cp:DECIDED vote whose justification is not a quorum certificate?
func TestDecidedVoteWithInitYesJust(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	td.changeProposerTimeout(td.consX)
	td.checkHeightRound(t, td.consX, 1, 0)

	// Validator B alone signs a DECIDED/Yes vote justified by JustInitYes (no certificate).
	td.addCPDecidedVote(t, td.consX, hash.UndefHash, 1, 0,
		vote.CPValueYes, &vote.JustInitYes{}, tIndexB)

	h, r := td.consX.HeightRound()
	t.Logf("after the vote: height=%d round=%d", h, r)
	assert.Equal(t, types.Round(0), r, "a single unjustified DECIDED vote must not end the round")
}
