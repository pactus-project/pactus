package firewall

import (
	"io"

	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

// Firewall check packets before passing them to sync module
type Firewall struct {
	config  *Config
	network network.Network
	peerSet *peerset.PeerSet
	state   state.Facade
	logger  *logger.Logger
}

func NewFirewall(conf *Config, net network.Network, peerSet *peerset.PeerSet, state state.Facade, logger *logger.Logger) *Firewall {
	return &Firewall{
		config:  conf,
		network: net,
		peerSet: peerSet,
		state:   state,
		logger:  logger,
	}
}

func (f *Firewall) OpenMessage(r io.Reader, from peer.ID) *message.Message {
	peer := f.peerSet.MustGetPeer(from)
	if f.shouldBanPeer(peer) {
		f.logger.Warn("Firewall: Peer banned", "pid", util.FingerprintPeerID(from))
		f.network.CloseConnection(peer.PeerID())
		return nil
	}

	peer.IncreaseReceivedMessage()
	msg := new(message.Message)
	bytesRead, err := msg.Decode(r)
	peer.IncreaseReceivedBytes(bytesRead)
	if err != nil {
		peer.IncreaseInvalidMessage()
		f.logger.Debug("Error decoding message", "from", util.FingerprintPeerID(from), "err", err)

		return nil
	}

	if err := msg.SanityCheck(); err != nil {
		peer.IncreaseInvalidMessage()
		f.logger.Debug("Peer sent us invalid msg", "from", util.FingerprintPeerID(from), "msg", msg, "err", err)
		return nil
	}

	if f.shouldDropMessage(msg) {
		// TODO: A better way for handshaking
		peer.IncreaseInvalidMessage()
		f.logger.Warn("Firewall: Message dropped", "msg", msg, "from", util.FingerprintPeerID(from))
		return nil
	}

	return msg
}

func (f *Firewall) shouldDropMessage(msg *message.Message) bool {
	if !f.config.Enabled {
		return false
	}

	initiatorPeer := f.peerSet.MustGetPeer(msg.Initiator)
	switch initiatorPeer.Status() {
	case peerset.StatusCodeBanned:
		return true
	}

	return false
}

func (f *Firewall) shouldBanPeer(peer *peerset.Peer) bool {
	if !f.config.Enabled {
		return false
	}

	switch peer.Status() {
	case peerset.StatusCodeBanned:
		return true
	}
	return false
}
