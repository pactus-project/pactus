package network

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

type peerInfo struct {
	MultiAddress multiaddr.Multiaddr
	Connected    bool
	Direction    lp2pnet.Direction
}

// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
type peerMgr struct {
	lk sync.RWMutex

	ctx           context.Context
	minConns      int
	numInbound    int
	numOutbound   int
	host          lp2phost.Host
	peers         map[lp2ppeer.ID]*peerInfo
	peerStorePath string
	logger        *logger.SubLogger
}

// newPeerMgr creates a new Peer Manager instance.
func newPeerMgr(ctx context.Context, h lp2phost.Host,
	conf *Config, log *logger.SubLogger,
) *peerMgr {
	peerStore, err := loadPeerStore(conf.PeerStorePath)
	if err != nil {
		log.Error("failed to load peer store", "err", err)
	}
	log.Info("peer store loaded successfully")

	peerStore = append(peerStore, conf.BootstrapAddrInfos()...)

	peers := make(map[lp2ppeer.ID]*peerInfo)
	for _, ai := range peerStore {
		peers[ai.ID] = &peerInfo{
			MultiAddress: ai.Addrs[0],
			Connected:    false,
			Direction:    lp2pnet.DirUnknown,
		}
	}

	pm := &peerMgr{
		ctx:           ctx,
		minConns:      conf.MinConns(),
		peers:         peers,
		peerStorePath: conf.PeerStorePath,
		host:          h,
		logger:        log,
	}

	log.Info("peer manager created", "minConns", pm.minConns)

	return pm
}

// Start starts the Peer  Manager.
func (mgr *peerMgr) Start() {
	mgr.CheckConnectivity()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
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
	if err := mgr.savePeerStore(); err != nil {
		mgr.logger.Error("can't save peer store", "err", err)
	}
}

func (mgr *peerMgr) SetPeerConnected(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction,
) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.setPeerConnected(pid, ma, direction)
}

func (mgr *peerMgr) setPeerConnected(pid lp2ppeer.ID, ma multiaddr.Multiaddr,
	direction lp2pnet.Direction,
) {
	pi, exists := mgr.peers[pid]
	if exists && pi.Connected {
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
		Connected:    true,
		MultiAddress: ma,
		Direction:    direction,
	}
}

func (mgr *peerMgr) SetPeerDisconnected(pid lp2ppeer.ID) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	mgr.setPeerDisconnected(pid)
}

func (mgr *peerMgr) setPeerDisconnected(pid lp2ppeer.ID) {
	pi, exists := mgr.peers[pid]
	if !exists {
		return
	}

	if !pi.Connected {
		return
	}

	switch pi.Direction {
	case lp2pnet.DirInbound:
		mgr.numInbound--

	case lp2pnet.DirOutbound:
		mgr.numOutbound--

	case lp2pnet.DirUnknown:
		//
	}

	pi.Connected = false
	pi.Direction = lp2pnet.DirUnknown
}

// CheckConnectivity performs the actual work of maintaining connections.
// It ensures that the number of connections stays within the minimum and maximum thresholds.
func (mgr *peerMgr) CheckConnectivity() {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	net := mgr.host.Network()

	// Check if some peers are disconnected
	numConnected := 0
	for pid := range mgr.peers {
		connectedness := net.Connectedness(pid)
		if connectedness == lp2pnet.Connected {
			numConnected++
		}
	}

	mgr.logger.Debug("check connectivity",
		"numConnected", numConnected,
		"inbound", mgr.numInbound,
		"outbound", mgr.numOutbound)

	if numConnected < mgr.minConns {
		mgr.logger.Info("peer count is less than minimum threshold",
			"numConnected", numConnected,
			"min", mgr.minConns)

		for id, pi := range mgr.peers {
			// preventing self dialing.
			if id == mgr.host.ID() {
				continue
			}

			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", id.String())

			// Don't try to connect to an already connected peer.
			if pi.Connected {
				continue
			}

			ai := lp2ppeer.AddrInfo{
				ID:    id,
				Addrs: []multiaddr.Multiaddr{pi.MultiAddress},
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

func (mgr *peerMgr) savePeerStore() error {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	ps := make([]string, 0)
	for id, info := range mgr.peers {
		ps = append(ps, fmt.Sprintf("%s/p2p/%s", info.MultiAddress.String(), id.String()))
	}

	data, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	return util.WriteFile(mgr.peerStorePath, data)
}

func loadPeerStore(path string) ([]lp2ppeer.AddrInfo, error) {
	peerStore := make([]lp2ppeer.AddrInfo, 0)

	data, err := util.ReadFile(path)
	if err != nil {
		return peerStore, err
	}

	ps := make([]string, 0)
	err = json.Unmarshal(data, &ps)
	if err != nil {
		return peerStore, err
	}

	peerStore, err = MakeAddrInfos(ps)
	if err != nil {
		return []lp2ppeer.AddrInfo{}, err
	}

	return peerStore, nil
}
