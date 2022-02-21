package firewall

import (
	"bytes"
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

func (f *Firewall) OpenGossipMessage(data []byte, source peer.ID, from peer.ID) *message.Message {
	if from != source {
		fromPeer := f.peerSet.MustGetPeer(from)
		if f.peerIsBanned(fromPeer) {
			f.logger.Warn("firewall: from peer banned", "from", util.FingerprintPeerID(from))
			f.closeConnection(from)
			return nil
		}
	}

	msg, err := f.openMessage(bytes.NewReader(data), source)
	if err != nil {
		f.logger.Warn("firewall: unable to open a gossip message", "err", err)
		f.closeConnection(from)
		return nil
	}

	// TODO: check if gossip flag is set
	// TODO: check if payload is a gossip payload

	return msg
}

func (f *Firewall) OpenStreamMessage(r io.Reader, from peer.ID) *message.Message {
	msg, err := f.openMessage(r, from)
	if err != nil {
		f.logger.Warn("firewall: unable to open a stream message", "err", err)
		f.closeConnection(from)
		return nil
	}

	// TODO: check if gossip flag is NOT set
	// TODO: check if payload is a stream payload

	return msg
}

func (f *Firewall) openMessage(r io.Reader, source peer.ID) (*message.Message, error) {
	peer := f.peerSet.MustGetPeer(source)
	peer.IncreaseReceivedMessage()

	if f.peerIsBanned(peer) {
		// If there is any connection to the source peer, close it
		f.closeConnection(source)
		return nil, errors.Errorf(errors.ErrInvalidMessage, "Source peer is banned: %s", source)
	}

	msg, err := f.decodeMessage(r, peer)
	if err != nil {
		peer.IncreaseInvalidMessage()
		return nil, err
	}

	if err := f.checkMessage(msg, peer); err != nil {
		peer.IncreaseInvalidMessage()
		return nil, err
	}

	return msg, nil
}

func (f *Firewall) decodeMessage(r io.Reader, source *peerset.Peer) (*message.Message, error) {
	msg := new(message.Message)
	bytesRead, err := msg.Decode(r)
	source.IncreaseReceivedBytes(bytesRead)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return msg, nil
}

func (f *Firewall) checkMessage(msg *message.Message, source *peerset.Peer) error {
	if err := msg.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	if msg.Initiator != source.PeerID() {
		return errors.Errorf(errors.ErrInvalidMessage,
			"source is not same as initiator. source: %v, initiator: %v", source.PeerID(), msg.Initiator)
	}

	return nil
}

func (f *Firewall) peerIsBanned(peer *peerset.Peer) bool {
	if !f.config.Enabled {
		return false
	}

	switch peer.Status() {
	case peerset.StatusCodeBanned:
		return true
	}
	return false
}

func (f *Firewall) closeConnection(pid peer.ID) {
	if !f.config.Enabled {
		return
	}

	f.network.CloseConnection(pid)
}
