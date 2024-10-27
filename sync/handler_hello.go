package sync

import (
	"fmt"
	"math"
	"time"

	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

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

		handler.acknowledge(response, pid)

		return
	}

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

	handler.peerSet.UpdateHeight(pid, msg.Height, msg.BlockHash)
	handler.peerSet.UpdateStatus(pid, status.StatusConnected)

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
