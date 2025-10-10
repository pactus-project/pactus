package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestHandlerQueryProposalParsingMessages(t *testing.T) {
	td := setup(t, nil)

	consHeight, consRound := td.sync.getConsMgr().HeightRound()

	t.Run("doesn't have the proposal", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewQueryProposalMessage(consHeight, consRound, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("should respond to the query proposal message", func(t *testing.T) {
		prop := td.GenerateTestProposal(consHeight, 0)
		pid := td.RandPeerID()
		td.sync.getConsMgr().SetProposal(prop)
		msg := message.NewQueryProposalMessage(consHeight, consRound, td.RandValAddress())
		td.receivingNewMessage(td.sync, msg, pid)

		bdl := td.shouldPublishMessageWithThisType(t, message.TypeProposal)
		assert.Equal(t, prop.Hash(), bdl.Message.(*message.ProposalMessage).Proposal.Hash())
	})
}

func TestHandlerQueryProposalBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	consHeight, consRound := td.sync.getConsMgr().HeightRound()
	msg := message.NewQueryProposalMessage(consHeight, consRound, td.RandValAddress())
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeQueryProposal)
}
