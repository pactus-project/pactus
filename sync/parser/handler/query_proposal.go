package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryProposalHandler struct {
	*HandlerContext
}

func NewQueryProposalHandler(ctx *HandlerContext) Handler {
	return &queryProposalHandler{
		ctx,
	}
}

func (h *queryProposalHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryProposalPayload)
	h.logger.Trace("Parsing query proposal payload", "pld", pld)

	if !h.peerIsInTheCommittee(initiator) {
		return errors.Errorf(errors.ErrInvalidMessage, "Peers is not in the commmittee")
	}

	height, _ := h.consensus.HeightRound()
	if pld.Height == height {
		p := h.consensus.RoundProposal(pld.Round)
		if p != nil {
			response := payload.NewProposalPayload(*p)
			h.publishFn(response)
		}
	}

	return nil
}
