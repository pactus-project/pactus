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

func (handler *blocksResponseHandler) ParseMessage(m message.Message, pid peer.ID) error {
	msg := m.(*message.BlocksResponseMessage)
	handler.logger.Trace("parsing BlocksResponse message", "msg", msg)

	if msg.IsRequestRejected() {
		handler.logger.Warn("blocks request is rejected", "pid", pid, "reason", msg.Reason, "sid", msg.SessionID)
	} else {
		// TODO:
		// It is good to check the latest height before adding blocks to the cache.
		// If they have already been committed, this message can be ignored.
		// Need to test!
		for _, data := range msg.CommittedBlocksData {
			blk, err := block.FromBytes(data)
			if err != nil {
				return err
			}
			handler.cache.AddBlock(blk)
		}
		handler.cache.AddCertificate(msg.LastCertificate)
		err := handler.tryCommitBlocks()
		if err != nil {
			return err
		}
	}

	handler.updateSession(msg.SessionID, msg.ResponseCode)

	return nil
}

func (handler *blocksResponseHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(m)
	bdl.CompressIt()

	return bdl
}

func (handler *blocksResponseHandler) updateSession(sid int, code message.ResponseCode) {
	switch code {
	case message.ResponseCodeOK:
		handler.logger.Debug("session accepted. keep session open", "sid", sid)
		handler.peerSet.UpdateSessionLastActivity(sid)

	case message.ResponseCodeRejected:
		handler.logger.Debug("session rejected, uncompleted session", "sid", sid)
		handler.peerSet.SetSessionUncompleted(sid)

	case message.ResponseCodeMoreBlocks:
		handler.logger.Debug("peer responding us. keep session open", "sid", sid)
		handler.peerSet.UpdateSessionLastActivity(sid)

	case message.ResponseCodeNoMoreBlocks:
		handler.logger.Debug("peer sent all blocks. close session", "sid", sid)
		handler.peerSet.SetSessionCompleted(sid) // TODO: test me
		handler.updateBlockchain()

	case message.ResponseCodeSynced:
		handler.logger.Debug("peer informed us we are synced. close session", "sid", sid)
		handler.peerSet.SetSessionCompleted(sid) // TODO: test me
		handler.moveConsensusToNewHeight()
	}
}
