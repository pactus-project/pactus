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
func newMdnsService(ctx context.Context, host lp2phost.Host, logger *logger.SubLogger) *mdnsService {
	mdns := &mdnsService{
		ctx:    ctx,
		host:   host,
		logger: logger,
	}
	// setup mDNS discovery to find local peers
	mdns.service = lp2pmdns.NewMdnsService(host, "pactus-mdns", mdns)

	return mdns
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (s *mdnsService) HandlePeerFound(pi lp2ppeer.AddrInfo) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*10)
	defer cancel()

	if pi.ID != s.host.ID() {
		s.logger.Debug("connecting to new peer", "addr", pi.Addrs, "id", pi.ID.Pretty())
		if err := s.host.Connect(ctx, pi); err != nil {
			s.logger.Error("error on connecting to peer", "id", pi.ID.Pretty(), "error", err)
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
