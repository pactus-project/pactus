package firewall

import (
	"bytes"
	"io"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util/ipblocker"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/ratelimit"
)

// Firewall check packets before passing them to sync module.
type Firewall struct {
	config               *Config
	network              network.Network
	peerSet              *peerset.PeerSet
	state                state.Facade
	ipBlocker            *ipblocker.IPBlocker
	blockRateLimit       *ratelimit.RateLimit
	transactionRateLimit *ratelimit.RateLimit
	consensusRateLimit   *ratelimit.RateLimit
	logger               *logger.SubLogger
}

func NewFirewall(conf *Config, network network.Network, peerSet *peerset.PeerSet, state state.Facade,
	log *logger.SubLogger,
) (*Firewall, error) {
	blocker, err := ipblocker.New(conf.BannedNets)
	if err != nil {
		return nil, err
	}

	blockRateLimit := ratelimit.NewRateLimit(conf.RateLimit.BlockTopic, time.Second)
	transactionRateLimit := ratelimit.NewRateLimit(conf.RateLimit.TransactionTopic, time.Second)
	consensusRateLimit := ratelimit.NewRateLimit(conf.RateLimit.ConsensusTopic, time.Second)

	return &Firewall{
		config:               conf,
		network:              network,
		peerSet:              peerSet,
		state:                state,
		ipBlocker:            blocker,
		blockRateLimit:       blockRateLimit,
		transactionRateLimit: transactionRateLimit,
		consensusRateLimit:   consensusRateLimit,
		logger:               log,
	}, nil
}

func (f *Firewall) OpenGossipBundle(data []byte, from peer.ID) (*bundle.Bundle, error) {
	bdl, err := f.openBundle(bytes.NewReader(data), from)
	if err != nil {
		return nil, err
	}

	if !bdl.Message.ShouldBroadcast() {
		f.logger.Warn("firewall: receive stream message as gossip message",
			"error", err, "bundle", bdl, "from", from)

		f.closeConnection(from)

		return nil, ErrGossipMessage
	}

	return bdl, nil
}

// IsBannedAddress checks if the remote IP address is banned.
func (f *Firewall) IsBannedAddress(remoteAddr string) bool {
	ip, err := f.getIPFromMultiAddress(remoteAddr)
	if err != nil {
		f.logger.Warn("firewall: unable to parse remote address", "error", err, "addr", remoteAddr)

		return false
	}

	return f.ipBlocker.IsBanned(ip)
}

func (f *Firewall) OpenStreamBundle(r io.Reader, from peer.ID) (*bundle.Bundle, error) {
	bdl, err := f.openBundle(r, from)
	if err != nil {
		f.logger.Debug("firewall: unable to open a stream bundle",
			"error", err, "bundle", bdl, "from", from)

		return nil, err
	}

	if bdl.Message.ShouldBroadcast() {
		f.logger.Warn("firewall: receive gossip message as stream message",
			"error", err, "bundle", bdl, "from", from)

		f.closeConnection(from)

		return nil, ErrStreamMessage
	}

	return bdl, nil
}

func (f *Firewall) openBundle(r io.Reader, from peer.ID) (*bundle.Bundle, error) {
	f.peerSet.UpdateLastReceived(from)

	peer := f.peerSet.GetPeer(from)
	if peer.Status.IsBanned() {
		f.closeConnection(from)

		return nil, PeerBannedError{
			PeerID:  peer.PeerID,
			Address: peer.Address,
		}
	}

	if f.IsBannedAddress(peer.Address) {
		f.closeConnection(from)
		f.peerSet.UpdateStatus(from, status.StatusBanned)

		return nil, PeerBannedError{
			PeerID:  peer.PeerID,
			Address: peer.Address,
		}
	}

	bdl, bytesRead, err := f.decodeBundle(r)
	if err != nil {
		f.peerSet.UpdateInvalidMetric(from, int64(bytesRead))

		return nil, err
	}

	if err := f.checkBundle(bdl); err != nil {
		f.peerSet.UpdateInvalidMetric(from, int64(bytesRead))

		return bdl, err
	}

	f.peerSet.UpdateReceivedMetric(from, bdl.Message.Type(), int64(bytesRead))

	return bdl, nil
}

func (*Firewall) decodeBundle(r io.Reader) (*bundle.Bundle, int, error) {
	bdl := new(bundle.Bundle)
	bytesRead, err := bdl.Decode(r)
	if err != nil {
		return nil, bytesRead, err
	}

	return bdl, bytesRead, nil
}

func (f *Firewall) checkBundle(bdl *bundle.Bundle) error {
	if err := bdl.BasicCheck(); err != nil {
		return err
	}

	switch f.state.Genesis().ChainType() {
	case genesis.Mainnet:
		if bdl.Flags&0x3 != bundle.BundleFlagNetworkMainnet {
			return ErrNetworkMismatch
		}

	case genesis.Testnet:
		if bdl.Flags&0x3 != bundle.BundleFlagNetworkTestnet {
			return ErrNetworkMismatch
		}

	case genesis.Localnet:
		if bdl.Flags&0x3 != 0 {
			return ErrNetworkMismatch
		}
	}

	return nil
}

func (f *Firewall) closeConnection(pid peer.ID) {
	f.network.CloseConnection(pid)
}

func (*Firewall) getIPFromMultiAddress(address string) (string, error) {
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

func (f *Firewall) AllowBlockRequest() bool {
	return f.blockRateLimit.AllowRequest()
}

func (f *Firewall) AllowTransactionRequest() bool {
	return f.transactionRateLimit.AllowRequest()
}

func (f *Firewall) AllowConsensusRequest() bool {
	return f.consensusRateLimit.AllowRequest()
}
