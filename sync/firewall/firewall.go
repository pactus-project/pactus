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

func (f *Firewall) OpenMessage(r io.Reader, source peer.ID, from peer.ID) *message.Message {
	sourcePeer := f.peerSet.MustGetPeer(source)
	fromPeer := f.peerSet.MustGetPeer(from)
	if from != source {
		if f.peerIsBanned(sourcePeer) {
			f.logger.Warn("firewall: source peer is banned", "source", util.FingerprintPeerID(source))
			// If there is any connection to the source peer, close it
			f.closeConnection(source)
			return nil
		}
	}

	if f.peerIsBanned(fromPeer) {
		f.logger.Warn("firewall: from peer banned", "from", util.FingerprintPeerID(from))
		f.closeConnection(from)
		return nil
	}

	msg, err := f.decodeMessage(r, sourcePeer)
	if err != nil {
		f.logger.Debug("unable to decode the message", "from", util.FingerprintPeerID(from), "err", err)
		return nil
	}

	if err := f.checkMessage(msg, sourcePeer); err != nil {
		f.logger.Warn("firewall: invalid message", "err", err, "msg", msg, "from", util.FingerprintPeerID(from))
		f.closeConnection(from)
		return nil
	}

	return msg
}

func (f *Firewall) decodeMessage(r io.Reader, source *peerset.Peer) (*message.Message, error) {
	source.IncreaseReceivedMessage()
	msg := new(message.Message)
	bytesRead, err := msg.Decode(r)
	source.IncreaseReceivedBytes(bytesRead)
	if err != nil {
		source.IncreaseInvalidMessage()
		return nil, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return msg, nil
}

func (f *Firewall) checkMessage(msg *message.Message, source *peerset.Peer) error {
	if err := msg.SanityCheck(); err != nil {
		source.IncreaseInvalidMessage()
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	if msg.Initiator != source.PeerID() {
		source.IncreaseInvalidMessage()
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
