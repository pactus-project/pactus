package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingQueryProposalMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	pid := util.RandomPeerID()
	msg := message.NewQueryProposalMessage(consensusHeight, 0)
	tConsensus.SetProposal(prop)

	t.Run("Not in the committee, should not respond to the query proposal message", func(t *testing.T) {
		assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))
	})

	testAddPeerToCommittee(t, pid, nil)

	t.Run("In the committee, should respond to the query proposal message", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeProposal)
		assert.Equal(t, bdl.Message.(*message.ProposalMessage).Proposal.Hash(), prop.Hash())
	})

	t.Run("In the committee, but doesn't have the proposal", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 1)
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

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

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signer.PublicKey())

	t.Run("In the committee, should send query proposal message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryProposal)
	})

	t.Run("Proposal set before", func(t *testing.T) {
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 0)
		tConsensus.SetProposal(prop)
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryProposal)
	})

	t.Run("Proposal found inside the cache", func(t *testing.T) {
		prop, _ := proposal.GenerateTestProposal(consensusHeight, 1)
		tSync.cache.AddProposal(prop)

		msg := message.NewQueryProposalMessage(consensusHeight, 1)
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryProposal)
	})
}
