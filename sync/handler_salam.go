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
	handler.logger.Trace("Parsing Salam payload", "pld", pld)

	peer := handler.peerSet.MustGetPeer(initiator)

	if !pld.GenesisHash.EqualsTo(handler.state.GenesisHash()) {
		handler.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash, "peer", util.FingerprintPeerID(initiator))
		// Response to salam
		peer.UpdateStatus(peerset.StatusCodeBanned)
		handler.broadcastAleyk(initiator, payload.ResponseCodeRejected, "Invalid genesis hash")
		return nil
	}

	peer.UpdateStatus(peerset.StatusCodeOK)
	peer.UpdateMoniker(pld.Moniker)
	peer.UpdateHeight(pld.Height)
	peer.UpdateAgent(pld.Agent)
	peer.UpdatePublicKey(pld.PublicKey)
	peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

	handler.peerSet.UpdateMaxClaimedHeight(pld.Height)

	// Response to salam
	handler.broadcastAleyk(initiator, payload.ResponseCodeOK, "Welcome!")

	handler.updateBlokchain()

	return nil
}

func (handler *salamHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}

func (handler *salamHandler) broadcastAleyk(target peer.ID, code payload.ResponseCode, resMsg string) {
	flags := 0
	if handler.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	response := payload.NewAleykPayload(
		handler.config.Moniker,
		handler.signer.PublicKey(),
		handler.signer.SignData(handler.signer.PublicKey().RawBytes()),
		handler.state.LastBlockHeight(),
		flags,
		target,
		code,
		resMsg)

	handler.broadcast(response)
}
