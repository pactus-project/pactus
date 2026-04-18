package session

import (
	"time"

	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/types"
)

type Status int

const (
	Open        = Status(0)
	Completed   = Status(2)
	Uncompleted = Status(1)
)

type Session struct {
	SessionID    int
	Status       Status
	PeerID       peer.ID
	From         types.Height
	Count        uint32
	LastActivity time.Time
}

func NewSession(id int, peerID peer.ID, from types.Height, count uint32) *Session {
	return &Session{
		SessionID:    id,
		Status:       Open,
		PeerID:       peerID,
		From:         from,
		Count:        count,
		LastActivity: time.Now(),
	}
}
