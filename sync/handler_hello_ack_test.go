package sync

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
)

func TestHandlerHelloAckParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Receiving HelloAck message: Rejected hello",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeRejected, "rejected", td.RandHeight())

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusUnknown)
		})

	t.Run("Receiving HelloAck message from unknown peer",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusUnknown)
		})
}

func TestHandlerHelloAckHandshaking(t *testing.T) {
	td := setup(t, nil)

	t.Run("Unknown Direction", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

		td.connectPeer(pid, lp2pnetwork.DirUnknown, false)
		td.receivingNewMessage(td.sync, msg, pid)
		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("Inbound Direction, Outbound Hello Not Sent", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

		td.connectPeer(pid, lp2pnetwork.DirInbound, false)
		td.receivingNewMessage(td.sync, msg, pid)
		td.shouldNotPublishAnyMessage(t)
		td.checkPeerStatus(t, pid, status.StatusConnected)
	})

	t.Run("Outbound Direction, Outbound Hello Not Sent", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

		td.connectPeer(pid, lp2pnetwork.DirOutbound, false)
		td.receivingNewMessage(td.sync, msg, pid)
		td.shouldNotPublishAnyMessage(t)
		td.checkPeerStatus(t, pid, status.StatusConnected)
	})

	t.Run("Inbound Direction, Outbound Hello Sent", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

		td.connectPeer(pid, lp2pnetwork.DirInbound, true)
		td.receivingNewMessage(td.sync, msg, pid)
		td.checkPeerStatus(t, pid, status.StatusKnown)
	})

	t.Run("Outbound Direction, Outbound Hello Sent", func(t *testing.T) {
		pid := td.RandPeerID()
		msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

		td.connectPeer(pid, lp2pnetwork.DirOutbound, true)
		td.receivingNewMessage(td.sync, msg, pid)
		td.shouldNotPublishAnyMessage(t)
		td.checkPeerStatus(t, pid, status.StatusKnown)
	})
}
