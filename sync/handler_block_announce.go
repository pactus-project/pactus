package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type blockAnnounceHandler struct {
	*synchronizer
}

func newBlockAnnounceHandler(sync *synchronizer) messageHandler {
	return &blockAnnounceHandler{
		sync,
	}
}

func (handler *blockAnnounceHandler) ParseMessage(m message.Message, pid peer.ID) error {
	msg := m.(*message.BlockAnnounceMessage)
	handler.logger.Trace("parsing BlockAnnounce message", "msg", msg)

	if handler.cache.HasBlockInCache(msg.Height()) {
		// We have processed this block before.

		return nil
	}

	handler.cache.AddCertificate(msg.Certificate)
	handler.cache.AddBlock(msg.Block)

	err := handler.tryCommitBlocks()
	if err != nil {
		return err
	}
	handler.moveConsensusToNewHeight()

	handler.peerSet.UpdateHeight(pid, msg.Height(), msg.Block.Hash())
	handler.updateBlockchain()

	return nil
}

func (handler *blockAnnounceHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)

	return bdl
}
