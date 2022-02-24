package firewall

import (
	"bytes"
	"io"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/bundle"
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

func (f *Firewall) OpenGossipBundle(data []byte, source peer.ID, from peer.ID) *bundle.Bundle {
	if from != source {
		fromPeer := f.peerSet.MustGetPeer(from)
		if f.peerIsBanned(fromPeer) {
			f.logger.Warn("firewall: from peer banned", "from", util.FingerprintPeerID(from))
			f.closeConnection(from)
			return nil
		}
	}

	bdl, err := f.openBundle(bytes.NewReader(data), source)
	if err != nil {
		f.logger.Warn("firewall: unable to open a gossip bundle", "err", err)
		f.closeConnection(from)
		return nil
	}

	// TODO: check if gossip flag is set
	// TODO: check if bundle is a gossip bundle

	return bdl
}

func (f *Firewall) OpenStreamBundle(r io.Reader, from peer.ID) *bundle.Bundle {
	bdl, err := f.openBundle(r, from)
	if err != nil {
		f.logger.Warn("firewall: unable to open a stream bundle", "err", err)
		f.closeConnection(from)
		return nil
	}

	// TODO: check if gossip flag is NOT set
	// TODO: check if bundle is a stream bundle

	return bdl
}

func (f *Firewall) openBundle(r io.Reader, source peer.ID) (*bundle.Bundle, error) {
	peer := f.peerSet.MustGetPeer(source)
	peer.IncreaseReceivedBundlesCounter()

	if f.peerIsBanned(peer) {
		// If there is any connection to the source peer, close it
		f.closeConnection(source)
		return nil, errors.Errorf(errors.ErrInvalidMessage, "Source peer is banned: %s", source)
	}

	bdl, err := f.decodeBundle(r, peer)
	if err != nil {
		peer.IncreaseInvalidBundlesCounter()
		return nil, err
	}

	if err := f.checkBundle(bdl, peer); err != nil {
		peer.IncreaseInvalidBundlesCounter()
		return nil, err
	}

	return bdl, nil
}

func (f *Firewall) decodeBundle(r io.Reader, source *peerset.Peer) (*bundle.Bundle, error) {
	bdl := new(bundle.Bundle)
	bytesRead, err := bdl.Decode(r)
	source.IncreaseReceivedBytesCounter(bytesRead)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return bdl, nil
}

func (f *Firewall) checkBundle(bdl *bundle.Bundle, source *peerset.Peer) error {
	if err := bdl.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	if bdl.Initiator != source.PeerID() {
		return errors.Errorf(errors.ErrInvalidMessage,
			"source is not same as initiator. source: %v, initiator: %v", source.PeerID(), bdl.Initiator)
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
