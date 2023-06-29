package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	td := setup(t, nil)

	td.consMocks[0].Round = 1
	h, _ := td.consMgr.HeightRound()
	pid := td.RandomPeerID()
	msg := message.NewHeartBeatMessage(h, 2, td.RandomHash())

	t.Run("Not in the committee, but processes heartbeat messages", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeQueryVotes)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[0].PublicKey())

	t.Run("In the committee, should query for votes", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for previous round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 0, td.RandomHash())
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for same round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 1, td.RandomHash())
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeQueryVotes)
	})
}

func TestBroadcastingHeartbeatMessages(t *testing.T) {
	config := testConfig()
	config.HeartBeatTimer = 1 * time.Second
	td := setup(t, config)

	t.Run("It is not in committee", func(t *testing.T) {
		td.sync.broadcastHeartBeat()
		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeHeartBeat)
		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeVote)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[1].PublicKey())

	t.Run("It is in committee", func(t *testing.T) {
		heightAlice, _ := td.consMgr.HeightRound()
		v1, _ := td.GenerateTestPrepareVote(heightAlice, 0)
		td.consMgr.AddVote(v1)

		td.sync.broadcastHeartBeat()
		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeHeartBeat)
		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeVote)
	})
}
