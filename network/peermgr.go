package network

import (
	"context"
	"time"

	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pswarm "github.com/libp2p/go-libp2p/p2p/net/swarm"
	"github.com/pactus-project/pactus/util/logger"
)

// peerMgr attempts to keep the p2p host connected to the network
// by keeping a minimum threshold of connections. If the threshold isn't met it
// connects to a random subset of the peerMgr peers. It does not use peer routing
// to discover new peers. To stop a peerMgr cancel the context passed in Start()
// or call Stop().
type peerMgr struct {
	ctx            context.Context
	bootstrapAddrs []lp2ppeer.AddrInfo
	minConns       int
	maxConns       int

	// Dependencies
	host   lp2phost.Host
	dialer lp2pnet.Dialer
	dht    *lp2pdht.IpfsDHT

	logger *logger.SubLogger
}

// newPeerMgr creates a new Peer Manager instance.
// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
func newPeerMgr(ctx context.Context, h lp2phost.Host, dht *lp2pdht.IpfsDHT,
	conf *Config, log *logger.SubLogger,
) *peerMgr {
	b := &peerMgr{
		ctx:            ctx,
		bootstrapAddrs: conf.BootstrapAddrInfos(),
		minConns:       conf.MinConns,
		maxConns:       conf.MaxConns,
		host:           h,
		dialer:         h.Network(),
		dht:            dht,
		logger:         log,
	}

	return b
}

// Start starts the Peer  Manager.
func (mgr *peerMgr) Start() {
	mgr.checkConnectivity()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-mgr.ctx.Done():
				return
			case <-ticker.C:
				mgr.checkConnectivity()
			}
		}
	}()
}

// Stop stops the Bootstrap.
func (mgr *peerMgr) Stop() {
	// TODO: complete me
}

// checkConnectivity performs the actual work of maintaining connections.
// It ensures that the number of connections stays within the minimum and maximum thresholds.
func (mgr *peerMgr) checkConnectivity() {
	currentPeers := mgr.dialer.Peers()
	mgr.logger.Debug("check connectivity", "peers", len(currentPeers))

	// Let's check if some peers are disconnected
	var connectedPeers []lp2ppeer.ID
	for _, p := range currentPeers {
		connectedness := mgr.dialer.Connectedness(p)
		if connectedness == lp2pnet.Connected {
			connectedPeers = append(connectedPeers, p)
		} else {
			mgr.logger.Debug("peer is not connected to us", "peer", p)
		}
	}

	if len(connectedPeers) > mgr.maxConns {
		mgr.logger.Debug("peer count is about maximum threshold",
			"count", len(connectedPeers),
			"max", mgr.maxConns)
		return
	}

	if len(connectedPeers) < mgr.minConns {
		mgr.logger.Info("peer count is less than minimum threshold",
			"count", len(connectedPeers),
			"min", mgr.minConns)

		for _, pi := range mgr.bootstrapAddrs {
			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", pi.String())

			// Don't try to connect to an already connected peer.
			if HasPID(connectedPeers, pi.ID) {
				mgr.logger.Trace("already connected", "peer", pi.String())
				continue
			}

			if swarm, ok := mgr.host.Network().(*lp2pswarm.Swarm); ok {
				swarm.Backoff().Clear(pi.ID)
			}

			ConnectAsync(mgr.ctx, mgr.host, pi, mgr.logger)
		}
	}
}
