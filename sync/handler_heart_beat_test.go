package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/vote"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	setup(t)

	tConsensus.Round = 1
	h, _ := tConsensus.HeightRound()
	pid := network.TestRandomPeerID()
	msg := message.NewHeartBeatMessage(h, 2, hash.GenerateTestHash())

	t.Run("Not in the committee, but processes hearbeat messages", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signer.PublicKey())

	t.Run("In the committee, should query for votes", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for previous round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 0, hash.GenerateTestHash())
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})

	t.Run("Should not query for votes for same round", func(t *testing.T) {
		msg := message.NewHeartBeatMessage(h, 1, hash.GenerateTestHash())
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryVotes)
	})
}

func TestBroadcastingHeartbeatMessages(t *testing.T) {
	setup(t)

	t.Run("It is not in committee", func(t *testing.T) {
		tSync.broadcastHeartBeat()
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHeartBeat)
		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signer.PublicKey())

	t.Run("It is in committee", func(t *testing.T) {
		heightAlice, _ := tConsensus.HeightRound()
		v1, _ := vote.GenerateTestPrepareVote(heightAlice, 0)
		tConsensus.Votes = []*vote.Vote{v1}

		tSync.broadcastHeartBeat()
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHeartBeat)
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeVote)
	})
}
