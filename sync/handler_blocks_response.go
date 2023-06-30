package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
)

type blocksResponseHandler struct {
	*synchronizer
}

func newBlocksResponseHandler(sync *synchronizer) messageHandler {
	return &blocksResponseHandler{
		sync,
	}
}

func (handler *blocksResponseHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.BlocksResponseMessage)
	handler.logger.Trace("parsing BlocksResponse message", "message", msg)

	if msg.IsRequestRejected() {
		handler.logger.Warn("blocks request is rejected", "pid", initiator, "response", msg.ResponseCode)
	} else {
		height := msg.From
		for _, d := range msg.BlocksData {
			b, err := block.FromBytes(d)
			if err != nil {
				return err
			}
			if err := b.SanityCheck(); err != nil {
				return err
			}
			handler.cache.AddBlock(height, b)
			height++
		}
		handler.cache.AddCertificate(msg.From, msg.LastCertificate)
		handler.tryCommitBlocks()
	}
	handler.updateSession(msg.SessionID, initiator, msg.ResponseCode)

	return nil
}

func (handler *blocksResponseHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.CompressIt()

	return bdl
}
