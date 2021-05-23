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
		p1 := NewQueryVotesPayload(-1, 0)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p1 := NewQueryVotesPayload(0, -1)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p2 := NewQueryVotesPayload(100, 0)
		assert.NoError(t, p2.SanityCheck())
	})
}
