package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	setup(t)

	t.Run("Alice is not in committee", func(t *testing.T) {
		tAliceSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
	})

	joinAliceToCommittee(t)

	t.Run("Alice is in committee", func(t *testing.T) {
		heightAlice, _ := tAliceConsensus.HeightRound()
		v1, _ := vote.GenerateTestPrepareVote(heightAlice, 0)
		tAliceConsensus.Votes = []*vote.Vote{v1}

		tAliceSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
	})

	t.Run("Bob processes Alice's HeartBeat but he is not in committee", func(t *testing.T) {
		h, r := tBobConsensus.HeightRound()
		pld := payload.NewHeartBeatPayload(h, r+2, hash.GenerateTestHash())
		simulatingReceiveingNewMessage(t, tBobSync, pld, util.RandomPeerID())

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryVotes)
	})

	joinBobToCommittee(t)

	t.Run("Bob should query for votes", func(t *testing.T) {
		h, r := tBobConsensus.HeightRound()
		pld := payload.NewHeartBeatPayload(h, r+2, hash.GenerateTestHash())
		simulatingReceiveingNewMessage(t, tBobSync, pld, util.RandomPeerID())

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryVotes)
	})

	t.Run("Bob should not query for votes for previous round", func(t *testing.T) {
		h, r := tBobConsensus.HeightRound()
		pld := payload.NewHeartBeatPayload(h, r+1, hash.GenerateTestHash())
		simulatingReceiveingNewMessage(t, tBobSync, pld, util.RandomPeerID())

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryVotes)
	})
}
