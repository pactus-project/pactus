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

	if !handler.consMgr.HasActiveInstance() {
		handler.logger.Debug("ignoring QueryProposal, not active", "msg", msg)

		return
	}

	if !handler.consMgr.HasProposer() {
		handler.logger.Debug("ignoring QueryProposal, not proposer", "msg", msg)

		return
	}

	height, round := handler.consMgr.HeightRound()
	if msg.Height != height || msg.Round != round {
		handler.logger.Debug("ignoring QueryProposal, not same height/round", "msg", msg,
			"height", height, "round", round)

		return
	}

	prop := handler.consMgr.Proposal()
	if prop != nil {
		response := message.NewProposalMessage(prop)
		handler.broadcast(response)
	}
}

func (*queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
