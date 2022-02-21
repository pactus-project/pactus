package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

type helloHandler struct {
	*synchronizer
}

func newHelloHandler(sync *synchronizer) payloadHandler {
	return &helloHandler{
		sync,
	}
}

func (handler *helloHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.HelloPayload)
	handler.logger.Trace("parsing Hello payload", "pld", pld)

	peer := handler.peerSet.MustGetPeer(initiator)

	if pld.PeerID != initiator {
		peer.UpdateStatus(peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage, "Peer ID is not same as initiator for Hello message. expected: %v, got: %v",
			pld.PeerID, initiator)
	}

	if !pld.GenesisHash.EqualsTo(handler.state.GenesisHash()) {

		peer.UpdateStatus(peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage, "received a message from different chain, expected: %v, got: %v",
			handler.state.GenesisHash(), pld.GenesisHash)
	}

	peer.UpdateStatus(peerset.StatusCodeKnown)
	peer.UpdateMoniker(pld.Moniker)
	peer.UpdateHeight(pld.Height)
	peer.UpdateAgent(pld.Agent)
	peer.UpdatePublicKey(pld.PublicKey)
	peer.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, payload.FlagInitialBlockDownload))

	handler.peerSet.UpdateMaxClaimedHeight(pld.Height)

	if util.IsFlagSet(pld.Flags, payload.FlagNeedResponse) {
		// TODO: Sends response only if there is a direct connection between two peers.
		// TODO: check if we have handshaked before. Ignore responding again
		// Response to Hello
		handler.sayHello(false)
	}

	handler.updateBlokchain()

	return nil
}

func (handler *helloHandler) PrepareMessage(p payload.Payload) *message.Message {
	msg := message.NewMessage(handler.SelfID(), p)
	msg.Flags = util.SetFlag(msg.Flags, message.MsgFlagHelloMessage)
	return msg
}
