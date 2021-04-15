package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func TestParsingQueryProposalMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	consensusHeight := tAliceState.LastBlockHeight() + 1
	p1, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	p2, _ := proposal.GenerateTestProposal(consensusHeight, 1)

	tAliceSync.cache.AddProposal(p1)
	tBobConsensus.SetProposal(p2)
	pld := payload.NewQueryProposalPayload(consensusHeight, 1)

	t.Run("Alice should not send query proposal message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice queries for a proposal, but she has it in her cache", func(t *testing.T) {
		tAliceBroadcastCh <- payload.NewQueryProposalPayload(consensusHeight, 0)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob should not process alice's message because she is not an active validator", func(t *testing.T) {
		tBobNet.ReceivingMessageFromOtherPeer(tAlicePeerID, pld)
		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})

	joinAliceToTheSet(t)

	t.Run("Alice should be able to send query proposal message because sh is an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob processes Alice's message", func(t *testing.T) {
		tBobNet.ReceivingMessageFromOtherPeer(tAlicePeerID, pld)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})
}
