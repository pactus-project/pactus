package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryProposalType(t *testing.T) {
	p := &QueryProposalPayload{}
	assert.Equal(t, p.Type(), PayloadTypeQueryProposal)
}

func TestQueryProposalPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		p := NewQueryProposalPayload(-1, 0)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p := NewQueryProposalPayload(0, -1)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewQueryProposalPayload(100, 0)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
