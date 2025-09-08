package sync

import (
	"fmt"
	"math"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

//
// Handshake protocol in Pactus
//
// Handshake protocol is used to exchange information between two peers
// when they first connect to each other. The handshake process involves
// the following steps:
//
// 1. When two peers connect, the outbound peer first sends a Hello message to the inbound peer.
// 2. The inbound peer responds with a HelloAck message.
// 3. Then the Inbound peer sends a Hello message to the outbound peer.
// 4. Finally, the outbound peer responds with a HelloAck message.
//
//                  |                        |
//    Outbound Peer |                        | Inbound Peer
//                  |                        |
//        Connected |                        | Connected
//                  |                        |
//                  |         Hello          |
//                  | ---------------------> |
//                  |         HelloAck       |
//                  | <--------------------- |
//                  |                        |
//                  |         Hello          |
//            Known | <--------------------- |
//                  |         HelloAck       |
//                  | ---------------------> | Known
//                  |                        |
//
// After the handshake process is complete, both peers have exchanged
// information about their status, such as their latest block height,
// genesis hash, and supported services. This information is used to
// determine whether the peers are compatible and can communicate
// effectively.
//

type helloHandler struct {
	*synchronizer
}

func newHelloHandler(sync *synchronizer) messageHandler {
	return &helloHandler{
		sync,
	}
}

func (handler *helloHandler) ParseMessage(m message.Message, pid peer.ID) {
	msg := m.(*message.HelloMessage)
	handler.logger.Trace("parsing Hello message", "msg", msg)

	handler.logger.Debug("updating peer info",
		"pid", pid,
		"moniker", msg.Moniker,
		"services", msg.Services)

	handler.peerSet.UpdateInfo(pid,
		msg.Moniker,
		msg.Agent,
		msg.PublicKeys,
		msg.Services)
	handler.peerSet.UpdateHeight(pid, msg.Height, msg.BlockHash)

	if msg.GenesisHash != handler.state.Genesis().Hash() {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			fmt.Sprintf("invalid genesis hash, expected: %v, got: %v",
				handler.state.Genesis().Hash(), msg.GenesisHash), 0)

		handler.acknowledge(response, pid)

		return
	}

	if math.Abs(time.Since(msg.MyTime()).Seconds()) > 10 {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			"time discrepancy exceeds 10 seconds", 0)

		handler.acknowledge(response, pid)

		return
	}

	agent, _ := version.ParseAgent(msg.Agent)
	if agent.Version.Compare(handler.config.LatestSupportingVer) == -1 {
		response := message.NewHelloAckMessage(message.ResponseCodeRejected,
			"not supporting version", 0)

		handler.acknowledge(response, pid)

		return
	}

	for _, pub := range msg.PublicKeys {
		handler.state.UpdateValidatorProtocolVersion(pub.ValidatorAddress(), agent.ProtocolVersion)
	}

	peer := handler.peerSet.GetPeer(pid)

	switch peer.Direction {
	case lp2pnetwork.DirUnknown:
		handler.logger.Warn("received unexpected Hello message",
			"pid", pid, "direction", peer.Direction)
		handler.network.CloseConnection(pid)

		return

	case lp2pnetwork.DirInbound:
		if peer.OutboundHelloSent {
			handler.logger.Warn("received unexpected Hello message",
				"pid", pid, "direction", peer.Direction)
			handler.network.CloseConnection(pid)

			return
		}

		// Mark that we've received the hello message from the outbound peer
		handler.peerSet.UpdateOutboundHelloSent(pid, true)

		handler.logger.Info("sending Hello message (inbound)", "to", pid)
		handler.sayHello(pid)

	case lp2pnetwork.DirOutbound:
		if !peer.OutboundHelloSent {
			handler.logger.Warn("received unexpected Hello message",
				"pid", pid, "direction", peer.Direction)
			handler.network.CloseConnection(pid)

			return
		}

		handler.peerSet.UpdateStatus(pid, status.StatusKnown)
	}

	response := message.NewHelloAckMessage(message.ResponseCodeOK, "Ok", handler.state.LastBlockHeight())
	handler.acknowledge(response, pid)
}

func (*helloHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHandshaking)

	return bdl
}

func (handler *helloHandler) acknowledge(msg *message.HelloAckMessage, pid peer.ID) {
	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.logger.Info("rejecting hello message", "msg", msg,
			"pid", pid, "reason", msg.Reason)

		handler.sendTo(msg, pid)
		handler.peerSet.UpdateStatus(pid, status.StatusBanned)
	} else {
		handler.logger.Info("acknowledging hello message", "msg", msg,
			"pid", pid)

		handler.sendTo(msg, pid)
	}
}
