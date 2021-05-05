package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

type salamHandler struct {
	*synchronizer
}

func newSalamHandler(sync *synchronizer) payloadHandler {
	return &salamHandler{
		sync,
	}
}

func (handler *salamHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.SalamPayload)
	handler.logger.Trace("Parsing salam payload", "pld", pld)

	peer := handler.peerSet.MustGetPeer(initiator)

	if !pld.GenesisHash.EqualsTo(handler.state.GenesisHash()) {
		handler.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash, "peer", util.FingerprintPeerID(initiator))
		// Response to salam
		peer.UpdateStatus(peerset.StatusCodeBanned)
		handler.broadcastAleyk(payload.ResponseCodeRejected, "Invalid genesis hash")
		return nil
	}

	peer.UpdateStatus(peerset.StatusCodeOK)
	peer.UpdateMoniker(pld.Moniker)
	peer.UpdateHeight(pld.Height)
	peer.UpdateNodeVersion(pld.NodeVersion)
	peer.UpdatePublicKey(pld.PublicKey)
	peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

	handler.peerSet.UpdateMaxClaimedHeight(pld.Height)

	// Response to salam
	handler.broadcastAleyk(payload.ResponseCodeOK, "Welcome!")

	return nil
}

func (handler *salamHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
