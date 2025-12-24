package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestProposalType(t *testing.T) {
	msg := &ProposalMessage{}
	assert.Equal(t, TypeProposal, msg.Type())
}

func TestProposalMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("OK", func(t *testing.T) {
		prop := ts.GenerateTestProposal(100, 0)
		msg := NewProposalMessage(prop)

		assert.NoError(t, msg.BasicCheck())
		assert.Contains(t, msg.LogString(), "100")
	})
}
