package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type proposalHandler struct {
	*HandlerContext
}

func NewProposalHandler(ctx *HandlerContext) Handler {
	return &proposalHandler{
		ctx,
	}
}

func (h *proposalHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.ProposalPayload)
	h.logger.Trace("Parsing proposal payload", "pld", pld)

	h.cache.AddProposal(&pld.Proposal)
	h.consensus.SetProposal(&pld.Proposal)

	return nil
}
