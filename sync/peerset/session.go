package peerset

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
)

type Session struct {
	lk sync.RWMutex

	data sessionData
}

type sessionData struct {
	SessionID   int
	PeerID      peer.ID
	Uncompleted bool
	From        uint32
	To          uint32
	StartedAt   time.Time
}

func newSession(id int, peerID peer.ID, from, to uint32) *Session {
	return &Session{
		data: sessionData{
			SessionID: id,
			PeerID:    peerID,
			From:      from,
			To:        to,
			StartedAt: util.Now(),
		},
	}
}

func (s *Session) PeerID() peer.ID {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.PeerID
}

func (s *Session) SessionID() int {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.SessionID
}

func (s *Session) StartedAt() time.Time {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.StartedAt
}

func (s *Session) Range() (uint32, uint32) {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.From, s.data.To
}

func (s *Session) SetUncompleted() {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.data.Uncompleted = true
}

func (s *Session) IsUncompleted() bool {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.Uncompleted
}
