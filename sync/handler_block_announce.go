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

func (handler *blockAnnounceHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.BlockAnnounceMessage)
	handler.logger.Trace("parsing BlockAnnounce message", "message", msg)

	handler.cache.AddCertificate(msg.Height, msg.Certificate)
	handler.cache.AddBlock(msg.Height, msg.Block)
	handler.tryCommitBlocks()
	handler.moveConsensusToNewHeight()

	handler.peerSet.UpdateHeight(initiator, msg.Height)
	handler.updateBlockchain()

	return nil
}

func (handler *blockAnnounceHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending BlockAnnounce ignored. We are not in the committee")
		return nil
	}
	bdl := bundle.NewBundle(handler.SelfID(), m)

	return bdl
}
