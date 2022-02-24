package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryProposalType(t *testing.T) {
	m := &QueryProposalMessage{}
	assert.Equal(t, m.Type(), MessageTypeQueryProposal)
}

func TestQueryProposalMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewQueryProposalMessage(-1, 0)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryProposalMessage(0, -1)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryProposalMessage(100, 0)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
