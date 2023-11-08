package network

import (
	"time"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
	"golang.org/x/exp/slices"
)

type NotifeeService struct {
	host             lp2phost.Host
	eventChannel     chan<- Event
	logger           *logger.SubLogger
	streamProtocolID protocol.ID
	peerMgr          *peerMgr
	bootstrapper     bool
}

func newNotifeeService(host lp2phost.Host, eventChannel chan<- Event, peerMgr *peerMgr,
	protocolID protocol.ID, bootstrapper bool, log *logger.SubLogger,
) *NotifeeService {
	notifee := &NotifeeService{
		host:             host,
		eventChannel:     eventChannel,
		streamProtocolID: protocolID,
		bootstrapper:     bootstrapper,
		peerMgr:          peerMgr,
		logger:           log,
	}
	host.Network().Notify(notifee)
	return notifee
}

func (n *NotifeeService) Connected(lp2pn lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	n.logger.Info("connected to peer", "pid", pid, "direction", conn.Stat().Direction)

	var protocols []lp2pcore.ProtocolID
	go func() {
		for i := 0; i < 10; i++ {
			protocols, _ = lp2pn.Peerstore().GetProtocols(pid)
			if len(protocols) > 0 {
				break
			}

			// TODO: better way?
			// Wait to complete libp2p identify
			time.Sleep(1 * time.Second)
		}

		if len(protocols) == 0 {
			n.logger.Info("unable to get supported protocols", "pid", pid)
		} else {
			n.logger.Debug("get supported protocols", "pid", pid, "protocols", protocols)
		}

		n.peerMgr.AddPeer(pid, conn.RemoteMultiaddr(), conn.Stat().Direction, protocols)

		supportStream := slices.Contains(protocols, n.streamProtocolID)
		n.eventChannel <- &ConnectEvent{
			PeerID:        pid,
			RemoteAddress: conn.RemoteMultiaddr().String(),
			SupportStream: supportStream,
		}
	}()
}

func (n *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	pid := conn.RemotePeer()
	n.logger.Info("disconnected from peer", "pid", pid)
	n.eventChannel <- &DisconnectEvent{PeerID: pid}

	n.peerMgr.RemovePeer(pid)
}

func (n *NotifeeService) Listen(_ lp2pnetwork.Network, ma multiaddr.Multiaddr) {
	// Handle listen event if needed.
	n.logger.Debug("notifee Listen event emitted", "addr", ma.String())
}

// ListenClose is called when your node stops listening on an address.
func (n *NotifeeService) ListenClose(_ lp2pnetwork.Network, ma multiaddr.Multiaddr) {
	// Handle listen close event if needed.
	n.logger.Debug("notifee ListenClose event emitted", "addr", ma.String())
}
