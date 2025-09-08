package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestHandlerBlockAnnounceParsingMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(10)

	pid := td.RandPeerID()
	lastHeight := td.state.LastBlockHeight()
	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1)

	blk2, cert2 := td.GenerateTestBlock(lastHeight + 2)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2)

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		td.receivingNewMessage(td.sync, msg2, pid)

		stateHeight := td.sync.state.LastBlockHeight()
		consHeight, _ := td.consMgr.HeightRound()
		assert.Equal(t, lastHeight, stateHeight)
		assert.Equal(t, lastHeight+1, consHeight)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		td.receivingNewMessage(td.sync, msg1, pid)

		assert.Equal(t, lastHeight+2, td.sync.state.LastBlockHeight())
	})
}

func TestHandlerBlockAnnounceBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert)
	td.sync.broadcast(msg)

	msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	assert.Equal(t, msg.Certificate.Height(), msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height())
}

func TestHandlerBlockAnnounceCacheBlock(t *testing.T) {
	td := setup(t, nil)

	height := td.RandHeight()
	blk1, cert1 := td.GenerateTestBlock(height)
	blk2, cert2 := td.GenerateTestBlock(height)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2)

	td.receivingNewMessage(td.sync, msg1, td.RandPeerID())
	td.receivingNewMessage(td.sync, msg2, td.RandPeerID())

	cachedBlock := td.sync.cache.GetBlock(height)
	cachedCert := td.sync.cache.GetCertificate(height)
	assert.Equal(t, blk1, cachedBlock)
	assert.Equal(t, cert1, cachedCert)
}
