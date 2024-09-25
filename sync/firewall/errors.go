package firewall

import (
	"errors"
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p/core"
)

// ErrGossipMessage is returned when a stream message sends as gossip message.
var ErrGossipMessage = errors.New("receive stream message as gossip message")

// ErrStreamMessage is returned when a gossip message sends as stream message.
var ErrStreamMessage = errors.New("receive gossip message as stream message")

// ErrNetworkMismatch is returned when the bundle doesn't belong to this network.
var ErrNetworkMismatch = errors.New("bundle is not for this network")

// PeerBannedError is returned when a message received from a banned peer-id or banned address.
type PeerBannedError struct {
	PeerID  lp2pcore.PeerID
	Address string
}

func (e PeerBannedError) Error() string {
	return fmt.Sprintf("peer is banned, peer-id: %s, remote-address: %s", e.PeerID, e.Address)
}
