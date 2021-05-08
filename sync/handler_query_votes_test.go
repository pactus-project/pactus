package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func TestParsingQueryVotesMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	consensusHeight := tAliceState.LastBlockHeight() + 1
	v1, _ := vote.GenerateTestPrecommitVote(consensusHeight, 0)

	tBobConsensus.AddVote(v1)
	pld := payload.NewQueryVotesPayload(consensusHeight, 1)

	t.Run("Alice should not send query votes message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryVotes)
	})

	t.Run("Bob should not process alice's message because she is not an active validator", func(t *testing.T) {
		tBobNet.ReceivingMessageFromOtherPeer(tAlicePeerID, pld)
		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeVote)
	})

	joinAliceToCommittee(t)

	t.Run("Alice should be able to send query votes message because she is an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryVotes)
	})

	t.Run("Bob processes Alice's message", func(t *testing.T) {
		tBobNet.ReceivingMessageFromOtherPeer(tAlicePeerID, pld)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeVote)
	})
}
