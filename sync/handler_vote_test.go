package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestHandlerVoteParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing vote message", func(t *testing.T) {
		v, _ := td.GenerateTestPrecommitVote(1, 0)
		msg := message.NewVoteMessage(v)
		pid := td.RandPeerID()

		td.receivingNewMessage(td.sync, msg, pid)
		assert.Contains(t, td.consMocks[0].AllVotes(), v)
	})
}
