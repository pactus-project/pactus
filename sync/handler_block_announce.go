package sync

import (
	"fmt"

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

	// The announcements come from any peer. They should be checked before they are cached,
	// otherwise a peer can replace the blocks and the certificates that we are waiting for.
	if err := handler.checkAnnouncement(msg); err != nil {
		handler.logger.Warn("invalid block announcement", "from", pid, "error", err)

		return
	}

	handler.peerSet.UpdateHeight(pid, msg.Height(), msg.Block.Hash())
	handler.cache.AddCertificate(msg.Certificate)
	handler.cache.AddBlock(msg.Block)

	handler.tryCommitBlocks()
	handler.moveConsensusToNewHeight()
	handler.updateBlockchain()
}

// checkAnnouncement checks that the certificate of the announcement is made for its block.
//
// The certificate can be verified only for the next block of the state, because the committee
// of the next blocks is not known yet. The other announcements are cached for the synchronization,
// and they are verified when they are committed.
func (handler *blockAnnounceHandler) checkAnnouncement(msg *message.BlockAnnounceMessage) error {
	if msg.Block.Height() != msg.Certificate.Height() {
		return fmt.Errorf("the block is for the height %d, but the certificate is for the height %d",
			msg.Block.Height(), msg.Certificate.Height())
	}

	if msg.Height() != handler.stateHeight()+1 {
		return nil
	}

	return msg.Certificate.ValidatePrecommit(handler.state.CommitteeValidators(), msg.Block.Hash())
}

func (*blockAnnounceHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
