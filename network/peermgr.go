package network

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

const PeerStorePath = "peers.json"

type peerInfo struct {
	AddrInfo  lp2ppeer.AddrInfo
	Direction lp2pnet.Direction
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

	if util.PathExists(PeerStorePath) {
		err := mgr.LoadPeerStore()
		if err != nil {
			mgr.logger.Error("failed to load peer store", "err", err)
		}
		mgr.logger.Info("peer store loaded successfully")
	}

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
	err := mgr.SavePeerStore()
	if err != nil {
		mgr.logger.Error("can't save peer store", "err", err)
	}
}

func (mgr *peerMgr) NumOfConnected() int {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	return len(mgr.peers) // TODO: try to keep record of all peers + connected peers
}

func (mgr *peerMgr) AddPeer(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction,
) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.addPeer(pid, ma, direction)
}

func (mgr *peerMgr) addPeer(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction,
) {
	_, exists := mgr.peers[pid]
	if exists {
		return
	}

	switch direction {
	case lp2pnet.DirInbound:
		mgr.numInbound++

	case lp2pnet.DirOutbound:
		mgr.numOutbound++

	case lp2pnet.DirUnknown:
		//
	}

	mgr.peers[pid] = &peerInfo{
		AddrInfo:  ma,
		Direction: direction,
	}
}

func (mgr *peerMgr) RemovePeer(pid lp2ppeer.ID) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.removePeer(pid)
}

func (mgr *peerMgr) removePeer(pid lp2ppeer.ID) {
	peerInfo, exists := mgr.peers[pid]
	if !exists {
		mgr.logger.Warn("unable to find a peer", "pid", pid)

		return
	}

	switch peerInfo.Direction {
	case lp2pnet.DirInbound:
		mgr.numInbound--

	case lp2pnet.DirOutbound:
		mgr.numOutbound--

	case lp2pnet.DirUnknown:
		//
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

	return peer.AddrInfo.Addrs[0]
}

// checkConnectivity performs the actual work of maintaining connections.
// It ensures that the number of connections stays within the minimum and maximum thresholds.
func (mgr *peerMgr) CheckConnectivity() {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	net := mgr.host.Network()

	// Let's check if some peers are disconnected
	var connectedPeers []lp2ppeer.ID
	for pid := range mgr.peers {
		connectedness := net.Connectedness(pid)
		if connectedness == lp2pnet.Connected {
			connectedPeers = append(connectedPeers, pid)
		} else {
			mgr.removePeer(pid)
		}
	}

	mgr.logger.Debug("check connectivity",
		"peers", len(connectedPeers),
		"inbound", mgr.numInbound,
		"outbound", mgr.numOutbound)

	switch {
	case len(connectedPeers) > mgr.maxConns:
		mgr.logger.Debug("peer count is about maximum threshold",
			"count", connectedPeers,
			"max", mgr.maxConns)

		return

	case len(connectedPeers) < mgr.minConns:
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
			if HasPID(connectedPeers, ai.ID) {
				mgr.logger.Trace("already connected", "peer", ai.String())

				continue
			}

			ConnectAsync(mgr.ctx, mgr.host, ai, mgr.logger)
		}

		for id, pi := range mgr.peers {
			// preventing self dialing.
			if id == mgr.host.ID() {
				continue
			}

			// Don't try to connect to an already connected peer.
			if HasPID(connectedPeers, id) {
				mgr.logger.Trace("already connected", "peer", pi)

				continue
			}

			ConnectAsync(mgr.ctx, mgr.host, pi.AddrInfo, mgr.logger)
		}
	}
}

type PeerStore struct {
	MultiAddr string `json:"multi_address"`
	Direction int    `json:"direction"`
}

func (mgr *peerMgr) SavePeerStore() error {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	ps := make(map[string]*PeerStore)
	for id, info := range mgr.peers {
		ps[id.String()] = &PeerStore{
			MultiAddr: info.AddrInfo.String(),
			Direction: int(info.Direction),
		}
	}

	data, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	return util.WriteFile(PeerStorePath, data)
}

func (mgr *peerMgr) LoadPeerStore() error {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	data, err := util.ReadFile(PeerStorePath)
	if err != nil {
		return err
	}

	ps := make(map[string]*PeerStore)

	err = json.Unmarshal(data, &ps)
	if err != nil {
		return err
	}

	for id, info := range ps {
		id, err := lp2ppeer.Decode(id)
		if err != nil {
			continue
		}
		addr, err := multiaddr.NewMultiaddr(info.MultiAddr)
		if err != nil {
			continue
		}

		ai, err := lp2ppeer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			continue
		}

		mgr.AddPeer(id, *ai, lp2pnet.Direction(info.Direction))
	}

	return nil
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
