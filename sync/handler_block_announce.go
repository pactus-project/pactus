package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
)

type blockAnnounceHandler struct {
	*synchronizer
}

func newBlockAnnounceHandler(sync *synchronizer) messageHandler {
	return &blockAnnounceHandler{
		sync,
	}
}

func (handler *blockAnnounceHandler) ParseMessage(m message.Message, pid peer.ID) {
	msg := m.(*message.BlockAnnounceMessage)
	handler.logger.Trace("parsing BlockAnnounce message", "msg", msg)

	if handler.cache.HasBlockInCache(msg.Height()) {
		// We have processed this block before.

		return
	}

	handler.peerSet.UpdateHeight(pid, msg.Height(), msg.Block.Hash())
	handler.cache.AddCertificate(msg.Certificate)
	handler.cache.AddBlock(msg.Block)

	handler.tryCommitBlocks()
	handler.moveConsensusToNewHeight()
	handler.updateBlockchain()
}

func (*blockAnnounceHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
