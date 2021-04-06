package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type latestBlocksRequestHandler struct {
	*HandlerContext
}

func NewLatestBlocksRequestHandler(ctx *HandlerContext) Handler {
	return &latestBlocksRequestHandler{
		ctx,
	}
}

func (h *latestBlocksRequestHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.LatestBlocksRequestPayload)
	h.logger.Trace("Parsing latest blocks request payload", "pld", pld)

	peer := h.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != h.selfID {
		return nil
	}
	ourHeight := h.state.LastBlockHeight()
	if pld.From < ourHeight-h.requestBlockInterval {
		return errors.Errorf(errors.ErrInvalidMessage, "The request height is not acceptable: %v", pld.From)
	}
	from := pld.From
	count := h.blockPerMessage

	// Help peer to catch up
	for {
		blocks, trxs := h.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, initiator, from, blocks, trxs, nil)
		h.publishFn(response)

		from += len(blocks)
	}

	lastCertificate := h.state.LastCertificate()
	response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, initiator, from, nil, nil, lastCertificate)
	h.publishFn(response)

	return nil
}
