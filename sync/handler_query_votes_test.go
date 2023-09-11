package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryVotesMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight, _ := td.consMgr.HeightRound()
	v1, _ := td.GenerateTestPrecommitVote(consensusHeight, 0)
	td.consMgr.AddVote(v1)
	pid := td.RandPeerID()
	msg := message.NewQueryVotesMessage(consensusHeight, 1)

	t.Run("Not known peer, should not respond to the query vote message", func(t *testing.T) {
		assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))
	})

	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, services.New(services.None))

	t.Run("Not in the committee, should not respond to the query vote message", func(t *testing.T) {
		assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))
	})

	td.addPeerToCommittee(t, pid, nil)

	t.Run("In the committee, should respond to the query vote message", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeVote)
		assert.Equal(t, bdl.Message.(*message.VoteMessage).Vote.Hash(), v1.Hash())
	})

	t.Run("In the committee, but doesn't have the vote", func(t *testing.T) {
		msg := message.NewQueryVotesMessage(consensusHeight+1, 1)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeVote)
	})
}

func TestBroadcastingQueryVotesMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight := td.state.LastBlockHeight() + 1
	msg := message.NewQueryVotesMessage(consensusHeight, 1)

	t.Run("Not in the committee, should not send query vote message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeQueryVotes)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[0].PublicKey())
	t.Run("In the committee, should send query vote message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldPublishMessageWithThisType(t, td.network, message.TypeQueryVotes)
	})
}
