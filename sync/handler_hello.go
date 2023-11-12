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

func (handler *helloHandler) ParseMessage(m message.Message, pid peer.ID) error {
	msg := m.(*message.HelloMessage)
	handler.logger.Trace("parsing Hello message", "msg", msg)

	handler.logger.Debug("updating peer info",
		"pid", msg.PeerID,
		"moniker", msg.Moniker,
		"services", msg.Services)

	handler.peerSet.UpdateInfo(pid,
		msg.Moniker,
		msg.Agent,
		msg.PublicKeys,
		msg.Services)

	if msg.PeerID != pid {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("peer ID is not matched, expected: %v, got: %v",
				msg.PeerID, pid), 0)

		return handler.acknowledge(response, pid)
	}

	if msg.GenesisHash != handler.state.Genesis().Hash() {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("invalid genesis hash, expected: %v, got: %v",
				handler.state.Genesis().Hash(), msg.GenesisHash), 0)

		return handler.acknowledge(response, pid)
	}

	if math.Abs(time.Since(msg.MyTime()).Seconds()) > 10 {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			"time discrepancy exceeds 10 seconds", 0)

		return handler.acknowledge(response, pid)
	}

	handler.peerSet.UpdateHeight(pid, msg.Height, msg.BlockHash)
	handler.peerSet.UpdateStatus(pid, peerset.StatusCodeKnown)

	if msg.Services.IsGossip() {
		handler.network.Protect(msg.PeerID, "GOSSIP")
	}

	response := message.NewHelloAckMessage(message.ResponseCodeOK, "Ok", handler.state.LastBlockHeight())
	return handler.acknowledge(response, pid)
}

func (handler *helloHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)
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
