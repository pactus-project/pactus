package network

import (
	"context"
	"sync"
	"sync/atomic"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2pevent "github.com/libp2p/go-libp2p/core/event"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peventbus "github.com/libp2p/go-libp2p/p2p/host/eventbus"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/pipeline"
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
	reachability     atomic.Int32

	// The events of the connections that are made before Start are kept in pending,
	// and they are sent by Start: the receivers of the pipe are registered by then.
	pendingLock sync.Mutex
	started     bool
	pending     []Event
}

func newNotifeeService(ctx context.Context, host lp2phost.Host, networkPipe pipeline.Pipeline[Event],
	peerMgr *peerMgr,
	protocolID lp2pcore.ProtocolID, log *logger.SubLogger,
) (*NotifeeService, error) {
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
		return nil, err
	}
	notifee := &NotifeeService{
		ctx:              ctx,
		host:             host,
		lp2pEventSub:     eventSub,
		networkPipe:      networkPipe,
		streamProtocolID: protocolID,
		peerMgr:          peerMgr,
		logger:           log,
	}

	return notifee, nil
}

func (s *NotifeeService) Start() {
	s.sendPendingEvents()

	go func() {
		for {
			select {
			case evt, ok := <-s.lp2pEventSub.Out():
				if !ok {
					// The subscription is closed.
					return
				}

				switch evt := evt.(type) {
				case lp2pevent.EvtLocalReachabilityChanged:
					s.logger.Info("reachability changed", "reachability", evt.Reachability)
					s.reachability.Store(int32(evt.Reachability))

				case lp2pevent.EvtPeerIdentificationCompleted:
					s.logger.Debug("identification completed", "pid", evt.Peer)
					s.sendProtocolsEvent(evt.Peer)

				case lp2pevent.EvtPeerIdentificationFailed:
					s.logger.Warn("identification failed", "pid", evt.Peer)

				case lp2pevent.EvtPeerProtocolsUpdated:
					s.logger.Debug("protocols updated", "pid", evt.Peer, "protocols", evt.Added)
					s.sendProtocolsEvent(evt.Peer)

				default:
					// unhandled libp2p event
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
	return lp2pnetwork.Reachability(s.reachability.Load())
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

	// libp2p calls Disconnected once per *connection*; a peer can hold more
	// than one at a time (simultaneous dial, TCP + QUIC), so only treat it
	// as gone, and tell the rest of the node so, once its last one closes.
	if s.peerMgr.SetPeerDisconnected(pid) {
		s.sendDisconnectEvent(pid)
	}
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
	protocolsStr := make([]string, 0, len(protocols))
	for _, p := range protocols {
		protocolsStr = append(protocolsStr, string(p))
	}

	slices.Sort(protocolsStr)
	event := &ProtocolsEvents{
		PeerID:    pid,
		Protocols: protocolsStr,
	}
	s.send(event)
}

func (s *NotifeeService) sendConnectEvent(pid lp2pcore.PeerID,
	remoteAddress multiaddr.Multiaddr, direction lp2pnetwork.Direction,
) {
	event := &ConnectEvent{
		PeerID:        pid,
		RemoteAddress: remoteAddress.String(),
		Direction:     direction,
	}
	s.send(event)
}

func (s *NotifeeService) sendDisconnectEvent(pid lp2pcore.PeerID) {
	event := &DisconnectEvent{
		PeerID: pid,
	}
	s.send(event)
}

// send sends an event to the pipe, or keeps it until the service is started.
func (s *NotifeeService) send(event Event) {
	s.pendingLock.Lock()
	if !s.started {
		s.pending = append(s.pending, event)
		s.pendingLock.Unlock()

		return
	}
	s.pendingLock.Unlock()

	s.networkPipe.Send(event)
}

// sendPendingEvents sends the events that were kept before the start, in their order.
// An event that comes at the same time waits for the lock, so it is sent after them.
func (s *NotifeeService) sendPendingEvents() {
	s.pendingLock.Lock()
	defer s.pendingLock.Unlock()

	for _, event := range s.pending {
		s.networkPipe.Send(event)
	}
	s.pending = nil
	s.started = true
}
