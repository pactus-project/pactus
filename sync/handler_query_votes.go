package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type queryVoteHandler struct {
	*synchronizer
}

func newQueryVoteHandler(sync *synchronizer) messageHandler {
	return &queryVoteHandler{
		sync,
	}
}

func (handler *queryVoteHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.QueryVoteMessage)
	handler.logger.Trace("parsing QueryVote message", "msg", msg)

	v := handler.getConsMgr().HandleQueryVote(msg.Height, msg.Round)
	if v != nil {
		response := message.NewVoteMessage(v)
		handler.broadcast(response)
	}
}

func (*queryVoteHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
