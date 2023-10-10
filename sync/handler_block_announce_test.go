package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastHeight := td.state.LastBlockHeight()
	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1)

	blk2, cert2 := td.GenerateTestBlock(lastHeight + 2)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2)

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))

		assert.Equal(t, td.sync.state.LastBlockHeight(), lastHeight)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg1, pid))

		assert.Equal(t, td.sync.state.LastBlockHeight(), lastHeight+2)
	})
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	td.sync.broadcast(msg)

	msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
	assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height(), msg.Certificate.Height())
}
