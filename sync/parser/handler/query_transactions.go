package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryTransactionsHandler struct {
	*HandlerContext
}

func NewQueryTransactionsHandler(ctx *HandlerContext) Handler {
	return &queryTransactionsHandler{
		ctx,
	}
}

func (h *queryTransactionsHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryTransactionsPayload)
	h.logger.Trace("Parsing query transactions payload", "pld", pld)

	if !h.peerIsInTheCommittee(initiator) {
		return errors.Errorf(errors.ErrInvalidMessage, "Peers is not in the commmittee")
	}

	trxs := h.prepareTransactions(pld.IDs)
	if len(trxs) > 0 {
		response := payload.NewTransactionsPayload(trxs)
		h.publishFn(response)
	}

	return nil
}
