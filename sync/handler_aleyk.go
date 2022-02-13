package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

type aleykHandler struct {
	*synchronizer
}

func newAleykHandler(sync *synchronizer) payloadHandler {
	return &aleykHandler{
		sync,
	}
}

func (handler *aleykHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.AleykPayload)
	handler.logger.Trace("parsing Aleyk payload", "pld", pld)

	peer := handler.peerSet.MustGetPeer(initiator)

	if pld.PeerID != initiator {
		peer.UpdateStatus(peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage, "Peer ID is not same as initiator for Aleyk message. expected: %v, got: %v",
			pld.PeerID, initiator)
	}

	if pld.ResponseTarget == handler.SelfID() {
		if pld.ResponseCode == payload.ResponseCodeOK {
			peer.UpdateStatus(peerset.StatusCodeGood)
		}
	}

	peer.UpdateMoniker(pld.Moniker)
	peer.UpdateHeight(pld.Height)
	peer.UpdateAgent(pld.Agent)
	peer.UpdatePublicKey(pld.PublicKey)
	peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

	handler.peerSet.UpdateMaxClaimedHeight(pld.Height)
	handler.updateBlokchain()

	return nil
}

func (handler *aleykHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
