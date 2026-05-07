package message

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		require.NoError(t, msg.BasicCheck())
		assert.Equal(t, types.Height(100), msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "100")
	})
}
