package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestHandlerBlockAnnounceBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	blk, cert := td.GenerateTestBlock(td.RandHeight())
	msg := message.NewBlockAnnounceMessage(blk, cert, nil)
	td.sync.broadcast(msg)

	msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	assert.Equal(t, msg.Certificate.Height(), msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height())
}

func TestHandlerBlockAnnounceCacheBlock(t *testing.T) {
	td := setup(t, nil)

	height := td.RandHeight()
	blk1, cert1 := td.GenerateTestBlock(height)
	blk2, cert2 := td.GenerateTestBlock(height)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1, nil)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2, nil)

	td.receivingNewMessage(td.sync, msg1, td.RandPeerID())
	td.receivingNewMessage(td.sync, msg2, td.RandPeerID())

	cachedBlock := td.sync.cache.GetBlock(height)
	cachedCert := td.sync.cache.GetCertificate(height)
	assert.Equal(t, blk1, cachedBlock)
	assert.Equal(t, cert1, cachedCert)
}
