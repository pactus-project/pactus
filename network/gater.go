package network

import (
	"sync"

	lp2pconnmgr "github.com/libp2p/go-libp2p/core/connmgr"
	lp2pcontrol "github.com/libp2p/go-libp2p/core/control"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

var _ lp2pconnmgr.ConnectionGater = &ConnectionGater{}

type ConnectionGater struct {
	lk sync.RWMutex

	filters     *multiaddr.Filters
	peerMgr     *peerMgr
	acceptLimit int
	dialLimit   int
	logger      *logger.SubLogger
}

func NewConnectionGater(conf *Config, log *logger.SubLogger) (*ConnectionGater, error) {
	filters := multiaddr.NewFilters()
	if !conf.ForcePrivateNetwork {
		privateSubnets := PrivateSubnets()
		filters = SubnetsToFilters(privateSubnets, multiaddr.ActionDeny)
	}

	acceptLimit := conf.ScaledMaxConns()
	dialLimit := conf.ScaledMaxConns() / 4
	log.Info("connection gater created", "listen", acceptLimit, "dial", dialLimit)

	return &ConnectionGater{
		filters:     filters,
		acceptLimit: acceptLimit,
		dialLimit:   dialLimit,
		logger:      log,
	}, nil
}

func (g *ConnectionGater) SetPeerManager(peerMgr *peerMgr) {
	g.lk.Lock()
	defer g.lk.Unlock()

	g.peerMgr = peerMgr
}

func (g *ConnectionGater) onDialLimit() bool {
	if g.peerMgr == nil {
		return false
	}

	return g.peerMgr.NumOutbound() > g.dialLimit
}

func (g *ConnectionGater) onAcceptLimit() bool {
	if g.peerMgr == nil {
		return false
	}

	return g.peerMgr.NumInbound() > g.acceptLimit
}

func (g *ConnectionGater) InterceptPeerDial(pid lp2ppeer.ID) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onDialLimit() {
		g.logger.Info("InterceptPeerDial rejected: many connections",
			"pid", pid, "outbound", g.peerMgr.NumOutbound())

		return false
	}

	return true
}

func (g *ConnectionGater) InterceptAddrDial(pid lp2ppeer.ID, ma multiaddr.Multiaddr) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onDialLimit() {
		g.logger.Info("InterceptAddrDial rejected: many connections",
			"pid", pid, "ma", ma.String(), "outbound", g.peerMgr.NumOutbound())

		return false
	}

	deny := g.filters.AddrBlocked(ma)
	if deny {
		g.logger.Info("InterceptAddrDial rejected", "pid", pid, "ma", ma.String())

		return false
	}

	return true
}

func (g *ConnectionGater) InterceptAccept(cma lp2pnetwork.ConnMultiaddrs) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onAcceptLimit() {
		g.logger.Info("InterceptAccept rejected: many connections",
			"inbound", g.peerMgr.NumInbound())

		return false
	}

	deny := g.filters.AddrBlocked(cma.RemoteMultiaddr())
	if deny {
		g.logger.Info("InterceptAccept rejected")

		return false
	}

	return true
}

func (g *ConnectionGater) InterceptSecured(_ lp2pnetwork.Direction, _ lp2ppeer.ID, _ lp2pnetwork.ConnMultiaddrs) bool {
	return true
}

func (g *ConnectionGater) InterceptUpgraded(_ lp2pnetwork.Conn) (bool, lp2pcontrol.DisconnectReason) {
	return true, 0
}
