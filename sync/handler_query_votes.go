package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type queryVotesHandler struct {
	*synchronizer
}

func newQueryVotesHandler(sync *synchronizer) messageHandler {
	return &queryVotesHandler{
		sync,
	}
}

func (handler *queryVotesHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.QueryVotesMessage)
	handler.logger.Trace("parsing QueryVotes message", "msg", msg)

	if !handler.consMgr.HasActiveInstance() {
		handler.logger.Debug("ignoring QueryVotes, not active", "msg", msg)

		return
	}

	height, _ := handler.consMgr.HeightRound()
	if msg.Height != height {
		handler.logger.Debug("ignoring QueryVotes, not same height", "msg", msg,
			"height", height)

		return
	}

	v := handler.consMgr.PickRandomVote(msg.Round)
	if v != nil {
		response := message.NewVoteMessage(v)
		handler.broadcast(response)
	}
}

func (*queryVotesHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
