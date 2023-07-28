package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	m := &QueryVotesMessage{}
	assert.Equal(t, m.Type(), TypeQueryVotes)
}

func TestQueryVotesMessage(t *testing.T) {
	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryVotesMessage(0, -1)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryVotesMessage(100, 0)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
