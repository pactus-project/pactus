package sync

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryProposalMessages(t *testing.T) {
	setup(t)

	consensusHeight, _ := tConsMgr.HeightRound()
	prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	pid := network.TestRandomPeerID()
	tConsMgr.SetProposal(prop)

	t.Run("Not in the committee, should not respond to the query proposal message", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 0)

		assert.Error(t, testReceivingNewMessage(tSync, msg, pid))
	})

	testAddPeerToCommittee(t, pid, nil)

	t.Run("In the committee, but not the same height", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight+1, 0)
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeProposal)
	})
	t.Run("In the committee, should respond to the query proposal message", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 0)
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeProposal)
		assert.Equal(t, bdl.Message.(*message.ProposalMessage).Proposal.Hash(), prop.Hash())
	})

	t.Run("In the committee, but doesn't have the proposal", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 1)
		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeProposal)
	})
}

func TestBroadcastingQueryProposalMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	msg := message.NewQueryProposalMessage(consensusHeight, 0)

	t.Run("Not in the committee, should not send query proposal message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryProposal)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signers[0].PublicKey())

	t.Run("In the committee, should send query proposal message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryProposal)
	})
}
