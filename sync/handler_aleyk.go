package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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
	handler.logger.Trace("Parsing Aleyk payload", "pld", pld)

	if pld.ResponseCode != payload.ResponseCodeOK {
		handler.logger.Warn("Our Salam is not welcomed!", "message", pld.ResponseMessage)
	} else {
		peer := handler.peerSet.MustGetPeer(initiator)
		peer.UpdateMoniker(pld.Moniker)
		peer.UpdateHeight(pld.Height)
		peer.UpdateNodeVersion(pld.NodeVersion)
		peer.UpdatePublicKey(pld.PublicKey)
		peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

		handler.peerSet.UpdateMaxClaimedHeight(pld.Height)
	}

	return nil
}

func (handler *aleykHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
