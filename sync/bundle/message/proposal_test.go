package message

import (
	"testing"

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

	t.Run("No proposal", func(t *testing.T) {
		msg := NewProposalMessage(nil)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{Reason: "no proposal"})
		assert.Zero(t, msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "{nil-proposal}")
	})

	t.Run("OK", func(t *testing.T) {
		prop := ts.GenerateTestProposal(ts.RandHeight(), ts.RandRound())
		msg := NewProposalMessage(prop)

		require.NoError(t, msg.BasicCheck())
		assert.Equal(t, prop.Height(), msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), prop.Block().LogString())
	})
}
