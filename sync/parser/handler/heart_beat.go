package handler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type heartBeatHandler struct {
	*HandlerContext
}

func NewHeartBeatHandler(ctx *HandlerContext) Handler {
	return &heartBeatHandler{
		ctx,
	}
}

func (h *heartBeatHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.HeartBeatPayload)
	h.logger.Trace("Parsing heartbeat payload", "pld", pld)

	height, round := h.consensus.HeightRound()

	if pld.Height == height {
		if pld.Round > round {
			if h.weAreInTheCommittee() {
				h.logger.Info("Our consensus is behind of this peer.", "ours", round, "peer", pld.Round)

				q := payload.NewQueryVotesPayload(height, round)
				h.publishFn(q)
			}
		} else if pld.Round < round {
			h.logger.Trace("Our consensus is ahead of this peer.")
		} else {
			h.logger.Trace("Our consensus is at the same step with this peer.")
		}
	}

	peer := h.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.Height - 1)
	h.peerSet.UpdateMaxClaimedHeight(pld.Height - 1)

	return nil
}
