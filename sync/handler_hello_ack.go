package sync

import (
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
