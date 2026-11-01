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
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/scheduler"
	"github.com/pactus-project/pactus/util"
)

type peerInfo struct {
	MultiAddress multiaddr.Multiaddr
	Connected    bool
	Direction    lp2pnet.Direction
	// connCount tracks how many simultaneous connections this peer has.
	// libp2p calls Connected/Disconnected once per *connection*, and a peer
	// can hold more than one at once (simultaneous dial, TCP + QUIC, a
	// relayed connection replaced by a direct one): the peer is only
	// genuinely disconnected once this reaches zero.
	connCount int
}

// Peer Manager attempts to establish connections with other nodes when the
// number of connections falls below the minimum threshold.
type peerMgr struct {
	lk sync.RWMutex

	ctx           context.Context
	logger        *logger.SubLogger
	minConns      int
	numInbound    int
	numOutbound   int
	host          lp2phost.Host
	peers         map[lp2ppeer.ID]*peerInfo
	peerStorePath string
	checkInterval time.Duration
}

// newPeerMgr creates a new Peer Manager instance.
func newPeerMgr(ctx context.Context, host lp2phost.Host,
	conf *Config, log *logger.SubLogger,
) *peerMgr {
	peerStore, err := loadPeerStore(conf.PeerStorePath)
	if err != nil {
		log.Debug("failed to load peer store", "err", err)
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

	mgr := &peerMgr{
		ctx:           ctx,
		minConns:      conf.MinConns(),
		peers:         peers,
		peerStorePath: conf.PeerStorePath,
		host:          host,
		logger:        log,
		checkInterval: conf.CheckConnectivityInterval,
	}

	log.Info("peer manager created", "minConns", mgr.minConns)

	return mgr
}

// Start starts the Peer  Manager.
func (mgr *peerMgr) Start() {
	mgr.CheckConnectivity()

	scheduler.Every(mgr.checkInterval).Do(mgr.ctx, func(context.Context) { mgr.CheckConnectivity() })
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
		pi.connCount++

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
		connCount:    1,
	}
}

// SetPeerDisconnected reports one of the peer's connections closing. It
// returns true only once the peer has no connection left, i.e. it is
// genuinely disconnected; the caller should act on that (e.g. notify other
// modules) only in that case.
func (mgr *peerMgr) SetPeerDisconnected(pid lp2ppeer.ID) bool {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	return mgr.setPeerDisconnected(pid)
}

func (mgr *peerMgr) setPeerDisconnected(pid lp2ppeer.ID) bool {
	peerInfo, exists := mgr.peers[pid]
	if !exists {
		return false
	}

	if !peerInfo.Connected {
		return false
	}

	peerInfo.connCount--
	if peerInfo.connCount > 0 {
		return false
	}

	switch peerInfo.Direction {
	case lp2pnet.DirInbound:
		mgr.numInbound--

		// An inbound peer's remote address is its ephemeral source port,
		// never useful for dialing it back. Unlike an outbound-dialable
		// address, there is nothing worth remembering once it disconnects:
		// keeping the entry would let any inbound stranger (a peer ID
		// costs nothing to generate) grow this table and peers.json
		// without limit, and get re-dialled by CheckConnectivity.
		delete(mgr.peers, pid)

		return true

	case lp2pnet.DirOutbound:
		mgr.numOutbound--

	case lp2pnet.DirUnknown:
		//
	}

	peerInfo.Connected = false
	peerInfo.Direction = lp2pnet.DirUnknown

	return true
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

		for pid, info := range mgr.peers {
			// preventing self dialing.
			if pid == mgr.host.ID() {
				continue
			}

			mgr.logger.Debug("try connecting to a bootstrap peer", "peer", pid.String())

			// Don't try to connect to an already connected peer.
			if info.Connected {
				continue
			}

			ai := lp2ppeer.AddrInfo{
				ID:    pid,
				Addrs: []multiaddr.Multiaddr{info.MultiAddress},
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

	ps := make([]string, 0, len(mgr.peers))
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

	addrs := make([]string, 0)
	err = json.Unmarshal(data, &addrs)
	if err != nil {
		return peerStore, err
	}

	peerStore, err = MakeAddrInfos(addrs)
	if err != nil {
		return []lp2ppeer.AddrInfo{}, err
	}

	return peerStore, nil
}
