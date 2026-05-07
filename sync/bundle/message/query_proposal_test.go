package message

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryProposalType(t *testing.T) {
	msg := &QueryProposalMessage{}
	assert.Equal(t, TypeQueryProposal, msg.Type())
}

func TestQueryProposalMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid round", func(t *testing.T) {
		msg := NewQueryProposalMessage(0, -1, ts.RandValAddress())

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{"invalid round"})
	})

	t.Run("OK", func(t *testing.T) {
		msg := NewQueryProposalMessage(100, 0, ts.RandValAddress())

		require.NoError(t, msg.BasicCheck())
		assert.Equal(t, types.Height(100), msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), "100")
	})
}
