package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestQueryVotesType(t *testing.T) {
	msg := &QueryVoteMessage{}
	assert.Equal(t, TypeQueryVote, msg.Type())
}

func TestQueryVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid round", func(t *testing.T) {
		msg := NewQueryVoteMessage(0, -1, ts.RandValAddress())

		err := msg.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "invalid round"})
	})

	t.Run("OK", func(t *testing.T) {
		msg := NewQueryVoteMessage(100, 0, ts.RandValAddress())

		assert.NoError(t, msg.BasicCheck())
		assert.Contains(t, msg.LogString(), "100")
	})
}
