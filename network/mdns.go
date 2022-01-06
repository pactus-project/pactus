package network

import (
	"context"
	"time"

	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2ppeer "github.com/libp2p/go-libp2p-core/peer"
	lp2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryServiceTag is used in our mDNS advertisements to discover other peers.
const DiscoveryServiceTag = "mdns-zarb"

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *network) HandlePeerFound(pi lp2ppeer.AddrInfo) {
	n.logger.Debug("discovered new peer", "id", pi.ID.Pretty())
	ctx, cancel := context.WithTimeout(n.ctx, time.Second*10)
	defer cancel()
	if err := n.host.Connect(ctx, pi); err != nil {
		n.logger.Error("error connecting to peer", "id", pi.ID.Pretty(), "err", err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func (n *network) setupMNSDiscovery(h lp2phost.Host) (lp2pmdns.Service, error) {
	// setup mDNS discovery to find local peers
	service := lp2pmdns.NewMdnsService(h, DiscoveryServiceTag, n)

	err := service.Start()
	if err != nil {
		return nil, err
	}

	return service, nil
}
