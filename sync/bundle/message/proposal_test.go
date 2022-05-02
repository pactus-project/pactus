package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestProposalType(t *testing.T) {
	m := &ProposalMessage{}
	assert.Equal(t, m.Type(), MessageTypeProposal)
}

func TestProposalMessage(t *testing.T) {
	t.Run("Invalid proposal", func(t *testing.T) {
		proposal, _ := proposal.GenerateTestProposal(100, -1)
		m := NewProposalMessage(proposal)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		proposal, _ := proposal.GenerateTestProposal(100, 0)
		m := NewProposalMessage(proposal)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
