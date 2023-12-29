package network

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pswarm "github.com/libp2p/go-libp2p/p2p/net/swarm"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

const peerStorePath = "peers.json"

type peerInfo struct {
	MultiAddress multiaddr.Multiaddr `json:"multi_address"`
	Direction    lp2pnet.Direction   `json:"direction"`
}

// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
type peerMgr struct {
	lk sync.RWMutex

	ctx            context.Context
	bootstrapAddrs []lp2ppeer.AddrInfo
	minConns       int
	maxConns       int
	host           lp2phost.Host
	peers          map[lp2ppeer.ID]*peerInfo
	logger         *logger.SubLogger
}

// newPeerMgr creates a new Peer Manager instance.
func newPeerMgr(ctx context.Context, h lp2phost.Host,
	conf *Config, log *logger.SubLogger,
) *peerMgr {
	peers := make(map[lp2ppeer.ID]*peerInfo)
	if util.PathExists(peerStorePath) {
		data, err := util.ReadFile(peerStorePath)
		if err != nil {
			log.Error("can't read peer list", "error", err)
		}
		peers, err = loadPeerStore(data)
		if err != nil {
			log.Error("can't load peer list", "error", err)
		}
	}

	b := &peerMgr{
		ctx:            ctx,
		bootstrapAddrs: conf.BootstrapAddrInfos(),
		minConns:       conf.ScaledMinConns(),
		maxConns:       conf.ScaledMaxConns(),
		peers:          peers,
		host:           h,
		logger:         log,
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
	data, err := json.Marshal(mgr.getPeerStore())
	if err != nil {
		mgr.logger.Error("can't marshal peer list", "error", err)
	}

	err = util.WriteFile(peerStorePath, data)
	if err != nil {
		mgr.logger.Error("can't save peer list", "error", err)
	}
}

func (mgr *peerMgr) NumOfConnected() int {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	return len(mgr.peers)
}

func (mgr *peerMgr) AddPeer(pid lp2ppeer.ID, ma multiaddr.Multiaddr, direction lp2pnet.Direction,
) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.peers[pid] = &peerInfo{
		MultiAddress: ma,
		Direction:    direction,
	}
}

func (mgr *peerMgr) RemovePeer(pid lp2ppeer.ID) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

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

	mgr.logger.Debug("check connectivity", "peers", len(mgr.peers))

	net := mgr.host.Network()

	// Let's check if some peers are disconnected
	var connectedPeers []lp2ppeer.ID
	for pid := range mgr.peers {
		connectedness := net.Connectedness(pid)
		if connectedness == lp2pnet.Connected {
			connectedPeers = append(connectedPeers, pid)
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

		for _, ai := range mgr.bootstrapAddrs {
			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", ai.String())

			// Don't try to connect to an already connected peer.
			if HasPID(connectedPeers, ai.ID) {
				mgr.logger.Trace("already connected", "peer", ai.String())
				continue
			}

			if swarm, ok := mgr.host.Network().(*lp2pswarm.Swarm); ok {
				swarm.Backoff().Clear(ai.ID)
			}

			ConnectAsync(mgr.ctx, mgr.host, ai, mgr.logger)
		}
	}
}

type PeerStore struct {
	Peers map[string]*peerInfo `json:"peers"`
}

func (mgr *peerMgr) getPeerStore() *PeerStore {
	ps := PeerStore{
		Peers: map[string]*peerInfo{},
	}

	for id, info := range mgr.peers {
		ps.Peers[id.String()] = info
	}
	return &ps
}

func loadPeerStore(data []byte) (map[lp2ppeer.ID]*peerInfo, error) {
	peerStorePeers := make(map[string]*peerInfo)

	err := json.Unmarshal(data, &peerStorePeers)
	if err != nil {
		return nil, err
	}

	peers := make(map[lp2ppeer.ID]*peerInfo)
	for id, info := range peerStorePeers {
		peerID, err := lp2ppeer.Decode(id)
		if err != nil {
			continue
		}

		peers[peerID] = info
	}

	return peers, nil
}
