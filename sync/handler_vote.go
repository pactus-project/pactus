package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type voteHandler struct {
	*synchronizer
}

func newVoteHandler(sync *synchronizer) messageHandler {
	return &voteHandler{
		sync,
	}
}

func (handler *voteHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.VoteMessage)
	handler.logger.Trace("parsing Vote message", "message", msg)

	handler.consMgr.AddVote(msg.Vote)

	return nil
}

func (handler *voteHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
