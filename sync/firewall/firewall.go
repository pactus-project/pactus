package firewall

import (
	"bytes"
	"io"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

// Firewall check packets before passing them to sync module.
type Firewall struct {
	network network.Network
	state   state.Facade
	config  *Config
	peerSet *peerset.PeerSet
	logger  *logger.Logger
}

func NewFirewall(conf *Config, net network.Network, peerSet *peerset.PeerSet, state state.Facade,
	logger *logger.Logger) *Firewall {
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
		f.peerSet.UpdateLastSeen(from)
		if f.isPeerBanned(from) {
			f.logger.Warn("firewall: from peer banned", "from", from)
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
	f.peerSet.UpdateLastSeen(source)
	f.peerSet.IncreaseReceivedBundlesCounter(source)

	if f.isPeerBanned(source) {
		// If there is any connection to the source peer, close it
		f.closeConnection(source)
		return nil, errors.Errorf(errors.ErrInvalidMessage, "source peer is banned: %s", source)
	}

	bdl, err := f.decodeBundle(r, source)
	if err != nil {
		f.peerSet.IncreaseInvalidBundlesCounter(source)
		return nil, err
	}

	if err := f.checkBundle(bdl, source); err != nil {
		f.peerSet.IncreaseInvalidBundlesCounter(source)
		return nil, err
	}

	return bdl, nil
}

func (f *Firewall) decodeBundle(r io.Reader, pid peer.ID) (*bundle.Bundle, error) {
	bdl := new(bundle.Bundle)
	bytesRead, err := bdl.Decode(r)
	f.peerSet.IncreaseReceivedBytesCounter(pid, bytesRead)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return bdl, nil
}

func (f *Firewall) checkBundle(bdl *bundle.Bundle, pid peer.ID) error {
	if err := bdl.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	if bdl.Initiator != pid {
		return errors.Errorf(errors.ErrInvalidMessage,
			"source is not same as initiator. source: %v, initiator: %v", pid, bdl.Initiator)
	}

	switch f.state.Genesis().ChainType() {
	case genesis.Mainnet:
		if bdl.Flags&0x3 != bundle.BundleFlagNetworkMainnet {
			return errors.Errorf(errors.ErrInvalidMessage,
				"bundle is not for the mainnet")
		}

	case genesis.Testnet:
		if bdl.Flags&0x3 != bundle.BundleFlagNetworkTestnet {
			return errors.Errorf(errors.ErrInvalidMessage,
				"bundle is not for the testnet")
		}

	default:
		if bdl.Flags&0x3 != 0 {
			return errors.Errorf(errors.ErrInvalidMessage,
				"bundle is not for the localnet")
		}
	}

	return nil
}

func (f *Firewall) isPeerBanned(pid peer.ID) bool {
	if !f.config.Enabled {
		return false
	}

	p := f.peerSet.GetPeer(pid)
	return p.IsBanned()
}

func (f *Firewall) closeConnection(pid peer.ID) {
	if !f.config.Enabled {
		return
	}

	f.network.CloseConnection(pid)
}
