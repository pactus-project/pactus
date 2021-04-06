package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type transactionsHandler struct {
	*HandlerContext
}

func NewTransactionsHandler(ctx *HandlerContext) Handler {
	return &transactionsHandler{
		ctx,
	}
}

func (h *transactionsHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.TransactionsPayload)
	h.logger.Trace("Parsing transactions payload", "pld", pld)

	h.addTransactionsToCache(pld.Transactions)

	for _, trx := range pld.Transactions {
		if err := h.state.AddPendingTx(&trx); err != nil {
			h.logger.Debug("Cannot append transaction", "tx", trx, "err", err)

			// TODO: set peer as bad peer?
		}
	}

	return nil
}
