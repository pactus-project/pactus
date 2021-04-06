package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type voteHandler struct {
	*HandlerContext
}

func NewVoteHandler(ctx *HandlerContext) Handler {
	return &voteHandler{
		ctx,
	}
}

func (h *voteHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.VotePayload)
	h.logger.Trace("Parsing vote payload", "pld", pld)

	h.consensus.AddVote(&pld.Vote)

	return nil
}
