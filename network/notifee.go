package network

import (
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

type NotifeeService struct {
	eventChannel chan<- Event
	logger       *logger.SubLogger
	protocolID   protocol.ID
}

func newNotifeeService(host lp2phost.Host, eventChannel chan<- Event,
	logger *logger.SubLogger, protocolID protocol.ID,
) *NotifeeService {
	notifee := &NotifeeService{
		eventChannel: eventChannel,
		logger:       logger,
		protocolID:   protocolID,
	}
	host.Network().Notify(notifee)
	return notifee
}

func (n *NotifeeService) Connected(lp2pn lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	peerID := conn.RemotePeer()

	go func() {
		for i := 0; i < 6; i++ {
			// TODO: better way?
			// Wait to complete libp2p identify
			time.Sleep(200 * time.Millisecond)

			protocols, _ := lp2pn.Peerstore().SupportsProtocols(peerID, n.protocolID)
			if len(protocols) > 0 {
				n.logger.Info("connected to peer", "pid", peerID)
				n.eventChannel <- &ConnectEvent{PeerID: peerID}

				return
			}
		}

		n.logger.Info("this node doesn't support stream protocol", "pid", peerID)
	}()
}

func (n *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	peerID := conn.RemotePeer()
	n.logger.Info("disconnected from peer", "pid", peerID)
	n.eventChannel <- &DisconnectEvent{PeerID: peerID}
}

func (n *NotifeeService) Listen(_ lp2pnetwork.Network, ma ma.Multiaddr) {
	// Handle listen event if needed.
	n.logger.Debug("notifee Listen event emitted", "addr", ma.String())
}

// ListenClose is called when your node stops listening on an address.
func (n *NotifeeService) ListenClose(_ lp2pnetwork.Network, ma ma.Multiaddr) {
	// Handle listen close event if needed.
	n.logger.Debug("notifee ListenClose event emitted", "addr", ma.String())
}
