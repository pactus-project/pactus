package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type blockAnnounceHandler struct {
	*HandlerContext
}

func NewBlockAnnounceHandler(ctx *HandlerContext) Handler {
	return &blockAnnounceHandler{
		ctx,
	}
}

func (h *blockAnnounceHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.BlockAnnouncePayload)
	h.logger.Trace("Parsing block announce payload", "pld", pld)

	h.cache.AddCertificate(&pld.Certificate)
	h.cache.AddBlock(pld.Height, &pld.Block)
	h.tryCommitBlocks()
	h.syncedFn()

	peer := h.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.Height)
	h.peerSet.UpdateMaxClaimedHeight(pld.Height)

	return nil
}
