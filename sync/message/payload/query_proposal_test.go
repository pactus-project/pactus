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
		p1 := NewQueryProposalPayload(-1, 0)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p1 := NewQueryProposalPayload(0, -1)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p2 := NewQueryProposalPayload(100, 0)
		assert.NoError(t, p2.SanityCheck())
	})
}
