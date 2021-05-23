package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	p := &QueryVotesPayload{}
	assert.Equal(t, p.Type(), PayloadTypeQueryVotes)
}

func TestQueryVotesPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		p := NewQueryVotesPayload(-1, 0)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p := NewQueryVotesPayload(0, -1)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewQueryVotesPayload(100, 0)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
