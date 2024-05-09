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

func (handler *proposalHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.ProposalMessage)
	handler.logger.Trace("parsing Proposal message", "msg", msg)

	handler.consMgr.SetProposal(msg.Proposal)

	return nil
}

func (handler *proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
