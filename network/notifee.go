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
	reachability     lp2pnetwork.Reachability
}

func newNotifeeService(ctx context.Context, host lp2phost.Host, eventChannel chan<- Event,
	peerMgr *peerMgr,
	protocolID lp2pcore.ProtocolID, log *logger.SubLogger,
) *NotifeeService {
	events := []interface{}{
		new(lp2pevent.EvtLocalReachabilityChanged),
		new(lp2pevent.EvtPeerIdentificationCompleted),
		new(lp2pevent.EvtPeerProtocolsUpdated),
	}
	eventSub, err := host.EventBus().Subscribe(events)
	if err != nil {
		logger.Error("failed to register for libp2p events")
	}
	notifee := &NotifeeService{
		ctx:              ctx,
		host:             host,
		lp2pEventSub:     eventSub,
		eventChannel:     eventChannel,
		streamProtocolID: protocolID,
		peerMgr:          peerMgr,
		logger:           log,
		reachability:     lp2pnetwork.ReachabilityUnknown,
	}

	return notifee
}

func (s *NotifeeService) Start() {
	go func() {
		for {
			select {
			case evt := <-s.lp2pEventSub.Out():
				switch e := evt.(type) {
				case lp2pevent.EvtLocalReachabilityChanged:
					s.logger.Info("reachability changed", "reachability", e.Reachability)
					s.reachability = e.Reachability

				case lp2pevent.EvtPeerIdentificationCompleted:
					s.logger.Debug("identification completed", "pid", e.Peer)
					s.sendProtocolsEvent(e.Peer)

				case lp2pevent.EvtPeerProtocolsUpdated:
					s.logger.Debug("protocols updated", "pid", e.Peer, "protocols", e.Added)
					s.sendProtocolsEvent(e.Peer)

				default:
					s.logger.Debug("unhandled libp2p event", "event", evt)
				}

			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *NotifeeService) Stop() {
	s.lp2pEventSub.Close()
}

func (s *NotifeeService) Reachability() lp2pnetwork.Reachability {
	return s.reachability
}

func (s *NotifeeService) Connected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("connected to peer", "pid", pid, "direction", conn.Stat().Direction, "addr", conn.RemoteMultiaddr())

	s.peerMgr.AddPeer(pid, conn.RemoteMultiaddr(), conn.Stat().Direction)
	s.sendConnectEvent(pid, conn.RemoteMultiaddr(), conn.Stat().Direction)
}

func (s *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("disconnected from peer", "pid", pid)

	s.peerMgr.RemovePeer(pid)
	s.sendDisconnectEvent(pid)
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

func (s *NotifeeService) sendProtocolsEvent(pid lp2pcore.PeerID) {
	protocols, _ := s.host.Peerstore().GetProtocols(pid)
	protocolsStr := []string{}
	for _, p := range protocols {
		protocolsStr = append(protocolsStr, string(p))
	}

	slices.Sort(protocolsStr)
	supportStream := slices.Contains(protocols, s.streamProtocolID)
	s.eventChannel <- &ProtocolsEvents{
		PeerID:        pid,
		Protocols:     protocolsStr,
		SupportStream: supportStream,
	}
}

func (s *NotifeeService) sendConnectEvent(pid lp2pcore.PeerID,
	remoteAddress multiaddr.Multiaddr, direction lp2pnetwork.Direction,
) {
	s.eventChannel <- &ConnectEvent{
		PeerID:        pid,
		RemoteAddress: remoteAddress.String(),
		Direction:     direction.String(),
	}
}

func (s *NotifeeService) sendDisconnectEvent(pid lp2pcore.PeerID) {
	s.eventChannel <- &DisconnectEvent{PeerID: pid}
}
