package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingVoteMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing vote message", func(t *testing.T) {
		v, _ := td.GenerateTestPrecommitVote(1, 0)
		msg := message.NewVoteMessage(v)

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, td.RandomPeerID()))
		assert.Equal(t, td.consMgr.PickRandomVote().Hash(), v.Hash())
	})
}
