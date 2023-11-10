package sync

import (
	"fmt"
	"math"
	"time"

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
	handler.logger.Trace("parsing Hello message", "msg", msg)

	handler.logger.Debug("updating peer info",
		"pid", msg.PeerID,
		"moniker", msg.Moniker,
		"services", msg.Services)

	handler.peerSet.UpdateInfo(initiator,
		msg.Moniker,
		msg.Agent,
		msg.PublicKeys,
		msg.Services)

	if msg.PeerID != initiator {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("peer ID is not matched, expected: %v, got: %v",
				msg.PeerID, initiator))

		return handler.acknowledge(response, initiator)
	}

	if msg.GenesisHash != handler.state.Genesis().Hash() {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("invalid genesis hash, expected: %v, got: %v",
				handler.state.Genesis().Hash(), msg.GenesisHash))

		return handler.acknowledge(response, initiator)
	}

	if math.Abs(time.Since(msg.MyTime()).Seconds()) > 10 {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			"time discrepancy exceeds 10 seconds")

		return handler.acknowledge(response, initiator)
	}

	handler.peerSet.UpdateHeight(initiator, msg.Height, msg.BlockHash)
	handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeKnown)

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

		handler.logger.Debug("rejecting hello message", "msg", msg,
			"to", to, "reason", msg.Reason)
		handler.network.CloseConnection(to)
	} else {
		handler.logger.Info("acknowledging hello message", "msg", msg,
			"to", to)
	}

	return handler.sendTo(msg, to)
}
