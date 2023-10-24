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

// newPeerMgr returns a new Bootstrap that will attempt to keep connected
// to the network by connecting to the given bootstrap peers.
func newPeerMgr(ctx context.Context, h lp2phost.Host, d lp2pnet.Dialer, dht *lp2pdht.IpfsDHT,
	bootstrapAddrs []lp2ppeer.AddrInfo, minConns int, maxConns int, logger *logger.SubLogger,
) *peerMgr {
	b := &peerMgr{
		ctx:            ctx,
		bootstrapAddrs: bootstrapAddrs,
		minConns:       minConns,
		maxConns:       maxConns,
		host:           h,
		dialer:         d,
		dht:            dht,
		logger:         logger,
	}

	return b
}

// Start starts the Bootstrap bootstrapping. Cancel `ctx` or call Stop() to stop it.
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
}

// checkConnectivity does the actual work. If the number of connected peers
// has fallen below b.MinPeerThreshold it will attempt to connect to
// a random subset of its bootstrap peers.
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
			mgr.logger.Warn("peer is not connected to us", "peer", p)
		}
	}

	if len(connectedPeers) > mgr.maxConns {
		mgr.logger.Debug("peer count is about maximum threshold",
			"count", len(connectedPeers),
			"max", mgr.maxConns)
		return
	}

	if len(connectedPeers) < mgr.minConns {
		mgr.logger.Debug("peer count is less than minimum threshold",
			"count", len(connectedPeers),
			"min", mgr.minConns)

		for _, pi := range mgr.bootstrapAddrs {
			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", pi.String())

			// Don't try to connect to an already connected peer.
			if hasPID(connectedPeers, pi.ID) {
				mgr.logger.Trace("already connected", "peer", pi.String())
				continue
			}

			if swarm, ok := mgr.host.Network().(*lp2pswarm.Swarm); ok {
				swarm.Backoff().Clear(pi.ID)
			}

			ConnectAsync(mgr.ctx, mgr.host, pi, mgr.logger)
		}

		mgr.logger.Debug("expanding the connections")

		err := mgr.dht.Bootstrap(mgr.ctx)
		if err != nil {
			mgr.logger.Warn("peer discovery may suffer", "error", err)
		}
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
