package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type voteHandler struct {
	*synchronizer
}

func newVoteHandler(sync *synchronizer) payloadHandler {
	return &voteHandler{
		sync,
	}
}

func (handler *voteHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.VotePayload)
	handler.logger.Trace("Parsing vote payload", "pld", pld)

	handler.consensus.AddVote(pld.Vote)

	return nil
}

func (handler *voteHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
