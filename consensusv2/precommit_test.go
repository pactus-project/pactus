package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(0)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	prop := td.makeProposal(t, h, r)
	propBlockHash := prop.Block().Hash()

	_, _, decidedJust := td.makeChangeProposerJusts(t, propBlockHash, h, r)

	decideVote := vote.NewCPDecidedVote(propBlockHash, h, r, 0, vote.CPValueNo, decidedJust, td.consX.valKey.Address())
	td.HelperSignVote(td.consX.valKey, decideVote)

	td.consP.AddVote(decideVote)
	assert.Equal(t, "precommit", td.consP.currentState.name())

	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexX)
	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexY)

	td.shouldPublishQueryProposal(t, td.consP, h, r)
}
