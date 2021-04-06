package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

type aleykHandler struct {
	*HandlerContext
}

func NewAleykHandler(ctx *HandlerContext) Handler {
	return &aleykHandler{
		ctx,
	}
}

func (h *aleykHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.AleykPayload)
	h.logger.Trace("Parsing Aleyk payload", "pld", pld)

	if pld.ResponseCode != payload.ResponseCodeOK {
		h.logger.Warn("Our Salam is not welcomed!", "message", pld.ResponseMessage)
	} else {
		peer := h.peerSet.MustGetPeer(initiator)
		peer.UpdateMoniker(pld.Moniker)
		peer.UpdateHeight(pld.Height)
		peer.UpdateNodeVersion(pld.NodeVersion)
		peer.UpdatePublicKey(pld.PublicKey)
		peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

		h.peerSet.UpdateMaxClaimedHeight(pld.Height)
	}

	return nil
}
