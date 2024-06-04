package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	m := &QueryVotesMessage{}
	assert.Equal(t, m.Type(), TypeQueryVote)
}

func TestQueryVotesMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryVotesMessage(0, -1, ts.RandValAddress())

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryVotesMessage(100, 0, ts.RandValAddress())

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
