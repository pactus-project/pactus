package network

import (
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	lp2phostbasic "github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
	"golang.org/x/exp/slices"
)

type NotifeeService struct {
	host         lp2phost.Host
	eventChannel chan<- Event
	logger       *logger.SubLogger
	protocolID   protocol.ID
	bootstrapper bool
}

func newNotifeeService(host lp2phost.Host, eventChannel chan<- Event,
	log *logger.SubLogger, protocolID protocol.ID, bootstrapper bool,
) *NotifeeService {
	notifee := &NotifeeService{
		host:         host,
		eventChannel: eventChannel,
		logger:       log,
		protocolID:   protocolID,
		bootstrapper: bootstrapper,
	}
	host.Network().Notify(notifee)
	return notifee
}

func (n *NotifeeService) Connected(lp2pn lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	peerID := conn.RemotePeer()
	n.logger.Info("connected to peer", "pid", peerID)

	go func() {
		for i := 0; i < 10; i++ {
			// TODO: better way?
			// Wait to complete libp2p identify
			time.Sleep(1 * time.Second)

			basicHost, ok := n.host.(*lp2phostbasic.BasicHost)
			if ok {
				idService := basicHost.IDService()
				idService.IdentifyConn(conn)
			}

			peerStore := lp2pn.Peerstore()
			protocols, _ := peerStore.GetProtocols(peerID)
			if len(protocols) > 0 {
				if slices.Contains(protocols, n.protocolID) {
					n.logger.Debug("peer supports the stream protocol",
						"pid", peerID, "protocols", protocols)

					n.eventChannel <- &ConnectEvent{PeerID: peerID}
				} else {
					n.logger.Debug("peer doesn't support the stream protocol",
						"pid", peerID, "protocols", protocols)
				}
				return
			}
		}

		n.logger.Info("unable to get supported protocols", "pid", peerID)
	}()
}

func (n *NotifeeService) Disconnected(_ lp2pnetwork.Network, conn lp2pnetwork.Conn) {
	peerID := conn.RemotePeer()
	n.logger.Info("disconnected from peer", "pid", peerID)
	n.eventChannel <- &DisconnectEvent{PeerID: peerID}
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
