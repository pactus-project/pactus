package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type proposalHandler struct {
	*synchronizer
}

func newProposalHandler(sync *synchronizer) messageHandler {
	return &proposalHandler{
		sync,
	}
}

func (handler *proposalHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.ProposalMessage)
	handler.logger.Trace("parsing Proposal message", "msg", msg)

	handler.cache.AddProposal(msg.Proposal)
	handler.consensus.SetProposal(msg.Proposal)

	return nil
}

func (handler *proposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
