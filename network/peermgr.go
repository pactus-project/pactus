package network

import (
	"context"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

type peerInfo struct {
	MultiAddress multiaddr.Multiaddr
	Direction    lp2pnet.Direction
}

// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
type peerMgr struct {
	lk sync.RWMutex

	ctx            context.Context
	bootstrapAddrs []lp2ppeer.AddrInfo
	minConns       int
	maxConns       int
	numInbound     int
	numOutbound    int
	host           lp2phost.Host
	peers          map[lp2ppeer.ID]*peerInfo
	logger         *logger.SubLogger
}

// newPeerMgr creates a new Peer Manager instance.
func newPeerMgr(ctx context.Context, h lp2phost.Host,
	conf *Config, log *logger.SubLogger,
) *peerMgr {
	b := &peerMgr{
		ctx:            ctx,
		bootstrapAddrs: conf.BootstrapAddrInfos(),
		minConns:       conf.ScaledMinConns(),
		maxConns:       conf.ScaledMaxConns(),
		peers:          make(map[lp2ppeer.ID]*peerInfo),
		host:           h,
		logger:         log,
	}

	return b
}

// Start starts the Peer  Manager.
func (mgr *peerMgr) Start() {
	mgr.CheckConnectivity()

	go func() {
		ticker := time.NewTicker(20 * time.Second)
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

func (mgr *peerMgr) Stop() {
}

func (mgr *peerMgr) AddPeer(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction,
) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	switch direction {
	case lp2pnet.DirInbound:
		mgr.numInbound++

	case lp2pnet.DirOutbound:
		mgr.numOutbound++
	}

	mgr.peers[pid] = &peerInfo{
		MultiAddress: ma,
		Direction:    direction,
	}
}

func (mgr *peerMgr) RemovePeer(pid lp2ppeer.ID) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	peer, ok := mgr.peers[pid]
	if !ok {
		mgr.logger.Warn("unable to find a peer", "pid", pid)

		return
	}

	switch peer.Direction {
	case lp2pnet.DirInbound:
		mgr.numInbound--

	case lp2pnet.DirOutbound:
		mgr.numOutbound--
	}

	delete(mgr.peers, pid)
}

func (mgr *peerMgr) GetMultiAddr(pid lp2ppeer.ID) multiaddr.Multiaddr {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	peer := mgr.peers[pid]
	if peer == nil {
		return nil
	}

	return peer.MultiAddress
}

// checkConnectivity performs the actual work of maintaining connections.
// It ensures that the number of connections stays within the minimum and maximum thresholds.
func (mgr *peerMgr) CheckConnectivity() {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	connectedPeers := len(mgr.peers)
	mgr.logger.Debug("check connectivity", "peers", connectedPeers)

	switch {
	case connectedPeers > mgr.maxConns:
		mgr.logger.Debug("peer count is about maximum threshold",
			"count", connectedPeers,
			"max", mgr.maxConns)

		return

	case connectedPeers < mgr.minConns:
		mgr.logger.Info("peer count is less than minimum threshold",
			"count", connectedPeers,
			"min", mgr.minConns)

		for _, ai := range mgr.bootstrapAddrs {
			// preventing self dialing.
			if ai.ID == mgr.host.ID() {
				continue
			}

			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", ai.String())

			// Don't try to connect to an already connected peer.
			if mgr.host.Network().Connectedness(ai.ID) == lp2pnet.Connected {
				mgr.logger.Trace("already connected", "peer", ai.String())

				continue
			}

			ConnectAsync(mgr.ctx, mgr.host, ai, mgr.logger)
		}
	}
}

func (mgr *peerMgr) NumInbound() int {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	return mgr.numInbound
}

func (mgr *peerMgr) NumOutbound() int {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	return mgr.numOutbound
}
