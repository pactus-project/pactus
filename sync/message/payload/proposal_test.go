package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/proposal"
)

func TestProposalType(t *testing.T) {
	p := &ProposalPayload{}
	assert.Equal(t, p.Type(), PayloadTypeProposal)
}

func TestProposalPayload(t *testing.T) {
	proposal1, _ := proposal.GenerateTestProposal(100, -1)
	p1 := NewProposalPayload(proposal1)
	assert.Error(t, p1.SanityCheck())

	proposal2, _ := proposal.GenerateTestProposal(100, 0)
	p2 := NewProposalPayload(proposal2)
	assert.NoError(t, p2.SanityCheck())
}
