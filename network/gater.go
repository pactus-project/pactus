package network

import (
	lp2pconnmgr "github.com/libp2p/go-libp2p/core/connmgr"
	lp2pcontrol "github.com/libp2p/go-libp2p/core/control"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pconngater "github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

var _ lp2pconnmgr.ConnectionGater = &ConnectionGater{}

type ConnectionGater struct {
	*lp2pconngater.BasicConnectionGater

	host    lp2phost.Host
	maxConn int
	logger  *logger.SubLogger
}

func NewConnectionGater(conf *Config, log *logger.SubLogger) (*ConnectionGater, error) {
	connGater, err := lp2pconngater.NewBasicConnectionGater(nil)
	if err != nil {
		return nil, err
	}

	if !conf.ForcePrivateNetwork {
		privateSubnets := PrivateSubnets()
		for _, sn := range privateSubnets {
			err := connGater.BlockSubnet(sn)
			if err != nil {
				return nil, LibP2PError{Err: err}
			}
		}
	}

	return &ConnectionGater{
		BasicConnectionGater: connGater,
		maxConn:              conf.MaxConns,
		logger:               log,
	}, nil
}

func (g *ConnectionGater) SetHost(host lp2phost.Host) {
	g.host = host
}

func (g *ConnectionGater) hasMaxConnections() bool {
	return len(g.host.Network().Peers()) > g.maxConn
}

func (g *ConnectionGater) InterceptPeerDial(p lp2ppeer.ID) bool {
	if g.hasMaxConnections() {
		g.logger.Debug("InterceptPeerDial rejected: many connections")
		return false
	}

	allow := g.BasicConnectionGater.InterceptPeerDial(p)
	if !allow {
		g.logger.Debug("InterceptPeerDial rejected", "p")
	}

	return allow
}

func (g *ConnectionGater) InterceptAddrDial(p lp2ppeer.ID, ma multiaddr.Multiaddr) bool {
	if g.hasMaxConnections() {
		g.logger.Debug("InterceptAddrDial rejected: many connections")
		return false
	}

	allow := g.BasicConnectionGater.InterceptAddrDial(p, ma)
	if !allow {
		g.logger.Debug("InterceptAddrDial rejected", "p", p, "ma", ma.String())
	}

	return allow
}

func (g *ConnectionGater) InterceptAccept(cma lp2pnetwork.ConnMultiaddrs) bool {
	if g.hasMaxConnections() {
		g.logger.Debug("InterceptAccept rejected: many connections")
		return false
	}

	allow := g.BasicConnectionGater.InterceptAccept(cma)
	if !allow {
		g.logger.Debug("InterceptAccept rejected")
	}

	return allow
}

func (g *ConnectionGater) InterceptSecured(dir lp2pnetwork.Direction, p lp2ppeer.ID,
	cma lp2pnetwork.ConnMultiaddrs,
) bool {
	allow := g.BasicConnectionGater.InterceptSecured(dir, p, cma)
	if !allow {
		g.logger.Debug("InterceptSecured rejected", "p", p)
	}

	return allow
}

func (g *ConnectionGater) InterceptUpgraded(con lp2pnetwork.Conn) (bool, lp2pcontrol.DisconnectReason) {
	return g.BasicConnectionGater.InterceptUpgraded(con)
}
