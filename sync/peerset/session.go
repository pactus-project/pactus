package peerset

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
)

type Session struct {
	data sessionData
	lk   sync.RWMutex
}

type sessionData struct {
	LastActivityAt   time.Time
	PeerID           peer.ID
	SessionID        int
	LastResponseCode message.ResponseCode
}

func newSession(id int, peerID peer.ID) *Session {
	return &Session{
		data: sessionData{
			SessionID:      id,
			PeerID:         peerID,
			LastActivityAt: util.Now(),
		},
	}
}

func (s *Session) SetLastResponseCode(code message.ResponseCode) {
	s.lk.Lock()
	defer s.lk.Unlock()

	s.data.LastResponseCode = code
	s.data.LastActivityAt = util.Now()
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

func (s *Session) LastActivityAt() time.Time {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.data.LastActivityAt
}
