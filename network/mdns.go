package network

import (
	"context"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/pactus-project/pactus/util/logger"
)

type mdnsService struct {
	ctx     context.Context
	host    lp2phost.Host
	service lp2pmdns.Service
	logger  *logger.SubLogger
}

// newMdnsService creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func newMdnsService(ctx context.Context, host lp2phost.Host, log *logger.SubLogger) *mdnsService {
	mdns := &mdnsService{
		ctx:    ctx,
		host:   host,
		logger: log,
	}
	// setup mDNS discovery to find local peers
	mdns.service = lp2pmdns.NewMdnsService(host, "pactus-mdns", mdns)

	return mdns
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (s *mdnsService) HandlePeerFound(addrInfo lp2ppeer.AddrInfo) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*10)
	defer cancel()

	if addrInfo.ID != s.host.ID() {
		s.logger.Debug("connecting to new peer", "addr", addrInfo.Addrs, "id", addrInfo.ID)
		if err := s.host.Connect(ctx, addrInfo); err != nil {
			s.logger.Error("error on connecting to peer", "id", addrInfo.ID, "error", err)
		}
	}
}

func (s *mdnsService) Start() error {
	err := s.service.Start()
	if err != nil {
		return LibP2PError{Err: err}
	}

	return nil
}

func (s *mdnsService) Stop() {
	err := s.service.Close()
	if err != nil {
		s.logger.Error("unable to close the network", "error", err)
	}
}
