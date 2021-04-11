package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

func TestParsingHeartbeatMessages(t *testing.T) {
	setup(t)

	t.Run("Alice is not in committee", func(t *testing.T) {
		tAliceSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
	})

	joinAliceToTheSet(t)

	t.Run("Alice is in committee", func(t *testing.T) {
		aliceH, _ := tAliceConsensus.HeightRound()
		v1, _ := vote.GenerateTestPrepareVote(aliceH, 0)
		tAliceConsensus.Votes = []*vote.Vote{v1}

		tAliceSync.broadcastHeartBeat()
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
	})
}
