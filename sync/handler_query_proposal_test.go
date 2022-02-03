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
	proposalRound0, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	proposalRound1, _ := proposal.GenerateTestProposal(consensusHeight, 1)

	pldRound0 := payload.NewQueryProposalPayload(consensusHeight, 0)
	pldRound1 := payload.NewQueryProposalPayload(consensusHeight, 1)

	tBobConsensus.Round = 0
	tBobConsensus.SetProposal(proposalRound0)

	t.Run("Alice should not send query proposal message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pldRound0
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice queries for a proposal, but she has it in her cache", func(t *testing.T) {
		tAliceBroadcastCh <- payload.NewQueryProposalPayload(consensusHeight, 0)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob should not process Alice's message because she is not an active validator", func(t *testing.T) {
		simulatingReceiveingNewMessage(t, tBobSync, pldRound0, tAlicePeerID)
		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})

	joinAliceToCommittee(t)
	tBobConsensus.Round = 1
	tBobConsensus.SetProposal(proposalRound1)

	t.Run("Alice should not query for proposal, because she has proposal in her cache", func(t *testing.T) {
		tAliceSync.cache.AddProposal(proposalRound0)

		tAliceBroadcastCh <- pldRound0
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob should not send Alice the stale proposal", func(t *testing.T) {
		// This case is importance for stability and performance of consensus

		simulatingReceiveingNewMessage(t, tBobSync, pldRound0, tAlicePeerID)
		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})

	t.Run("Alice should be able to send query proposal message because she is an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pldRound1
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob processes Alice's message", func(t *testing.T) {
		simulatingReceiveingNewMessage(t, tBobSync, pldRound1, tAlicePeerID)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})
}
