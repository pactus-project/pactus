package network

import (
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

type NotifeeService struct {
	eventChannel chan<- Event
	logger       *logger.SubLogger
}

func newNotifeeService(host lp2phost.Host, eventChannel chan<- Event, logger *logger.SubLogger) *NotifeeService {
	notifee := &NotifeeService{
		eventChannel: eventChannel,
		logger:       logger,
	}
	host.Network().Notify(notifee)
	return notifee
}

func (n *NotifeeService) Connected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	n.eventChannel <- &ConnectEvent{PeerID: conn.RemotePeer()}
}

func (n *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	n.eventChannel <- &DisconnectEvent{PeerID: conn.RemotePeer()}
}

func (n *NotifeeService) Listen(_ lp2pnetwork.Network, ma ma.Multiaddr) {
	// Handle listen event if needed.
	n.logger.Debug("Notifee Listen event emitted", "Multiaddr", ma.String())
}

// ListenClose is called when your node stops listening on an address.
func (n *NotifeeService) ListenClose(_ lp2pnetwork.Network, ma ma.Multiaddr) {
	// Handle listen close event if needed.
	n.logger.Debug("Notifee ListenClose event emitted", "Multiaddr", ma.String())
}
