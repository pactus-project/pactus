package firewall

import (
	"io"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
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
		f.logger.Warn("firewall: from peer banned", "pid", util.FingerprintPeerID(from))
		f.network.CloseConnection(from)
		return nil
	}

	peer.IncreaseReceivedMessage()
	msg := new(message.Message)
	bytesRead, err := msg.Decode(r)
	peer.IncreaseReceivedBytes(bytesRead)
	if err != nil {
		f.logger.Debug("error decoding message", "from", util.FingerprintPeerID(from), "err", err)
		peer.IncreaseInvalidMessage()
		return nil
	}

	if err := msg.SanityCheck(); err != nil {
		f.logger.Debug("peer sent us invalid msg", "from", util.FingerprintPeerID(from), "msg", msg, "err", err)
		peer.IncreaseInvalidMessage()
		return nil
	}

	if err := f.checkMessage(msg, from); err != nil {
		f.logger.Warn("firewall: message dropped", "err", err, "msg", msg)
		f.network.CloseConnection(from)
		peer.IncreaseInvalidMessage()
		return nil
	}

	return msg
}

func (f *Firewall) checkMessage(msg *message.Message, from peer.ID) error {
	if !f.config.Enabled {
		return nil
	}

	if msg.Initiator != from {
		return errors.Errorf(errors.ErrInvalidMessage,
			"source is not same as initiator. from: %v, initiator: %v", from, msg.Initiator)
	}

	return nil
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
