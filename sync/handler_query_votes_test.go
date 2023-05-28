package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryVotesMessages(t *testing.T) {
	setup(t)

	consensusHeight, _ := tConsMgr.HeightRound()
	v1, _ := vote.GenerateTestPrecommitVote(consensusHeight, 0)
	tConsMgr.AddVote(v1)
	pid := network.TestRandomPeerID()
	msg := message.NewQueryVotesMessage(consensusHeight, 1)

	t.Run("Not known peer, should not respond to the query vote message", func(t *testing.T) {
		assert.Error(t, testReceivingNewMessage(tSync, msg, pid))
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid, false)

	t.Run("Not in the committee, should not respond to the query vote message", func(t *testing.T) {
		assert.Error(t, testReceivingNewMessage(tSync, msg, pid))
	})

	testAddPeerToCommittee(t, pid, nil)

	t.Run("In the committee, should respond to the query vote message", func(t *testing.T) {
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
		assert.Equal(t, bdl.Message.(*message.VoteMessage).Vote.Hash(), v1.Hash())
	})

	t.Run("In the committee, but doesn't have the vote", func(t *testing.T) {
		msg := message.NewQueryVotesMessage(consensusHeight+1, 1)
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
	})
}

func TestBroadcastingQueryVotesMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	msg := message.NewQueryVotesMessage(consensusHeight, 1)

	t.Run("Not in the committee, should not send query vote message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signers[0].PublicKey())
	t.Run("In the committee, should send query vote message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})
}
