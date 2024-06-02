package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util/ratelimit"
)

type queryProposalHandler struct {
	*synchronizer

	rateLimit *ratelimit.RateLimit
}

func newQueryProposalHandler(sync *synchronizer) messageHandler {
	rateLimit := ratelimit.NewRateLimit(1, sync.config.QueryProposalWindow)

	return &queryProposalHandler{
		synchronizer: sync,
		rateLimit:    rateLimit,
	}
}

func (handler *queryProposalHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.QueryProposalMessage)
	handler.logger.Trace("parsing QueryProposal message", "msg", msg)

	if !handler.consMgr.HasActiveInstance() {
		handler.logger.Debug("ignoring QueryProposal, not active", "msg", msg)

		return nil
	}

	if !handler.consMgr.HasProposer() {
		handler.logger.Debug("ignoring QueryProposal, not proposer", "msg", msg)

		return nil
	}

	if !handler.rateLimit.AllowRequest() {
		handler.logger.Debug("ignoring QueryProposal, rate limit exceeded", "msg", msg)

		return nil
	}

	height, round := handler.consMgr.HeightRound()
	if msg.Height != height || msg.Round != round {
		handler.logger.Debug("ignoring QueryProposal, not same height/round", "msg", msg,
			"height", height, "round", round)

		return nil
	}

	prop := handler.consMgr.Proposal()
	if prop != nil {
		response := message.NewProposalMessage(prop)
		handler.broadcast(response)
	}

	return nil
}

func (handler *queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
