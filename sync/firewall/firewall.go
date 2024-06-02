package firewall

import (
	"bytes"
	"io"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/ipblocker"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/ratelimit"
)

// Firewall check packets before passing them to sync module.
type Firewall struct {
	config               *Config
	network              network.Network
	peerSet              *peerset.PeerSet
	blockRateLimit       *ratelimit.RateLimit
	transactionRateLimit *ratelimit.RateLimit
	consensusRateLimit   *ratelimit.RateLimit
	state                state.Facade
	ipBlocker *ipblocker.IPBlocker
	logger               *logger.SubLogger
}

func NewFirewall(conf *Config, net network.Network, peerSet *peerset.PeerSet, st state.Facade,
	log *logger.SubLogger,
) (*Firewall, error) {
	if err := conf.LoadDefaultBlackListAddresses(); err != nil {
		return nil, err
	}

	blocker, err := ipblocker.New(conf.BlackListAddresses)
	if err != nil {
		return nil, err
	}

	blockRateLimit := ratelimit.NewRateLimit(conf.RateLimit.BlockTopic, time.Second)
	transactionRateLimit := ratelimit.NewRateLimit(conf.RateLimit.TransactionTopic, time.Second)
	consensusRateLimit := ratelimit.NewRateLimit(conf.RateLimit.ConsensusTopic, time.Second)

	return &Firewall{
		config:               conf,
		network:              net,
		peerSet:              peerSet,
		blockRateLimit:       blockRateLimit,
		transactionRateLimit: transactionRateLimit,
		consensusRateLimit:   consensusRateLimit,
		ipBlocker: blocker,
		state:                st,
		logger:               log,
	}
}

func (f *Firewall) OpenGossipBundle(data []byte, from peer.ID) *bundle.Bundle {
	bdl, err := f.openBundle(bytes.NewReader(data), from)
	if err != nil {
		f.logger.Debug("firewall: unable to open a gossip bundle",
			"error", err, "bundle", bdl, "from", from)

		return nil
	}

	// TODO: check if gossip flag is set
	// TODO: check if bundle is a gossip bundle

	return bdl
}

func (f *Firewall) IsBlackListAddress(remoteAddr string) bool {
	ip, err := f.parseP2PAddr(remoteAddr)
	if err != nil {
		f.logger.Debug("firewall: unable to parse remote address", "err", err)
		return false
	}

	return f.ipBlocker.IsBlocked(ip)
}

func (f *Firewall) OpenStreamBundle(r io.Reader, from peer.ID) *bundle.Bundle {
	bdl, err := f.openBundle(r, from)
	if err != nil {
		f.logger.Debug("firewall: unable to open a stream bundle",
			"error", err, "bundle", bdl, "from", from)

		return nil
	}

	// TODO: check if gossip flag is NOT set
	// TODO: check if bundle is a stream bundle

	return bdl
}

func (f *Firewall) openBundle(r io.Reader, from peer.ID) (*bundle.Bundle, error) {
	f.peerSet.UpdateLastReceived(from)
	f.peerSet.IncreaseReceivedBundlesCounter(from)

	if f.isPeerBanned(from) {
		f.closeConnection(from)

		return nil, errors.Errorf(errors.ErrInvalidMessage, "peer is banned: %s", from)
	}

	bdl, err := f.decodeBundle(r, from)
	if err != nil {
		f.peerSet.IncreaseInvalidBundlesCounter(from)

		return nil, err
	}

	if err := f.checkBundle(bdl); err != nil {
		f.peerSet.IncreaseInvalidBundlesCounter(from)

		return bdl, err
	}

	return bdl, nil
}

func (f *Firewall) decodeBundle(r io.Reader, pid peer.ID) (*bundle.Bundle, error) {
	bdl := new(bundle.Bundle)
	bytesRead, err := bdl.Decode(r)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}
	f.peerSet.IncreaseReceivedBytesCounter(pid, bdl.Message.Type(), int64(bytesRead))

	return bdl, nil
}

func (f *Firewall) checkBundle(bdl *bundle.Bundle) error {
	if err := bdl.BasicCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
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

	case genesis.Localnet:
		if bdl.Flags&0x3 != 0 {
			return errors.Errorf(errors.ErrInvalidMessage,
				"bundle is not for the localnet")
		}
	}

	return nil
}

func (f *Firewall) isPeerBanned(pid peer.ID) bool {
	p := f.peerSet.GetPeer(pid)

	return p.IsBanned()
}

func (f *Firewall) closeConnection(pid peer.ID) {
	f.network.CloseConnection(pid)
}

func (f *Firewall) AllowBlockRequest() bool {
	return f.blockRateLimit.AllowRequest()
}

func (f *Firewall) AllowTransactionRequest() bool {
	return f.transactionRateLimit.AllowRequest()
}

func (f *Firewall) AllowConsensusRequest() bool {
	return f.consensusRateLimit.AllowRequest()
}

func (f *Firewall) parseP2PAddr(address string) (string, error) {
	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return "", err
	}

	components := addr.Protocols()

	var ip string
	for _, comp := range components {
		switch comp.Name {
		// TODO: can parse dns address and find ip??
		case "ip4", "ip6":
			ipComponent, err := addr.ValueForProtocol(comp.Code)
			if err != nil {
				return "", err
			}
			ip = ipComponent
		}
	}

	return ip, nil
}
