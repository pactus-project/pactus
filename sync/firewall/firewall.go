package firewall

import (
	"bytes"
	"encoding/binary"
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
		return bdl, err
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

func (f *Firewall) isExpiredMessage(msgData []byte) bool {
	msgLen := len(msgData)
	if msgLen < 6 {
		return true
	}

	consensusHeightExtracted := false
	var consensusHeight uint32
	consensusHeightBytes := msgData[msgLen-6:]
	// Check if consensus height is set. Refer to the bundle encoding for more details.
	if consensusHeightBytes[0] == 0x04 && consensusHeightBytes[1] == 0x1a {
		consensusHeight = binary.BigEndian.Uint32(consensusHeightBytes[2:])

		if consensusHeight > 2_900_000 {
			consensusHeightExtracted = true
		}
	}
	if !consensusHeightExtracted {
		// Decoding the message at this level is costly, and we should avoid it.
		// In future versions, this code can be removed.
		// However, at the time of writing this code, we need it to prevent replay attacks.
		bdl := new(bundle.Bundle)
		_, err := bdl.Decode(bytes.NewReader(msgData))
		if err != nil {
			return true
		}

		consensusHeight = bdl.Message.ConsensusHeight()
	}

	// The message is expired, or the consensus height is behind the network's current height.
	// In either case, the message is dropped and won't be propagated.
	if consensusHeight < f.state.LastBlockHeight()-1 {
		f.logger.Warn("firewall: expired message", "message height", consensusHeight, "our height", f.state.LastBlockHeight())

		return true
	}

	return false
}

func (f *Firewall) AllowBlockRequest(gossipMsg *network.GossipMessage) network.PropagationPolicy {
	if f.isExpiredMessage(gossipMsg.Data) {
		return network.Drop
	}

	if !f.blockRateLimit.AllowRequest() {
		return network.DropButConsume
	}

	return network.Propagate
}

func (f *Firewall) AllowTransactionRequest(_ *network.GossipMessage) network.PropagationPolicy {
	if !f.transactionRateLimit.AllowRequest() {
		return network.DropButConsume
	}

	return network.Propagate
}

func (f *Firewall) AllowConsensusRequest(gossipMsg *network.GossipMessage) network.PropagationPolicy {
	if f.isExpiredMessage(gossipMsg.Data) {
		return network.Drop
	}

	if !f.consensusRateLimit.AllowRequest() {
		return network.DropButConsume
	}

	return network.Propagate
}
