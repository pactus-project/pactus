package peerset

import (
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
)

type Session struct {
	SessionID      int
	PeerID         peer.ID
	Active         bool
	LastActivityAt time.Time
}

func newSession(id int, peerID peer.ID) *Session {
	return &Session{
		SessionID:      id,
		PeerID:         peerID,
		LastActivityAt: time.Now(),
	}
}
