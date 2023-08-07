package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type heartBeatHandler struct {
	*synchronizer
}

func newHeartBeatHandler(sync *synchronizer) messageHandler {
	return &heartBeatHandler{
		sync,
	}
}

func (handler *heartBeatHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HeartBeatMessage)
	handler.logger.Trace("parsing HeartBeat message", "message", msg)

	height, round := handler.consMgr.HeightRound()

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

	handler.peerSet.UpdateHeight(initiator, msg.Height, msg.PrevBlockHash)

	return nil
}

func (handler *heartBeatHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
