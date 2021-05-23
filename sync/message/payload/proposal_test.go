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
	t.Run("Invalid proposal", func(t *testing.T) {
		proposal, _ := proposal.GenerateTestProposal(100, -1)
		p := NewProposalPayload(proposal)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		proposal, _ := proposal.GenerateTestProposal(100, 0)
		p := NewProposalPayload(proposal)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
