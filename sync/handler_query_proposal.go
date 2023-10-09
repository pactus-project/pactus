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

func (handler *queryProposalHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.QueryProposalMessage)
	handler.logger.Trace("parsing QueryProposal message", "message", msg)

	height, _ := handler.consMgr.HeightRound()
	if msg.Height == height {
		prop := handler.consMgr.Proposal()
		if prop != nil {
			response := message.NewProposalMessage(prop)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)

	return bdl
}
