package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryVoteMessages(t *testing.T) {
	td := setup(t, nil)

	consHeight, consRound := td.consMgr.HeightRound()
	t.Run("doesn't have any votes", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewQueryVoteMessage(consHeight, consRound, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("should respond to the query votes message", func(t *testing.T) {
		vote, _ := td.GenerateTestPrecommitVote(consHeight, consRound)
		td.consMgr.AddVote(vote)
		pid := td.RandPeerID()
		msg := message.NewQueryVoteMessage(consHeight, consRound, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeVote)
		assert.Equal(t, vote.Hash(), bdl.Message.(*message.VoteMessage).Vote.Hash())
	})
}

func TestBroadcastingQueryVoteMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight := td.state.LastBlockHeight() + 1
	msg := message.NewQueryVoteMessage(consensusHeight, 1, td.RandValAddress())
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeQueryVote)
}
