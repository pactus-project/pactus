package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRemoveInvalidProposal(t *testing.T) {
	cons, pvals = newTestConsensus(VAL_2)

	cons.enterNewHeight(1)

	addr := pvals[VAL_1].Address()
	block, _ := block.GenerateTestBlock(&addr)
	invalidProposal := vote.NewProposal(1, 0, block)
	pvals[VAL_1].SignMsg(invalidProposal)
	cons.setProposal(invalidProposal)
	assert.Nil(t, cons.votes.RoundProposal(0))
}
