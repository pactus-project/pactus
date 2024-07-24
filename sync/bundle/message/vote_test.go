package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestVoteType(t *testing.T) {
	m := &VoteMessage{}
	assert.Equal(t, TypeVote, m.Type())
}

func TestVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid vote", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), ts.RandHeight(), -1, ts.RandValAddress())
		m := NewVoteMessage(v)

		assert.ErrorIs(t, m.BasicCheck(), vote.BasicCheckError{Reason: "invalid round"})
	})

	t.Run("OK", func(t *testing.T) {
		v, _ := ts.GenerateTestPrepareVote(100, 0)
		m := NewVoteMessage(v)

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), v.String())
	})
}
