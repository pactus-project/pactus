package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type voteHandler struct {
	*synchronizer
}

func newVoteHandler(sync *synchronizer) messageHandler {
	return &voteHandler{
		sync,
	}
}

func (handler *voteHandler) ParseMessage(m message.Message, _ peer.ID) {
	msg := m.(*message.VoteMessage)
	handler.logger.Trace("parsing Vote message", "msg", msg)

	handler.getConsMgr().AddVote(msg.Vote)
}

func (*voteHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
