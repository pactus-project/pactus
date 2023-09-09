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

func (handler *helloAckHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HelloAckMessage)
	handler.logger.Trace("parsing HelloAck message", "message", msg)

	if msg.ResponseCode != message.ResponseCodeOK {
		handler.logger.Warn("hello message rejected", "from", initiator, "reason", msg.Reason)

		handler.network.CloseConnection(initiator)
		return nil
	}
	handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeKnown)
	handler.logger.Debug("hello message acknowledged", "from", initiator)

	return nil
}

func (handler *helloAckHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHandshaking)
	return bdl
}
