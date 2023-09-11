package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
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
		height := msg.From
		for _, data := range msg.CommittedBlocksData {
			committedBlock := handler.state.MakeCommittedBlock(data, height, hash.UndefHash)
			b, err := committedBlock.ToBlock()
			if err != nil {
				return err
			}
			if err := b.BasicCheck(); err != nil {
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
