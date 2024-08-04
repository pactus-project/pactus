package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	m := &QueryVoteMessage{}
	assert.Equal(t, TypeQueryVote, m.Type())
}

func TestQueryVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryVoteMessage(0, -1, ts.RandValAddress())

		assert.Equal(t, errors.ErrInvalidRound, errors.Code(m.BasicCheck()))
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryVoteMessage(100, 0, ts.RandValAddress())

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
