package peerset

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/util"
)

// TODO:
// - Add tests for peerset
// - Implementing garbage collection for peerset

type PeerSet struct {
	lk deadlock.RWMutex

	peers            map[peer.ID]*Peer
	sessions         map[int]*Session
	nextSessionID    int
	maxClaimedHeight int
	sessionTimeout   time.Duration
}

func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*Peer),
		sessions:       make(map[int]*Session),
		sessionTimeout: sessionTimeout,
	}
}

func (ps *PeerSet) GetPeer(peerID peer.ID) *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.getPeer(peerID)
}

func (ps *PeerSet) OpenSession(peerID peer.ID) *Session {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	s := newSession(ps.nextSessionID, peerID)
	ps.sessions[s.SessionID()] = s
	ps.nextSessionID++

	return s
}

func (ps *PeerSet) FindSession(id int) *Session {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	s, ok := ps.sessions[id]
	if ok {
		return s
	}

	return nil
}

func (ps *PeerSet) NumberOfOpenSessions() int {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.removeExpiredSessions()

	return len(ps.sessions)
}

func (ps *PeerSet) HasAnyOpenSession() bool {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.removeExpiredSessions()

	return len(ps.sessions) != 0
}

func (ps *PeerSet) removeExpiredSessions() {
	// First remove old sessions
	for id, s := range ps.sessions {
		if ps.sessionTimeout < util.Now().Sub(s.LastActivityAt()) {
			delete(ps.sessions, id)
		}
	}
}

func (ps *PeerSet) CloseSession(id int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	delete(ps.sessions, id)
}

func (ps *PeerSet) getPeer(peerID peer.ID) *Peer {
	if peer, ok := ps.peers[peerID]; ok {
		return peer
	}
	return nil
}

func (ps *PeerSet) Len() int {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return len(ps.peers)
}

// MaxClaimedHeight returns the maximum calimed height
//
// Note: This value might not be accurate
// A bad peer can claim invalid height
//
func (ps *PeerSet) MaxClaimedHeight() int {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.maxClaimedHeight
}

func (ps *PeerSet) UpdateMaxClaimedHeight(h int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.maxClaimedHeight = util.Max(ps.maxClaimedHeight, h)
}

func (ps *PeerSet) MustGetPeer(peerID peer.ID) *Peer {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.getPeer(peerID)
	if p == nil {
		p = NewPeer(peerID)
		ps.peers[peerID] = p
	}
	return p
}

// TODO: write test for me
func (ps *PeerSet) Clear() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.peers = make(map[peer.ID]*Peer)
	ps.sessions = make(map[int]*Session)
	ps.maxClaimedHeight = 0
}

func (ps *PeerSet) RemovePeer(peerID peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	delete(ps.peers, peerID)
}

func (ps *PeerSet) GetPeerList() []*Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	l := make([]*Peer, len(ps.peers))
	i := 0
	for _, p := range ps.peers {
		l[i] = p
		i++
	}
	return l
}

func (ps *PeerSet) GetRandomPeer() *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	i := util.RandInt(len(ps.peers))
	for _, p := range ps.peers {
		i--
		if i <= 0 {
			return p
		}
	}

	return nil
}
