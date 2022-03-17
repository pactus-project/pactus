package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

type helloHandler struct {
	*synchronizer
}

func newHelloHandler(sync *synchronizer) messageHandler {
	return &helloHandler{
		sync,
	}
}

func (handler *helloHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HelloMessage)
	handler.logger.Trace("parsing Hello message", "msg", msg)

	if msg.PeerID != initiator {
		handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage, "peer ID is not same as initiator for Hello message, expected: %v, got: %v",
			msg.PeerID, initiator)
	}

	if !msg.GenesisHash.EqualsTo(handler.state.GenesisHash()) {
		handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage, "received a message from different chain, expected: %v, got: %v",
			handler.state.GenesisHash(), msg.GenesisHash)
	}

	handler.peerSet.UpdatePeerInfo(initiator,
		peerset.StatusCodeKnown,
		msg.Moniker,
		msg.Agent,
		msg.PublicKey,
		util.IsFlagSet(msg.Flags, message.FlagNodeNetwork))
	handler.peerSet.UpdateHeight(initiator, msg.Height)

	if util.IsFlagSet(msg.Flags, message.FlagNeedResponse) {
		// TODO: Sends response only if there is a direct connection between two peers.
		// TODO: check if we have handshaked before. Ignore responding again
		// Response to Hello
		handler.sayHello(false)
	}

	handler.updateBlokchain()

	return nil
}

func (handler *helloHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	msg := bundle.NewBundle(handler.SelfID(), m)
	msg.Flags = util.SetFlag(msg.Flags, bundle.BundleFlagHelloMessage)
	return msg
}
