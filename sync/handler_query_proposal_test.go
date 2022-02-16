package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingQueryProposalMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	pid := util.RandomPeerID()
	pld := payload.NewQueryProposalPayload(consensusHeight, 0)
	tConsensus.SetProposal(prop)

	t.Run("Not in the committee, should not respond to the query proposal message", func(t *testing.T) {
		assert.Error(t, testReceiveingNewMessage(t, tSync, pld, pid))
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeerToCommittee(t, pub, pid)

	t.Run("In the committee, should respond to the query proposal message", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeProposal)
		assert.Equal(t, msg.Payload.(*payload.ProposalPayload).Proposal.Hash(), prop.Hash())
	})

	t.Run("In the committee, but doesn't have the proposal", func(t *testing.T) {
		pld := payload.NewQueryProposalPayload(consensusHeight, 1)
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeProposal)
	})
}

func TestBroadcastingQueryProposalMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	pld := payload.NewQueryProposalPayload(consensusHeight, 0)

	t.Run("Not in the committee, should not send query proposal message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryProposal)
	})

	testAddPeerToCommittee(t, tSync.signer.PublicKey(), tSync.SelfID())

	t.Run("In the committee, should send query proposal message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryProposal)
	})

	t.Run("Proposal set before", func(t *testing.T) {
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
		tConsensus.SetProposal(prop)
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryProposal)
	})

	t.Run("Proposal found inside the cache", func(t *testing.T) {
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 1)
		tSync.cache.AddProposal(prop)

		pld := payload.NewQueryProposalPayload(consensusHeight, 1)
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryProposal)
	})
}
