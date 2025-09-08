package sync

import (
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
)

type helloAckHandler struct {
	*synchronizer
}

func newHelloAckHandler(sync *synchronizer) messageHandler {
	return &helloAckHandler{
		sync,
	}
}

func (handler *helloAckHandler) ParseMessage(m message.Message, pid peer.ID) {
	msg := m.(*message.HelloAckMessage)
	handler.logger.Trace("parsing HelloAck message", "msg", msg)

	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.logger.Warn("hello message rejected",
			"from", pid, "reason", msg.Reason)
		handler.network.CloseConnection(pid)

		return
	}

	peer := handler.peerSet.GetPeer(pid)
	if peer == nil {
		handler.logger.Warn("received HelloAck from unknown peer", "pid", pid)
		handler.network.CloseConnection(pid)

		return
	}

	switch peer.Direction {
	case lp2pnetwork.DirUnknown:
		handler.logger.Warn("received unexpected HelloAc message",
			"pid", pid, "direction", peer.Direction)

		return

	case lp2pnetwork.DirInbound:
		if !peer.OutboundHelloSent {
			handler.logger.Warn("received unexpected HelloAc message",
				"pid", pid, "direction", peer.Direction)
			handler.network.CloseConnection(pid)

			return
		}

	case lp2pnetwork.DirOutbound:
		if !peer.OutboundHelloSent {
			handler.logger.Warn("received unexpected HelloAc message",
				"pid", pid, "direction", peer.Direction)
			handler.network.CloseConnection(pid)

			return
		}
	}

	handler.peerSet.UpdateStatus(pid, status.StatusKnown)
	handler.logger.Info("hello message acknowledged", "pid", pid, "reason", msg.Reason)

	if msg.Height > handler.state.LastBlockHeight() {
		handler.updateBlockchain()
	}
}

func (*helloAckHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHandshaking)

	return bdl
}
