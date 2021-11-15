package peerset

import (
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

type Session struct {
	lk   sync.RWMutex
	data sessionData
}

type sessionData struct {
	SessionID        int
	PeerID           peer.ID
	LastResponseCode payload.ResponseCode
	LastActivityAt   time.Time
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

func (s *Session) SetLastResponseCode(code payload.ResponseCode) {
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

func (s *Session) Fingerprint() string {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return fmt.Sprintf("{id %d %v}",
		s.data.SessionID,
		util.FingerprintPeerID(s.data.PeerID))
}
