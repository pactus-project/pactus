package firewall

import (
	"encoding/hex"

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
	Enabled bool
	network network.Network
	peerSet *peerset.PeerSet
	state   state.StateFacade
	logger  *logger.Logger
}

func NewFirewall(enabled bool, net network.Network, peerSet *peerset.PeerSet, state state.StateFacade, logger *logger.Logger) *Firewall {
	return &Firewall{
		Enabled: enabled,
		network: net,
		peerSet: peerSet,
		state:   state,
		logger:  logger,
	}
}

func (f *Firewall) OpenMessage(data []byte, from peer.ID) *message.Message {
	peer := f.peerSet.MustGetPeer(from)
	if f.shouldBanPeer(peer) {
		f.logger.Warn("Firewall: Peer banned", "pid", util.FingerprintPeerID(from))
		f.network.CloseConnection(peer.PeerID())
		return nil
	}

	peer.IncreaseReceivedMessage()
	peer.IncreaseReceivedBytes(len(data))

	msg := new(message.Message)
	if err := msg.Decode(data); err != nil {
		peer.IncreaseInvalidMessage()
		f.logger.Debug("Error decoding message", "from", util.FingerprintPeerID(from), "data", hex.EncodeToString(data), "err", err)

		return nil
	}

	if err := msg.SanityCheck(); err != nil {
		peer.IncreaseInvalidMessage()
		f.logger.Debug("Peer sent us invalid msg", "from", util.FingerprintPeerID(from), "msg", msg, "err", err)
		return nil
	}

	if f.Enabled {
		if f.shouldDropMessage(msg) {
			// TODO: A better way for handshaking
			// peer.IncreaseInvalidMessage()
			f.logger.Warn("Firewall: Message dropped", "msg", msg, "from", util.FingerprintPeerID(from))
			return nil
		}
	}

	return msg
}

func (f *Firewall) shouldDropMessage(msg *message.Message) bool {
	initiatorPeer := f.peerSet.MustGetPeer(msg.Initiator)
	switch initiatorPeer.Status() {
	case peerset.StatusCodeBanned:
		return true
	}

	return false
}

func (f *Firewall) shouldBanPeer(peer *peerset.Peer) bool {
	if !f.Enabled {
		return false
	}

	// Ban peers after sending many invalid messages
	if peer.ReceivedMessages() > 1000 {
		ratio := (peer.InvalidMessages() * 100) / peer.ReceivedMessages()
		if ratio >= 10 {
			peer.UpdateStatus(peerset.StatusCodeBanned)
		}
	}

	switch peer.Status() {
	case peerset.StatusCodeBanned:
		return true
	}
	return false
}
