package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
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

func (handler *helloAckHandler) ParseMessage(m message.Message, pid peer.ID) error {
	msg := m.(*message.HelloAckMessage)
	handler.logger.Trace("parsing HelloAck message", "msg", msg)

	if msg.ResponseCode != message.ResponseCodeOK {
		handler.logger.Warn("hello message rejected",
			"from", pid, "reason", msg.Reason)

		handler.network.CloseConnection(pid)

		return nil
	}

	handler.peerSet.UpdateStatus(pid, peerset.StatusCodeKnown)
	handler.logger.Debug("hello message acknowledged", "from", pid)

	if msg.Height > handler.state.LastBlockHeight() {
		handler.updateBlockchain()
	}

	return nil
}

func (handler *helloAckHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHandshaking)

	return bdl
}
