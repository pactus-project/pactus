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
		pid := td.RandPeerID()

		td.receivingNewMessage(td.sync, msg, pid)
		assert.Equal(t, v.Hash(), td.consMgr.PickRandomVote(0).Hash())
	})
}
