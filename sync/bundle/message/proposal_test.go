package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/proposal"
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
		prop := proposal.NewProposal(0, 0, nil)
		m := NewProposalMessage(prop)

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidBlock)
	})

	t.Run("OK", func(t *testing.T) {
		prop, _ := ts.GenerateTestProposal(100, 0)
		m := NewProposalMessage(prop)

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
