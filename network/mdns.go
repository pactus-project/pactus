package network

import (
	"context"
	"sync"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/pactus-project/gopkg/logger"
)

// maxMdnsConnects is the maximum number of peers that are dialed at the same
// time after being announced via mDNS.
const maxMdnsConnects = 8

type mdnsService struct {
	ctx     context.Context
	host    lp2phost.Host
	service lp2pmdns.Service
	logger  *logger.SubLogger

	lk       sync.Mutex
	inFlight map[lp2ppeer.ID]struct{}
}

// newMdnsService creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func newMdnsService(ctx context.Context, host lp2phost.Host, log *logger.SubLogger) *mdnsService {
	mdns := &mdnsService{
		ctx:      ctx,
		host:     host,
		logger:   log,
		inFlight: make(map[lp2ppeer.ID]struct{}),
	}
	// setup mDNS discovery to find local peers
	mdns.service = lp2pmdns.NewMdnsService(host, "pactus-mdns", mdns)

	return mdns
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
// A peer that is already being dialed is not dialed again, and the number of
// simultaneous dials is bounded. A skipped peer is dialed on a later announcement.
func (s *mdnsService) HandlePeerFound(addrInfo lp2ppeer.AddrInfo) {
	if addrInfo.ID == s.host.ID() {
		return
	}

	if !s.acquire(addrInfo.ID) {
		s.logger.Debug("skip connecting to new peer", "id", addrInfo.ID)

		return
	}
	defer s.release(addrInfo.ID)

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*10)
	defer cancel()

	s.logger.Debug("connecting to new peer", "addr", addrInfo.Addrs, "id", addrInfo.ID)
	if err := s.host.Connect(ctx, addrInfo); err != nil {
		s.logger.Debug("error on connecting to peer", "id", addrInfo.ID, "error", err)
	}
}

func (s *mdnsService) acquire(pid lp2ppeer.ID) bool {
	s.lk.Lock()
	defer s.lk.Unlock()

	if _, ok := s.inFlight[pid]; ok || len(s.inFlight) >= maxMdnsConnects {
		return false
	}
	s.inFlight[pid] = struct{}{}

	return true
}

func (s *mdnsService) release(pid lp2ppeer.ID) {
	s.lk.Lock()
	defer s.lk.Unlock()

	delete(s.inFlight, pid)
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
