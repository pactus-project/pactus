package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util/ratelimit"
)

type queryVotesHandler struct {
	*synchronizer

	rateLimit *ratelimit.RateLimit
}

func newQueryVotesHandler(sync *synchronizer) messageHandler {
	rateLimit := ratelimit.NewRateLimit(1, sync.config.QueryVoteWindow)

	return &queryVotesHandler{
		synchronizer: sync,
		rateLimit:    rateLimit,
	}
}

func (handler *queryVotesHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.QueryVotesMessage)
	handler.logger.Trace("parsing QueryVotes message", "msg", msg)

	if !handler.consMgr.HasActiveInstance() {
		handler.logger.Debug("ignoring QueryVotes, not active", "msg", msg)

		return nil
	}

	if !handler.rateLimit.AllowRequest() {
		handler.logger.Debug("ignoring QueryVotes, rate limit exceeded", "msg", msg)

		return nil
	}

	height, _ := handler.consMgr.HeightRound()
	if msg.Height != height {
		handler.logger.Debug("ignoring QueryVotes, not same height", "msg", msg,
			"height", height)

		return nil
	}

	v := handler.consMgr.PickRandomVote(msg.Round)
	if v != nil {
		response := message.NewVoteMessage(v)
		handler.broadcast(response)
	}

	return nil
}

func (handler *queryVotesHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
