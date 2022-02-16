package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	setup(t)

	tConsensus.Round = 1
	h, _ := tConsensus.HeightRound()
	pid := util.RandomPeerID()
	pld := payload.NewHeartBeatPayload(h, 2, hash.GenerateTestHash())

	t.Run("Not in the committee, but processes hearbeat messages", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})

	testAddPeerToCommittee(t, tSync.signer.PublicKey(), tSync.SelfID())

	t.Run("Bob should query for votes", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})

	t.Run("Bob should not query for votes for previous round", func(t *testing.T) {
		pld := payload.NewHeartBeatPayload(h, 0, hash.GenerateTestHash())
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})

	t.Run("Bob should not query for votes for same round", func(t *testing.T) {
		pld := payload.NewHeartBeatPayload(h, 1, hash.GenerateTestHash())
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})
}

func TestBroadcastingHeartbeatMessages(t *testing.T) {
	setup(t)

	t.Run("It is not in committee", func(t *testing.T) {
		tSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHeartBeat)
		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeVote)
	})

	testAddPeerToCommittee(t, tSync.signer.PublicKey(), tSync.SelfID())

	t.Run("It is in committee", func(t *testing.T) {
		heightAlice, _ := tConsensus.HeightRound()
		v1, _ := vote.GenerateTestPrepareVote(heightAlice, 0)
		tConsensus.Votes = []*vote.Vote{v1}

		tSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHeartBeat)
		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeVote)
	})
}
