package network

import (
	lp2pconnmgr "github.com/libp2p/go-libp2p/core/connmgr"
	lp2pcontrol "github.com/libp2p/go-libp2p/core/control"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	lp2pconngater "github.com/libp2p/go-libp2p/p2p/net/conngater"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

var _ lp2pconnmgr.ConnectionGater = &ConnectionGater{}

type ConnectionGater struct {
	*lp2pconngater.BasicConnectionGater

	logger *logger.SubLogger
}

func NewConnectionGater(conf *Config) (*ConnectionGater, error) {
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
	}, nil
}

func (g *ConnectionGater) SetLogger(log *logger.SubLogger) {
	g.logger = log
}

func (g *ConnectionGater) InterceptPeerDial(p peer.ID) bool {
	allow := g.BasicConnectionGater.InterceptPeerDial(p)
	if !allow {
		g.logger.Debug("InterceptPeerDial not allowed", "p")
	}

	return allow
}

func (g *ConnectionGater) InterceptAddrDial(p peer.ID, ma multiaddr.Multiaddr) bool {
	allow := g.BasicConnectionGater.InterceptAddrDial(p, ma)
	if !allow {
		g.logger.Debug("InterceptAddrDial not allowed", "p", p, "ma", ma.String())
	}

	return allow
}

func (g *ConnectionGater) InterceptAccept(cma lp2pnetwork.ConnMultiaddrs) bool {
	allow := g.BasicConnectionGater.InterceptAccept(cma)
	if !allow {
		g.logger.Debug("InterceptAccept not allowed")
	}

	return allow
}

func (g *ConnectionGater) InterceptSecured(dir lp2pnetwork.Direction, p peer.ID, cma lp2pnetwork.ConnMultiaddrs) bool {
	allow := g.BasicConnectionGater.InterceptSecured(dir, p, cma)
	if !allow {
		g.logger.Debug("InterceptSecured not allowed", "p", p)
	}

	return allow
}

func (g *ConnectionGater) InterceptUpgraded(con lp2pnetwork.Conn) (bool, lp2pcontrol.DisconnectReason) {
	return g.BasicConnectionGater.InterceptUpgraded(con)
}
