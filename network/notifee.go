package network

import (
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

type NotifeeService struct {
	eventChannel chan<- Event
}

func NewNotifeeService(host lp2phost.Host, eventChannel chan<- Event) *NotifeeService {
	notifee := &NotifeeService{
		eventChannel: eventChannel,
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
func (n *NotifeeService) Listen(_ lp2pnetwork.Network, _ ma.Multiaddr) {
	// Handle listen event if needed.
}

// ListenClose is called when your node stops listening on an address.
func (n *NotifeeService) ListenClose(_ lp2pnetwork.Network, _ ma.Multiaddr) {
	// Handle listen close event if needed.
}
