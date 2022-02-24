package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type heartBeatHandler struct {
	*synchronizer
}

func newHeartBeatHandler(sync *synchronizer) messageHandler {
	return &heartBeatHandler{
		sync,
	}
}

func (handler *heartBeatHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HeartBeatMessage)
	handler.logger.Trace("parsing HeartBeat message", "msg", msg)

	height, round := handler.consensus.HeightRound()

	if msg.Height == height {
		if msg.Round > round {
			if handler.weAreInTheCommittee() {
				handler.logger.Info("our consensus is behind of this peer", "ours", round, "peer", msg.Round)

				query := message.NewQueryVotesMessage(height, round)
				handler.broadcast(query)
			}
		} else if msg.Round < round {
			handler.logger.Trace("our consensus is ahead of this peer", "ours", round, "peer", msg.Round)
		} else {
			handler.logger.Trace("our consensus is at the same round with this peer", "ours", round, "peer", msg.Round)
		}
	}

	peer := handler.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(msg.Height - 1)
	handler.peerSet.UpdateMaxClaimedHeight(msg.Height - 1)

	return nil
}

func (handler *heartBeatHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
