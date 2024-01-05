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
	pid := td.RandPeerID()
	td.consMgr.SetProposal(prop)

	t.Run("not the same height", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight+1, td.RandValAddress())
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, message.TypeProposal)
	})
	t.Run("should respond to the query proposal message", func(t *testing.T) {
		msg := message.NewQueryProposalMessage(consensusHeight, td.RandValAddress())
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeProposal)
		assert.Equal(t, bdl.Message.(*message.ProposalMessage).Proposal.Hash(), prop.Hash())
	})

	t.Run("doesn't have the proposal", func(t *testing.T) {
		td.consMocks[0].CurProposal = nil
		msg := message.NewQueryProposalMessage(consensusHeight, td.RandValAddress())
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, message.TypeProposal)
	})
}

func TestBroadcastingQueryProposalMessages(t *testing.T) {
	td := setup(t, nil)

	consensusHeight := td.state.LastBlockHeight() + 1
	msg := message.NewQueryProposalMessage(consensusHeight, td.RandValAddress())
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeQueryProposal)
}
