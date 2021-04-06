package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

type salamHandler struct {
	*HandlerContext
}

func NewSalamHandler(ctx *HandlerContext) Handler {
	return &salamHandler{
		ctx,
	}
}

func (h *salamHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.SalamPayload)
	h.logger.Trace("Parsing salam payload", "pld", pld)

	if !pld.GenesisHash.EqualsTo(h.state.GenesisHash()) {
		h.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash)
		// Response to salam
		h.responseSalam(payload.ResponseCodeRejected, "Invalid genesis hash")
		return nil
	}

	peer := h.peerSet.MustGetPeer(initiator)
	peer.UpdateMoniker(pld.Moniker)
	peer.UpdateHeight(pld.Height)
	peer.UpdateNodeVersion(pld.NodeVersion)
	peer.UpdatePublicKey(pld.PublicKey)
	peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

	h.peerSet.UpdateMaxClaimedHeight(pld.Height)

	// Response to salam
	h.responseSalam(payload.ResponseCodeOK, "Welcome!")

	return nil
}

func (h *salamHandler) responseSalam(code payload.ResponseCode, resMsg string) {
	flags := 0
	if h.initialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	response := payload.NewAleykPayload(
		code,
		resMsg,
		h.moniker,
		h.publicKey,
		h.state.LastBlockHeight(),
		flags)

	h.publishFn(response)
}
