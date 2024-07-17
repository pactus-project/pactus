package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
)

func TestParsingHelloAckMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Receiving HelloAck message: Rejected hello",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeRejected, "rejected", td.RandHeight())

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusUnknown)
		})

	t.Run("Receiving HelloAck message: OK hello",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", td.RandHeight())

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusKnown)
		})
}
