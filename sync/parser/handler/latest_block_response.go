package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type latestBlocksResponseHandler struct {
	*HandlerContext
}

func NewLatestBlocksResponseHandler(ctx *HandlerContext) Handler {
	return &latestBlocksResponseHandler{
		ctx,
	}
}

func (h *latestBlocksResponseHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.LatestBlocksResponsePayload)
	h.logger.Trace("Parsing latest blocks response payload", "pld", pld)

	ourHeight := h.state.LastBlockHeight()
	if pld.To() == 0 || ourHeight < pld.To() {
		h.cache.AddCertificate(pld.LastCertificate)
		h.addBlocksToCache(pld.From, pld.Blocks)
		h.addTransactionsToCache(pld.Transactions)
		h.tryCommitBlocks()
	}
	h.updateSession(pld.ResponseCode, pld.SessionID, initiator, pld.Target)

	return nil
}
