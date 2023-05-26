package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	setup(t)

	tConsMocks[0].Round = 1
	h, _ := tConsMgr.HeightRound()
	pid := network.TestRandomPeerID()
	msg := message.NewHeartBeatMessage(h, 2, hash.GenerateTestHash())

	t.Run("Not in the committee, but processes heartbeat messages", func(t *testing.T) {
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signers[0].PublicKey())

	t.Run("In the committee, should query for votes", func(t *testing.T) {
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for previous round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 0, hash.GenerateTestHash())
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for same round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 1, hash.GenerateTestHash())
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})
}

func TestBroadcastingHeartbeatMessages(t *testing.T) {
	tConfig.HeartBeatTimer = 1 * time.Second
	setup(t)

	t.Run("It is not in committee", func(t *testing.T) {
		tSync.broadcastHeartBeat()
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHeartBeat)
		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signers[1].PublicKey())

	t.Run("It is in committee", func(t *testing.T) {
		heightAlice, _ := tConsMgr.HeightRound()
		v1, _ := vote.GenerateTestPrepareVote(heightAlice, 0)
		tConsMgr.AddVote(v1)

		tSync.broadcastHeartBeat()
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHeartBeat)
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
	})
}
