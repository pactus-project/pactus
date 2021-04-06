package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type downloadResponseHandler struct {
	*HandlerContext
}

func NewDownloadResponseHandler(ctx *HandlerContext) Handler {
	return &downloadResponseHandler{
		ctx,
	}
}

func (h *downloadResponseHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.DownloadResponsePayload)
	h.logger.Trace("Parsing download response payload", "pld", pld)

	ourHeight := h.state.LastBlockHeight()
	if pld.To() == 0 || ourHeight < pld.To() {
		h.addBlocksToCache(pld.From, pld.Blocks)
		h.addTransactionsToCache(pld.Transactions)
		h.tryCommitBlocks()
	}
	h.updateSession(pld.ResponseCode, pld.SessionID, initiator, pld.Target)

	return nil
}
