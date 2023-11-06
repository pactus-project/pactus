package network

import (
	"context"
	"sync"
	"time"

	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pswarm "github.com/libp2p/go-libp2p/p2p/net/swarm"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
	"golang.org/x/exp/slices"
)

type peerInfo struct {
	MultiAddress multiaddr.Multiaddr
	Direction    lp2pnet.Direction
	Protocols    []lp2pcore.ProtocolID
}

// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
type peerMgr struct {
	lk sync.RWMutex

	ctx              context.Context
	bootstrapAddrs   []lp2ppeer.AddrInfo
	minConns         int
	maxConns         int
	host             lp2phost.Host
	dht              *lp2pdht.IpfsDHT
	peers            map[lp2ppeer.ID]*peerInfo
	streamProtocolID lp2pcore.ProtocolID
	logger           *logger.SubLogger
}

// newPeerMgr creates a new Peer Manager instance.
func newPeerMgr(ctx context.Context, h lp2phost.Host, dht *lp2pdht.IpfsDHT,
	streamProtocolID lp2pcore.ProtocolID, conf *Config, log *logger.SubLogger,
) *peerMgr {
	b := &peerMgr{
		ctx:              ctx,
		bootstrapAddrs:   conf.BootstrapAddrInfos(),
		minConns:         conf.MinConns,
		maxConns:         conf.MaxConns,
		streamProtocolID: streamProtocolID,
		peers:            make(map[lp2ppeer.ID]*peerInfo),
		host:             h,
		dht:              dht,
		logger:           log,
	}

	return b
}

// Start starts the Peer  Manager.
func (mgr *peerMgr) Start() {
	mgr.CheckConnectivity()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-mgr.ctx.Done():
				return
			case <-ticker.C:
				mgr.CheckConnectivity()
			}
		}
	}()
}

// Stop stops the Bootstrap.
func (mgr *peerMgr) Stop() {
	// TODO: complete me
}

func (mgr *peerMgr) AddPeer(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction, protocols []lp2pcore.ProtocolID,
) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.peers[pid] = &peerInfo{
		MultiAddress: ma,
		Direction:    direction,
		Protocols:    protocols,
	}
}

func (mgr *peerMgr) RemovePeer(pid lp2ppeer.ID) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	delete(mgr.peers, pid)
}

// checkConnectivity performs the actual work of maintaining connections.
// It ensures that the number of connections stays within the minimum and maximum thresholds.
func (mgr *peerMgr) CheckConnectivity() {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.logger.Debug("check connectivity", "peers", len(mgr.peers))

	net := mgr.host.Network()

	// Make sure we have connected to at least one peer that supports the stream protocol.
	hasStreamConn := 0
	for _, pi := range mgr.peers {
		if slices.Contains(pi.Protocols, mgr.streamProtocolID) {
			hasStreamConn++
		}
	}

	if hasStreamConn == 0 {
		// TODO: is it possible?
		mgr.logger.Warn("no stream connection")

		for pid := range mgr.peers {
			_ = net.ClosePeer(pid)
		}

		time.Sleep(1 * time.Second)
	}

	// Let's check if some peers are disconnected
	var connectedPeers []lp2ppeer.ID
	for pid := range mgr.peers {
		connectedness := net.Connectedness(pid)
		if connectedness == lp2pnet.Connected {
			connectedPeers = append(connectedPeers, pid)
		} else {
			mgr.logger.Debug("peer is not connected to us", "peer", pid)
			delete(mgr.peers, pid)
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
