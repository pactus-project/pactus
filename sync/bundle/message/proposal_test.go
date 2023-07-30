package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestProposalType(t *testing.T) {
	m := &ProposalMessage{}
	assert.Equal(t, m.Type(), TypeProposal)
}

func TestProposalMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid proposal", func(t *testing.T) {
		proposal, _ := ts.GenerateTestProposal(100, -1)
		m := NewProposalMessage(proposal)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		proposal, _ := ts.GenerateTestProposal(100, 0)
		m := NewProposalMessage(proposal)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.String(), "100")
	})
}
