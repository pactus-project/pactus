package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/certificate"
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

func TestInvalidBlockAnnounce(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	height := td.state.LastBlockHeight() + 1
	blk, _ := td.GenerateTestBlock(height)
	invCert := certificate.NewBlockCertificate(height, 0, false)
	msg := message.NewBlockAnnounceMessage(blk, invCert)

	err := td.receivingNewMessage(td.sync, msg, pid)
	assert.Error(t, err)
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	td.sync.broadcast(msg)

	msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height(), msg.Certificate.Height())
}
