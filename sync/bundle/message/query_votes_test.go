package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	m := &QueryVotesMessage{}
	assert.Equal(t, m.Type(), MessageTypeQueryVotes)
}

func TestQueryVotesMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewQueryVotesMessage(-1, 0)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryVotesMessage(0, -1)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryVotesMessage(100, 0)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
