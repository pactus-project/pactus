package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRemoveInvalidProposal(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)

	addr := signers[VAL1].Address()
	b, _ := block.GenerateTestBlock(&addr)
	invalidProposal := vote.NewProposal(1, 0, *b)
	signers[VAL1].SignMsg(invalidProposal)
	cons.setProposal(invalidProposal)
	assert.Nil(t, cons.votes.RoundProposal(0))
}
