package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type queryVotesHandler struct {
	*synchronizer
}

func newQueryVotesHandler(sync *synchronizer) messageHandler {
	return &queryVotesHandler{
		sync,
	}
}

func (handler *queryVotesHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.QueryVotesMessage)
	handler.logger.Trace("parsing QueryVotes message", "msg", msg)

	height, _ := handler.consMgr.HeightRound()
	if msg.Height == height {
		v := handler.consMgr.PickRandomVote(msg.Round)
		if v != nil {
			response := message.NewVoteMessage(v)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryVotesHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)

	return bdl
}
