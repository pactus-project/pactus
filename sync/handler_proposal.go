package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
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

	handler.getConsMgr().SetProposal(msg.Proposal)

	handler.state.UpdateValidatorProtocolVersion(
		msg.Proposal.Block().Header().ProposerAddress(),
		msg.ProtocolVersion,
	)
}

func (*proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
