package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestHandlerQueryVoteParsingMessages(t *testing.T) {
	td := setup(t, nil)

	consHeight, consRound := td.sync.getConsMgr().HeightRound()
	t.Run("doesn't have any votes", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewQueryVoteMessage(consHeight, consRound, td.RandValAddress())

		td.consV1Mgr.EXPECT().HandleQueryVote(consHeight, consRound).Return(nil).Times(1)

		td.receivingNewMessage(td.sync, msg, pid)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("should respond to the query votes message", func(t *testing.T) {
		vote, _ := td.GenerateTestPrecommitVote(consHeight, consRound)
		pid := td.RandPeerID()
		msg := message.NewQueryVoteMessage(consHeight, consRound, td.RandValAddress())

		td.consV1Mgr.EXPECT().HandleQueryVote(consHeight, consRound).Return(vote).Times(1)
		td.receivingNewMessage(td.sync, msg, pid)

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeVote)
		assert.Equal(t, vote.Hash(), bdl.Message.(*message.VoteMessage).Vote.Hash())
	})
}
