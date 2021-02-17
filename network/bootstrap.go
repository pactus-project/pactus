package network

import (
	"context"
	"sync"
	"time"

	host "github.com/libp2p/go-libp2p-core/host"
	inet "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/zarbchain/zarb-go/logger"
)

// Bootstrapper attempts to keep the p2p host connected to the network
// by keeping a minimum threshold of connections. If the threshold isn't met it
// connects to a random subset of the bootstrap peers. It does not use peer routing
// to discover new peers. To stop a Bootstrapper cancel the context passed in Start()
// or call Stop().
type Bootstrapper struct {
	config *BootstrapConfig

	bootstrapPeers []peer.AddrInfo

	// Dependencies
	h host.Host
	d inet.Dialer
	r routing.Routing

	// Bookkeeping
	ticker         *time.Ticker
	ctx            context.Context
	cancel         context.CancelFunc
	dhtBootStarted bool

	logger *logger.Logger
}

// NewBootstrapper returns a new Bootstrapper that will attempt to keep connected
// to the network by connecting to the given bootstrap peers.
func NewBootstrapper(ctx context.Context, h host.Host, d inet.Dialer, r routing.Routing, conf *BootstrapConfig, logger *logger.Logger) *Bootstrapper {
	b := &Bootstrapper{
		ctx:    ctx,
		config: conf,
		h:      h,
		d:      d,
		r:      r,
		logger: logger,
	}

	addresses, err := PeerAddrsToAddrInfo(conf.Addresses)
	if err != nil {
		b.logger.Panic("couldn't parse bootstrap addresses", "addressed", conf.Addresses)
	}

	b.bootstrapPeers = addresses
	b.checkConnectivity()

	return b
}

// Start starts the Bootstrapper bootstrapping. Cancel `ctx` or call Stop() to stop it.
func (b *Bootstrapper) Start() {
	b.ctx, b.cancel = context.WithCancel(b.ctx)
	b.ticker = time.NewTicker(b.config.Period)

	go func() {
		defer b.ticker.Stop()

		for {
			select {
			case <-b.ctx.Done():
				return
			case <-b.ticker.C:
				b.checkConnectivity()
			}
		}
	}()
}

// Stop stops the Bootstrapper.
func (b *Bootstrapper) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
}

// checkConnectivity does the actual work. If the number of connected peers
// has fallen below b.MinPeerThreshold it will attempt to connect to
// a random subset of its bootstrap peers.
func (b *Bootstrapper) checkConnectivity() {
	currentPeers := b.d.Peers()
	b.logger.Debug("Check connectivity", "peers", len(currentPeers), "threshold", b.config.MinThreshold, "timeout", b.config.Timeout)

	peersNeeded := b.config.MinThreshold - len(currentPeers)
	if peersNeeded < 1 {
		return
	}

	ctx, cancel := context.WithTimeout(b.ctx, b.config.Timeout)
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		// After connecting to bootstrap peers, bootstrap the DHT.
		// DHT Bootstrap is a persistent process so only do this once.
		if !b.dhtBootStarted {
			b.dhtBootStarted = true
			err := b.bootstrapIpfsRouting()
			if err != nil {
				b.logger.Warn("Peer discovery may suffer.", "err", err)
			}
		}
		cancel()
	}()

	for _, pinfo := range b.bootstrapPeers {
		b.logger.Trace("Try connecting to a bootstrap peer.", "peer", pinfo.String())

		// Don't try to connect to an already connected peer.
		if hasPID(currentPeers, pinfo.ID) {
			b.logger.Trace("Already connected.", "peer", pinfo.String())
			continue
		}

		wg.Add(1)
		go func(pi peer.AddrInfo) {
			if err := b.h.Connect(ctx, pi); err != nil {
				b.logger.Error("got error trying to connect to bootstrap node ", "info", pi, "err", err.Error())
			}
			wg.Done()
		}(pinfo)
	}
}

func hasPID(pids []peer.ID, pid peer.ID) bool {
	for _, p := range pids {
		if p == pid {
			return true
		}
	}
	return false
}

func (b *Bootstrapper) bootstrapIpfsRouting() error {
	dht, ok := b.r.(*dht.IpfsDHT)
	if !ok {
		b.logger.Warn("No bootstrapping to do exit quietly.")
		// No bootstrapping to do exit quietly.
		return nil
	}

	return dht.Bootstrap(b.ctx)
}
