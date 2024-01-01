package network

import (
	"context"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/util/logger"
)

type relayService struct {
	lk sync.RWMutex

	ctx          context.Context
	host         lp2phost.Host
	reachability lp2pnetwork.Reachability
	conf         *Config
	logger       *logger.SubLogger
}

func newRelayService(ctx context.Context, host lp2phost.Host, conf *Config, log *logger.SubLogger) *relayService {
	return &relayService{
		ctx:          ctx,
		host:         host,
		conf:         conf,
		logger:       log,
		reachability: lp2pnetwork.ReachabilityUnknown,
	}
}

func (rs *relayService) Start() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-rs.ctx.Done():
				return
			case <-ticker.C:
				rs.checkConnectivity()
			}
		}
	}()
}

func (rs *relayService) Stop() {
}

func (rs *relayService) SetReachability(reachability lp2pnetwork.Reachability) {
	rs.lk.Lock()
	rs.reachability = reachability
	rs.lk.Unlock()

	rs.checkConnectivity()
}

func (rs *relayService) Reachability() lp2pnetwork.Reachability {
	rs.lk.RLock()
	defer rs.lk.RUnlock()

	return rs.reachability
}

func (rs *relayService) checkConnectivity() {
	rs.lk.Lock()
	defer rs.lk.Unlock()

	net := rs.host.Network()
	if rs.conf.EnableRelay &&
		rs.reachability == lp2pnetwork.ReachabilityPrivate {
		for _, ai := range rs.conf.RelayAddrInfos() {
			if net.Connectedness(ai.ID) != lp2pnetwork.Connected {
				rs.logger.Info("try connecting relay node", "addr", ai.Addrs)
				err := ConnectSync(rs.ctx, rs.host, ai)
				if err != nil {
					rs.logger.Warn("unable to connect to relay node", "error", err, "addr", ai.Addrs)
				} else {
					rs.logger.Info("connect to relay node", "addr", ai.Addrs)
				}
			}
		}
	} else {
		// It is public node or relay is disabled.
		for _, ai := range rs.conf.RelayAddrInfos() {
			_ = net.ClosePeer(ai.ID)
		}
	}
}
