package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
)

type helloHandler struct {
	*synchronizer
}

func newHelloHandler(sync *synchronizer) messageHandler {
	return &helloHandler{
		sync,
	}
}

func (handler *helloHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.HelloMessage)
	handler.logger.Trace("parsing Hello message", "message", msg)

	if msg.PeerID != initiator {
		handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage,
			"peer ID is not same as initiator for Hello message, expected: %v, got: %v",
			msg.PeerID, initiator)
	}

	if !msg.GenesisHash.EqualsTo(handler.state.Genesis().Hash()) {
		handler.peerSet.UpdateStatus(initiator, peerset.StatusCodeBanned)
		return errors.Errorf(errors.ErrInvalidMessage,
			"received a message from different chain, expected: %v, got: %v",
			handler.state.Genesis().Hash(), msg.GenesisHash)
	}

	handler.logger.Debug("updating peer info",
		"pid", initiator,
		"moniker", msg.Moniker,
		"flags", msg.Flags)

	handler.peerSet.UpdatePeerInfo(initiator,
		peerset.StatusCodeKnown,
		msg.Moniker,
		msg.Agent,
		msg.PublicKey,
		util.IsFlagSet(msg.Flags, message.FlagNodeNetwork))
	handler.peerSet.UpdateHeight(initiator, msg.Height, msg.BlockHash)

	if !util.IsFlagSet(msg.Flags, message.FlagHelloAck) {
		// TODO: Sends response only if there is a direct connection between two peers.
		// TODO: check if we have handshaked before. Ignore responding again
		// Response to Hello
		handler.sayHello(true, initiator)
	}

	handler.updateBlockchain()

	return nil
}

func (handler *helloHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	bdl := bundle.NewBundle(handler.SelfID(), m)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagHelloMessage)
	return bdl
}
