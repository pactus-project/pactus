package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryVotesMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight, _ := td.consMgr.HeightRound()
	v1, _ := td.GenerateTestPrecommitVote(consensusHeight, 0)
	td.consMgr.AddVote(v1)
	pid := td.RandPeerID()

	t.Run("doesn't have active validator", func(t *testing.T) {
		msg := message.NewQueryVotesMessage(consensusHeight, 1, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		td.shouldNotPublishMessageWithThisType(t, message.TypeVote)
	})

	td.consMocks[0].Active = true

	t.Run("should respond to the query votes message", func(t *testing.T) {
		msg := message.NewQueryVotesMessage(consensusHeight, 1, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeVote)
		assert.Equal(t, v1.Hash(), bdl.Message.(*message.VoteMessage).Vote.Hash())
	})

	t.Run("doesn't have any votes", func(t *testing.T) {
		msg := message.NewQueryVotesMessage(consensusHeight+1, 1, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		td.shouldNotPublishMessageWithThisType(t, message.TypeVote)
	})
}

func TestBroadcastingQueryVotesMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight := td.state.LastBlockHeight() + 1
	msg := message.NewQueryVotesMessage(consensusHeight, 1, td.RandValAddress())
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeQueryVote)
}
