package firewall

import (
	"errors"
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p/core"
)

// PeerBannedError is returned when a message received from a banned peer-id or banned address.
type PeerBannedError struct {
	PeerID  lp2pcore.PeerID
	Address string
}

func (e PeerBannedError) Error() string {
	return fmt.Sprintf("peer is banned, peer-id: %s, remote-address: %s", e.PeerID, e.Address)
}

// ErrGossipMessage is returned when a stream message sends as gossip message.
var ErrGossipMessage = errors.New("receive stream message as gossip message")

// ErrStreamMessage is returned when a gossip message sends as stream message.
var ErrStreamMessage = errors.New("receive gossip message as stream message")
