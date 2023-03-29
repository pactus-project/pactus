package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type proposalHandler struct {
	*synchronizer
}

func newProposalHandler(sync *synchronizer) messageHandler {
	return &proposalHandler{
		sync,
	}
}

func (handler *proposalHandler) ParsMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.ProposalMessage)
	handler.logger.Trace("parsing Proposal message", "message", msg)

	handler.consensus.SetProposal(msg.Proposal)

	return nil
}

func (handler *proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
