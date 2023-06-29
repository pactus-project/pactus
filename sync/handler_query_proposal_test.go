package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingQueryProposalMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight, _ := td.consMgr.HeightRound()
	prop, _ := td.GenerateTestProposal(consensusHeight, 0)
	pid := td.RandomPeerID()
	td.consMgr.SetProposal(prop)

	t.Run("Not in the committee, should not respond to the query proposal message", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 0)

		assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))
	})

	td.addPeerToCommittee(t, pid, nil)

	t.Run("In the committee, but not the same height", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight+1, 0)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeProposal)
	})
	t.Run("In the committee, should respond to the query proposal message", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 0)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeProposal)
		assert.Equal(t, bdl.Message.(*message.ProposalMessage).Proposal.Hash(), prop.Hash())
	})

	t.Run("In the committee, but doesn't have the proposal", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, 1)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeProposal)
	})
}

func TestBroadcastingQueryProposalMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight := td.state.LastBlockHeight() + 1
	msg := message.NewQueryProposalMessage(consensusHeight, 0)

	t.Run("Not in the committee, should not send query proposal message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, td.network, message.MessageTypeQueryProposal)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[0].PublicKey())

	t.Run("In the committee, should send query proposal message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeQueryProposal)
	})
}
