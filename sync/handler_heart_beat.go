package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type heartBeatHandler struct {
	*synchronizer
}

func newHeartBeatHandler(sync *synchronizer) payloadHandler {
	return &heartBeatHandler{
		sync,
	}
}

func (handler *heartBeatHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.HeartBeatPayload)
	handler.logger.Trace("Parsing heartbeat payload", "pld", pld)

	height, round := handler.consensus.HeightRound()

	if pld.Height == height {
		if pld.Round > round {
			if handler.weAreInTheCommittee() {
				handler.logger.Info("Our consensus is shorter than this peer.", "ours", round, "peer", pld.Round)

				query := payload.NewQueryVotesPayload(height, round)
				handler.broadcast(query)
			}
		} else if pld.Round < round {
			handler.logger.Trace("Our consensus is ahead of this peer.", "ours", round, "peer", pld.Round)
		} else {
			handler.logger.Trace("Our consensus is about the same round with this peer.", "ours", round, "peer", pld.Round)
		}
	}

	peer := handler.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.Height - 1)
	handler.peerSet.UpdateMaxClaimedHeight(pld.Height - 1)

	return nil
}

func (handler *heartBeatHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
