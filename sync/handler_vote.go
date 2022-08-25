package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type voteHandler struct {
	*synchronizer
}

func newVoteHandler(sync *synchronizer) messageHandler {
	return &voteHandler{
		sync,
	}
}

func (handler *voteHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.VoteMessage)
	handler.logger.Trace("parsing Vote message", "message", msg)

	handler.consensus.AddVote(msg.Vote)

	return nil
}

func (handler *voteHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
