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
		handler.logger.Warn("blocks request is rejected", "pid", initiator, "reason", msg.Reason)
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
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.CompressIt()

	return bdl
}

func (handler *blocksResponseHandler) updateSession(sessionID int, code message.ResponseCode) {
	s := handler.peerSet.FindSession(sessionID)
	if s == nil {
		// TODO: test me
		// Probably session was expired before.
		handler.logger.Debug("session not found or closed", "sid", sessionID)
		return
	}

	switch code {
	case message.ResponseCodeRejected:
		handler.logger.Debug("session rejected, uncompleted session", "sid", sessionID, "peer", s.PeerID())
		s.SetUncompleted() // TODO: test me
		handler.updateBlockchain()

	case message.ResponseCodeMoreBlocks:
		handler.logger.Debug("peer responding us. keep session open", "sid", sessionID, "peer", s.PeerID())

	case message.ResponseCodeNoMoreBlocks:
		handler.logger.Debug("peer has no more block. close session", "sid", sessionID, "peer", s.PeerID())
		handler.peerSet.RemoveSession(sessionID) // TODO: test me
		handler.updateBlockchain()

	case message.ResponseCodeSynced:
		handler.logger.Debug("peer informed us we are synced. close session", "sid", sessionID, "peer", s.PeerID())
		handler.peerSet.RemoveSession(sessionID) // TODO: test me
		handler.moveConsensusToNewHeight()
	}
}
