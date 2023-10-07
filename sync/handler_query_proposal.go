package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type queryProposalHandler struct {
	*synchronizer
}

func newQueryProposalHandler(sync *synchronizer) messageHandler {
	return &queryProposalHandler{
		sync,
	}
}

func (handler *queryProposalHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.QueryProposalMessage)
	handler.logger.Trace("parsing QueryProposal message", "message", msg, "initiator", initiator)

	height, _ := handler.consMgr.HeightRound()
	if msg.Height == height {
		// TODO: this should be refactored
		// if !handler.peerIsInTheCommittee(initiator) {
		// 	return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the committee")
		// }
		prop := handler.consMgr.Proposal()
		if prop != nil {
			response := message.NewProposalMessage(prop)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending QueryProposal ignored. We are not in the committee")
		return nil
	}
	bdl := bundle.NewBundle(handler.SelfID(), m)

	return bdl
}
