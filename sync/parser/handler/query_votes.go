package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryVotesHandler struct {
	*HandlerContext
}

func NewQueryVotesHandler(ctx *HandlerContext) Handler {
	return &queryVotesHandler{
		ctx,
	}
}

func (h *queryVotesHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryVotesPayload)
	h.logger.Trace("Parsing query votes payload", "pld", pld)

	if !h.peerIsInTheCommittee(initiator) {
		return errors.Errorf(errors.ErrInvalidMessage, "Peers is not in the commmittee")
	}

	height, _ := h.consensus.HeightRound()
	if pld.Height == height {
		v := h.consensus.PickRandomVote()
		if v != nil {
			response := payload.NewVotePayload(*v)
			h.publishFn(response)
		}
	}

	return nil
}
