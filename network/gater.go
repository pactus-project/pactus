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

	filters    *multiaddr.Filters
	peerMgr    *peerMgr
	connsLimit int
	logger     *logger.SubLogger
}

func NewConnectionGater(conf *Config, log *logger.SubLogger) (*ConnectionGater, error) {
	filters := multiaddr.NewFilters()
	if !conf.ForcePrivateNetwork {
		privateSubnets := PrivateSubnets()
		filters = SubnetsToFilters(privateSubnets, multiaddr.ActionDeny)
	}

	connsLimit := conf.ScaledMaxConns() + conf.ConnsThreshold()
	log.Info("connection gater created", "connsLimit", connsLimit)
	return &ConnectionGater{
		filters:    filters,
		connsLimit: connsLimit,
		logger:     log,
	}, nil
}

func (g *ConnectionGater) SetPeerManager(peerMgr *peerMgr) {
	g.lk.Lock()
	defer g.lk.Unlock()

	g.peerMgr = peerMgr
}

func (g *ConnectionGater) onConnectionLimit() bool {
	if g.peerMgr == nil {
		return false
	}

	return g.peerMgr.NumOfConnected() > g.connsLimit
}

func (g *ConnectionGater) InterceptPeerDial(pid lp2ppeer.ID) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onConnectionLimit() {
		g.logger.Debug("InterceptPeerDial rejected: many connections", "pid", pid)
		return false
	}

	return true
}

func (g *ConnectionGater) InterceptAddrDial(pid lp2ppeer.ID, ma multiaddr.Multiaddr) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onConnectionLimit() {
		g.logger.Debug("InterceptAddrDial rejected: many connections", "pid", pid, "ma", ma.String())
		return false
	}

	deny := g.filters.AddrBlocked(ma)
	if deny {
		g.logger.Debug("InterceptAddrDial rejected", "pid", pid, "ma", ma.String())
		return false
	}

	return true
}

func (g *ConnectionGater) InterceptAccept(cma lp2pnetwork.ConnMultiaddrs) bool {
	g.lk.RLock()
	defer g.lk.RUnlock()

	if g.onConnectionLimit() {
		g.logger.Debug("InterceptAccept rejected: many connections")
		return false
	}

	deny := g.filters.AddrBlocked(cma.RemoteMultiaddr())
	if deny {
		g.logger.Debug("InterceptAccept rejected")
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
