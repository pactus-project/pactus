package session

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
)

type Status int

const (
	Open        = Status(0)
	Completed   = Status(2)
	Uncompleted = Status(1)
)

type Session struct {
	SessionID int
	Status    Status
	PeerID    peer.ID
	From      uint32
	Count     uint32
	StartedAt time.Time
}

func NewSession(id int, peerID peer.ID, from, count uint32) *Session {
	return &Session{
		SessionID: id,
		Status:    Open,
		PeerID:    peerID,
		From:      from,
		Count:     count,
		StartedAt: util.Now(),
	}
}
