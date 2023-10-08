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
		height := msg.From
		for _, data := range msg.CommittedBlocksData {
			blk, err := block.FromBytes(data)
			if err != nil {
				return err
			}
			handler.cache.AddBlock(blk)

			height++
		}
		handler.cache.AddCertificate(msg.LastCertificate)
		err := handler.tryCommitBlocks()
		if err != nil {
			return err
		}
	}

	handler.updateSession(msg.SessionID, initiator, msg.ResponseCode)

	return nil
}

func (handler *blocksResponseHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.CompressIt()

	return bdl
}

func (handler *blocksResponseHandler) updateSession(sessionID int, pid peer.ID, code message.ResponseCode) {
	s := handler.peerSet.FindSession(sessionID)
	if s == nil {
		// TODO: test me
		handler.logger.Debug("session not found or closed", "session-id", sessionID)
		return
	}

	if s.PeerID() != pid {
		// TODO: test me
		handler.logger.Warn("unknown peer", "session-id", sessionID, "pid", pid)
		return
	}

	s.SetLastResponseCode(code)

	switch code {
	case message.ResponseCodeRejected:
		handler.logger.Debug("session rejected, close session", "session-id", sessionID)
		handler.peerSet.CloseSession(sessionID)
		handler.updateBlockchain()

	case message.ResponseCodeMoreBlocks:
		handler.logger.Debug("peer responding us. keep session open", "session-id", sessionID)

	case message.ResponseCodeNoMoreBlocks:
		handler.logger.Debug("peer has no more block. close session", "session-id", sessionID)
		handler.peerSet.CloseSession(sessionID)
		handler.updateBlockchain()

	case message.ResponseCodeSynced:
		handler.logger.Debug("peer informed us we are synced. close session", "session-id", sessionID)
		handler.peerSet.CloseSession(sessionID)
		handler.moveConsensusToNewHeight()
	}
}
