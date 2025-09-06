package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/types/protocol"
)

type proposalHandler struct {
	*synchronizer
}

func newProposalHandler(sync *synchronizer) messageHandler {
	return &proposalHandler{
		sync,
	}
}

func (handler *proposalHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.ProposalMessage)
	handler.logger.Trace("parsing Proposal message", "msg", msg)

	handler.consMgr.SetProposal(msg.Proposal)

	// TODO: This condition can be removed in future releases.
	// This helps to support old nodes that don't specify their protocol version in proposal.
	if msg.ProtocolVersion != protocol.ProtocolVersionUnknown {
		handler.state.UpdateValidatorProtocolVersion(
			msg.Proposal.Block().Header().ProposerAddress(),
			msg.ProtocolVersion,
		)
	}
}

func (*proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
