package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
)

func TestHandlerVoteParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing vote message", func(t *testing.T) {
		vte, _ := td.GenerateTestPrecommitVote(1, 0)
		msg := message.NewVoteMessage(vte)
		pid := td.RandPeerID()

		td.consV1Mgr.EXPECT().AddVote(vte).Return().Times(1)

		td.receivingNewMessage(td.sync, msg, pid)
	})
}
