package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVoteType(t *testing.T) {
	msg := &VoteMessage{}
	assert.Equal(t, TypeVote, msg.Type())
}

func TestVoteMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No vote", func(t *testing.T) {
		msg := NewVoteMessage(nil)

		require.ErrorIs(t, msg.BasicCheck(), BasicCheckError{Reason: "no vote"})
	})

	t.Run("Invalid vote", func(t *testing.T) {
		vte := vote.NewPrepareVote(ts.RandHash(), ts.RandHeight(), -1, ts.RandValAddress())
		msg := NewVoteMessage(vte)

		require.ErrorIs(t, msg.BasicCheck(), vote.BasicCheckError{Reason: "invalid round"})
	})

	t.Run("OK", func(t *testing.T) {
		vte, _ := ts.GenerateTestPrepareVote(100, 0)
		msg := NewVoteMessage(vte)

		require.NoError(t, msg.BasicCheck())
		assert.Contains(t, msg.LogString(), vte.LogString())
	})
}
