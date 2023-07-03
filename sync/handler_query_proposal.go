package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util/errors"
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
	handler.logger.Trace("parsing QueryProposal message", "message", msg)

	height, _ := handler.consMgr.HeightRound()
	if msg.Height == height {
		if !handler.peerIsInTheCommittee(initiator) {
			return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the committee")
		}
		p := handler.consMgr.RoundProposal(msg.Round)
		if p != nil {
			response := message.NewProposalMessage(p)
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
