package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type downloadRequestHandler struct {
	*HandlerContext
}

func NewDownloadRequestHandler(ctx *HandlerContext) Handler {
	return &downloadRequestHandler{
		ctx,
	}
}

func (h *downloadRequestHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.DownloadRequestPayload)
	h.logger.Trace("Parsing download request payload", "pld", pld)

	peer := h.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != h.selfID {
		return nil
	}
	if pld.To-pld.From > h.requestBlockInterval {
		return errors.Errorf(errors.ErrInvalidMessage, "Peer request interval is not acceptable: %v", pld.To-pld.From)
	}

	from := pld.From
	count := h.blockPerMessage

	for {
		if from+count >= pld.To {
			// Last packet has one extra block, for confirming last block
			count++
		}
		blocks, trxs := h.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		response := payload.NewDownloadResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, initiator, from, blocks, trxs)
		h.publishFn(response)

		from += len(blocks)
		if from >= pld.To {
			break
		}
	}

	response := payload.NewDownloadResponsePayload(payload.ResponseCodeNoMoreBlocks, pld.SessionID, initiator, 0, nil, nil)
	h.publishFn(response)

	return nil
}
