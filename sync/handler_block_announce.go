package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type blockAnnounceHandler struct {
	*synchronizer
}

func newBlockAnnounceHandler(sync *synchronizer) payloadHandler {
	return &blockAnnounceHandler{
		sync,
	}
}

func (handler *blockAnnounceHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.BlockAnnouncePayload)
	handler.logger.Trace("parsing block announce payload", "pld", pld)

	handler.cache.AddCertificate(pld.Certificate)
	handler.cache.AddBlock(pld.Height, pld.Block)
	handler.tryCommitBlocks()
	handler.synced()

	peer := handler.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.Height)
	handler.peerSet.UpdateMaxClaimedHeight(pld.Height)

	handler.updateBlokchain()

	return nil
}

func (handler *blockAnnounceHandler) PrepareMessage(p payload.Payload) *message.Message {
	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending BlockAnnounce ignored. We are not in the committee")
		return nil
	}
	msg := message.NewMessage(handler.SelfID(), p)

	return msg
}
