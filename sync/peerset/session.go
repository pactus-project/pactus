package peerset

import (
	"fmt"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

type Session struct {
	SessionID        int
	PeerID           peer.ID
	LastResponseCode payload.ResponseCode
	LastActivityAt   time.Time
}

func newSession(id int, peerID peer.ID) *Session {
	return &Session{
		SessionID:      id,
		PeerID:         peerID,
		LastActivityAt: util.Now(),
	}
}

func (s *Session) Fingerprint() string {
	return fmt.Sprintf("{id %d %v}",
		s.SessionID,
		util.FingerprintPeerID(s.PeerID))
}
