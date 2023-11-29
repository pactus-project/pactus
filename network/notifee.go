package network

import (
	"context"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2pevent "github.com/libp2p/go-libp2p/core/event"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
	"golang.org/x/exp/slices"
)

type NotifeeService struct {
	ctx              context.Context
	host             lp2phost.Host
	lp2pEventSub     lp2pevent.Subscription
	eventChannel     chan<- Event
	logger           *logger.SubLogger
	streamProtocolID lp2pcore.ProtocolID
	peerMgr          *peerMgr
	bootstrapper     bool
}

func newNotifeeService(ctx context.Context, host lp2phost.Host, eventChannel chan<- Event, peerMgr *peerMgr,
	protocolID lp2pcore.ProtocolID, bootstrapper bool, log *logger.SubLogger,
) *NotifeeService {
	eventSub, err := host.EventBus().Subscribe(lp2pevent.WildcardSubscription)
	if err != nil {
		logger.Error("failed to register for libp2p events")
	}
	notifee := &NotifeeService{
		ctx:              ctx,
		host:             host,
		lp2pEventSub:     eventSub,
		eventChannel:     eventChannel,
		streamProtocolID: protocolID,
		bootstrapper:     bootstrapper,
		peerMgr:          peerMgr,
		logger:           log,
	}
	host.Network().Notify(notifee)

	return notifee
}

func (s *NotifeeService) Start() {
	go func() {
		defer s.lp2pEventSub.Close()

		for {
			select {
			case evt := <-s.lp2pEventSub.Out():
				switch e := evt.(type) {
				case lp2pevent.EvtLocalReachabilityChanged:
					s.logger.Info("reachability changed", "reachability", e.Reachability)

				case lp2pevent.EvtPeerConnectednessChanged:
					s.logger.Debug("connectedness changed", "pid", e.Peer, "connectedness", e.Connectedness)

				case lp2pevent.EvtPeerIdentificationCompleted:
					s.logger.Debug("identification completed", "pid", e.Peer)
					s.sendConnectEven(e.Peer)

				case lp2pevent.EvtPeerProtocolsUpdated:
					s.logger.Debug("protocols updated", "pid", e.Peer, "protocols", e.Added)
					s.sendConnectEven(e.Peer)

				default:
					s.logger.Info("unhandled libp2p event", "evt", evt)
				}

			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *NotifeeService) Stop() {
}

func (s *NotifeeService) Connected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("connected to peer", "pid", pid, "direction", conn.Stat().Direction)
	s.peerMgr.AddPeer(pid, conn.RemoteMultiaddr(), conn.Stat().Direction)
}

func (s *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("disconnected from peer", "pid", pid)
	s.eventChannel <- &DisconnectEvent{PeerID: pid}

	s.peerMgr.RemovePeer(pid)
}

func (s *NotifeeService) Listen(_ lp2pnetwork.Network, ma multiaddr.Multiaddr) {
	// Handle listen event if needed.
	s.logger.Debug("notifee Listen event emitted", "addr", ma.String())
}

// ListenClose is called when your node stops listening on an address.
func (s *NotifeeService) ListenClose(_ lp2pnetwork.Network, ma multiaddr.Multiaddr) {
	// Handle listen close event if needed.
	s.logger.Debug("notifee ListenClose event emitted", "addr", ma.String())
}

func (s *NotifeeService) sendConnectEven(pid lp2pcore.PeerID) {
	protocols, err := s.host.Peerstore().GetProtocols(pid)
	if err != nil {
		s.logger.Error("unable to get supported protocols", "pid", pid)
	}
	supportStream := slices.Contains(protocols, s.streamProtocolID)
	if supportStream {
		addr := s.peerMgr.GetMultiAddr(pid)
		if supportStream && addr != nil {
			s.eventChannel <- &ConnectEvent{
				PeerID:        pid,
				RemoteAddress: addr.String(),
			}
		}
	}
}
