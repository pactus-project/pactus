package network

import (
	"context"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2pevent "github.com/libp2p/go-libp2p/core/event"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peventbus "github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
	"golang.org/x/exp/slices"
)

type NotifeeService struct {
	ctx              context.Context
	host             lp2phost.Host
	lp2pEventSub     lp2pevent.Subscription
	networkPipe      pipeline.Pipeline[Event]
	logger           *logger.SubLogger
	streamProtocolID lp2pcore.ProtocolID
	peerMgr          *peerMgr
	reachability     lp2pnetwork.Reachability
}

func newNotifeeService(ctx context.Context, host lp2phost.Host, networkPipe pipeline.Pipeline[Event],
	peerMgr *peerMgr,
	protocolID lp2pcore.ProtocolID, log *logger.SubLogger,
) *NotifeeService {
	events := []any{
		new(lp2pevent.EvtLocalReachabilityChanged),
		new(lp2pevent.EvtPeerIdentificationCompleted),
		new(lp2pevent.EvtPeerIdentificationFailed),
		new(lp2pevent.EvtPeerProtocolsUpdated),
	}
	subOptions := []lp2pevent.SubscriptionOpt{
		lp2peventbus.BufSize(1024),
	}
	eventSub, err := host.EventBus().Subscribe(events, subOptions...)
	if err != nil {
		logger.Error("failed to register for libp2p events")
	}
	notifee := &NotifeeService{
		ctx:              ctx,
		host:             host,
		lp2pEventSub:     eventSub,
		networkPipe:      networkPipe,
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
				switch evt := evt.(type) {
				case lp2pevent.EvtLocalReachabilityChanged:
					s.logger.Info("reachability changed", "reachability", evt.Reachability)
					s.reachability = evt.Reachability

				case lp2pevent.EvtPeerIdentificationCompleted:
					s.logger.Debug("identification completed", "pid", evt.Peer)
					s.sendProtocolsEvent(evt.Peer)

				case lp2pevent.EvtPeerIdentificationFailed:
					s.logger.Warn("identification failed", "pid", evt.Peer)

				case lp2pevent.EvtPeerProtocolsUpdated:
					s.logger.Debug("protocols updated", "pid", evt.Peer, "protocols", evt.Added)
					s.sendProtocolsEvent(evt.Peer)

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
	_ = s.lp2pEventSub.Close()
}

func (s *NotifeeService) Reachability() lp2pnetwork.Reachability {
	return s.reachability
}

func (s *NotifeeService) Connected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("connected to peer", "pid", pid, "direction", conn.Stat().Direction, "addr", conn.RemoteMultiaddr())

	s.peerMgr.SetPeerConnected(pid, conn.RemoteMultiaddr(), conn.Stat().Direction)
	s.sendConnectEvent(pid, conn.RemoteMultiaddr(), conn.Stat().Direction)
}

func (s *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	s.logger.Info("disconnected from peer", "pid", pid)

	s.peerMgr.SetPeerDisconnected(pid)
	s.sendDisconnectEvent(pid)
}

func (s *NotifeeService) Listen(_ lp2pnetwork.Network, ma multiaddr.Multiaddr) {
	// Handle listen event if needed.
	s.logger.Debug("notifee Listen event emitted", "addr", ma.String())
}

// ListenClose is called when the peer stops listening on an address.
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
	event := &ProtocolsEvents{
		PeerID:    pid,
		Protocols: protocolsStr,
	}
	s.networkPipe.Send(event)
}

func (s *NotifeeService) sendConnectEvent(pid lp2pcore.PeerID,
	remoteAddress multiaddr.Multiaddr, direction lp2pnetwork.Direction,
) {
	event := &ConnectEvent{
		PeerID:        pid,
		RemoteAddress: remoteAddress.String(),
		Direction:     direction.String(),
	}
	s.networkPipe.Send(event)
}

func (s *NotifeeService) sendDisconnectEvent(pid lp2pcore.PeerID) {
	event := &DisconnectEvent{
		PeerID: pid,
	}
	s.networkPipe.Send(event)
}
