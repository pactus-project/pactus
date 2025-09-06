package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type queryProposalHandler struct {
	*synchronizer
}

func newQueryProposalHandler(sync *synchronizer) messageHandler {
	return &queryProposalHandler{
		sync,
	}
}

func (handler *queryProposalHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.QueryProposalMessage)
	handler.logger.Trace("parsing QueryProposal message", "msg", msg)

	prop := handler.getConsMgr().HandleQueryProposal(msg.Height, msg.Round)
	if prop != nil {
		response := message.NewProposalMessage(prop)
		handler.broadcast(response)
	}
}

func (*queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
