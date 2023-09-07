package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloAckMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Receiving HelloAck message: Rejected hello",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeRejected, "rejected")

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
		})

	t.Run("Receiving HelloAck message: OK hello",
		func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok")

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.checkPeerStatus(t, pid, peerset.StatusCodeKnown)
		})
}
