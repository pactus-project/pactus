package sync

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
)

type helloHandler struct {
	*synchronizer
}

func newHelloHandler(sync *synchronizer) messageHandler {
	return &helloHandler{
		sync,
	}
}

func (handler *helloHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HelloMessage)
	handler.logger.Trace("parsing Hello message", "message", msg)

	if msg.PeerID != initiator {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("peer ID is not matched, expected: %v, got: %v", msg.PeerID, initiator))

		return handler.acknowledge(response, initiator)
	}

	if !msg.GenesisHash.EqualsTo(handler.state.Genesis().Hash()) {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("peer ID is not matched, expected: %v, got: %v", msg.PeerID, initiator))

		return handler.acknowledge(response, initiator)
	}

	handler.logger.Debug("updating peer info",
		"pid", initiator,
		"moniker", msg.Moniker,
		"services", msg.Services)

	handler.peerSet.UpdateInfo(initiator,
		msg.Moniker,
		msg.Agent,
		msg.PublicKeys,
		msg.Services)
	handler.peerSet.UpdateHeight(initiator, msg.Height, msg.BlockHash)
	handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeConnected)

	response := message.NewHelloAckMessage(message.ResponseCodeOK, "Ok")
	return handler.acknowledge(response, initiator)
}

func (handler *helloHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHandshaking)
	return bdl
}

func (handler *helloHandler) acknowledge(msg *message.HelloAckMessage, to peer.ID) error {
	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.peerSet.UpdateStatus(to, peerset.StatusCodeBanned)

		handler.logger.Warn("rejecting hello message", "message", msg, "to", to, "reason", msg.Reason)
	} else {
		handler.logger.Info("acknowledging hello message", "message", msg, "to", to)
	}

	return handler.sendTo(msg, to)
}
