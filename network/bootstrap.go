package network

import (
	"context"
	"time"

	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2prouting "github.com/libp2p/go-libp2p/core/routing"
	"github.com/pactus-project/pactus/util/logger"
)

// bootstrap attempts to keep the p2p host connected to the network
// by keeping a minimum threshold of connections. If the threshold isn't met it
// connects to a random subset of the bootstrap peers. It does not use peer routing
// to discover new peers. To stop a bootstrap cancel the context passed in Start()
// or call Stop().
type bootstrap struct {
	ctx    context.Context
	config *BootstrapConfig

	bootstrapPeers []lp2ppeer.AddrInfo

	// Dependencies
	host    lp2phost.Host
	dialer  lp2pnet.Dialer
	routing lp2prouting.Routing

	logger *logger.SubLogger
}

// newBootstrap returns a new Bootstrap that will attempt to keep connected
// to the network by connecting to the given bootstrap peers.
func newBootstrap(ctx context.Context, h lp2phost.Host, d lp2pnet.Dialer, r lp2prouting.Routing,
	conf *BootstrapConfig, logger *logger.SubLogger,
) *bootstrap {
	b := &bootstrap{
		ctx:     ctx,
		config:  conf,
		host:    h,
		dialer:  d,
		routing: r,
		logger:  logger,
	}

	addresses, err := PeerAddrsToAddrInfo(conf.Addresses)
	if err != nil {
		b.logger.Panic("couldn't parse bootstrap addresses", "error", err, "addresses", conf.Addresses)
	}
	b.bootstrapPeers = addresses

	return b
}

// Start starts the Bootstrap bootstrapping. Cancel `ctx` or call Stop() to stop it.
func (b *bootstrap) Start() {
	// Protecting bootstrap peers
	for _, a := range b.bootstrapPeers {
		b.host.ConnManager().Protect(a.ID, "bootstrap")
	}

	b.checkConnectivity()

	go func() {
		ticker := time.NewTicker(b.config.Period)
		defer ticker.Stop()

		for {
			select {
			case <-b.ctx.Done():
				return
			case <-ticker.C:
				b.checkConnectivity()
			}
		}
	}()
}

// Stop stops the Bootstrap.
func (b *bootstrap) Stop() {
}

// checkConnectivity does the actual work. If the number of connected peers
// has fallen below b.MinPeerThreshold it will attempt to connect to
// a random subset of its bootstrap peers.
func (b *bootstrap) checkConnectivity() {
	currentPeers := b.dialer.Peers()
	b.logger.Debug("check connectivity", "peers", len(currentPeers))

	// Let's check if some peers are disconnected
	var connectedPeers []lp2ppeer.ID
	for _, p := range currentPeers {
		connectedness := b.dialer.Connectedness(p)
		if connectedness == lp2pnet.Connected {
			connectedPeers = append(connectedPeers, p)
		} else {
			b.logger.Warn("peer is not connected to us", "peer", p)
		}
	}

	if len(connectedPeers) > b.config.MaxThreshold {
		b.logger.Debug("peer count is about maximum threshold",
			"count", len(connectedPeers),
			"threshold", b.config.MaxThreshold)
		return
	}

	if len(connectedPeers) < b.config.MinThreshold {
		b.logger.Debug("peer count is less than minimum threshold",
			"count", len(connectedPeers),
			"threshold", b.config.MinThreshold)

		for _, pi := range b.bootstrapPeers {
			b.logger.Debug("try connecting to a bootstrap peer", "peer", pi.String())

			// Don't try to connect to an already connected peer.
			if hasPID(connectedPeers, pi.ID) {
				b.logger.Trace("already connected", "peer", pi.String())
				continue
			}

			if err := b.host.Connect(b.ctx, pi); err != nil {
				b.logger.Error("error trying to connect to bootstrap node", "info", pi, "error", err)
			}
		}

		b.logger.Debug("expanding the connections")
		b.expand()
	}
}

func hasPID(pids []lp2ppeer.ID, pid lp2ppeer.ID) bool {
	for _, p := range pids {
		if p == pid {
			return true
		}
	}
	return false
}

func (b *bootstrap) expand() {
	dht, ok := b.routing.(*lp2pdht.IpfsDHT)
	if !ok {
		b.logger.Warn("no bootstrapping to do exit quietly.")
		return
	}

	err := dht.Bootstrap(b.ctx)
	if err != nil {
		b.logger.Warn("peer discovery may suffer", "error", err)
	}
}
