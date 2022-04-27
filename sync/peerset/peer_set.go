package peerset

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

// TODO:
// - Add tests for peerset
// - Implementing garbage collection for peerset

type PeerSet struct {
	lk sync.RWMutex

	peers            map[peer.ID]*Peer
	sessions         map[int]*Session
	nextSessionID    int
	maxClaimedHeight int32
	sessionTimeout   time.Duration
}

func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*Peer),
		sessions:       make(map[int]*Session),
		sessionTimeout: sessionTimeout,
	}
}

/// GetPeer returns a cloned peer
func (ps *PeerSet) GetPeer(pid peer.ID) Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	p := ps.getPeer(pid)
	if p != nil {
		return *p
	}

	return Peer{}
}

func (ps *PeerSet) OpenSession(pid peer.ID) *Session {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	s := newSession(ps.nextSessionID, pid)
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
func (ps *PeerSet) MaxClaimedHeight() int32 {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.maxClaimedHeight
}

func (ps *PeerSet) Clear() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.peers = make(map[peer.ID]*Peer)
	ps.sessions = make(map[int]*Session)
	ps.maxClaimedHeight = 0
}

func (ps *PeerSet) RemovePeer(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	delete(ps.peers, pid)
}

func (ps *PeerSet) GetPeerList() []Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	l := make([]Peer, len(ps.peers))
	i := 0
	for _, p := range ps.peers {
		l[i] = *p
		i++
	}
	return l
}

func (ps *PeerSet) GetRandomPeer() Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	i := util.RandInt32(int32(len(ps.peers)))
	for _, p := range ps.peers {
		i--
		if i <= 0 {
			return *p
		}
	}

	return Peer{}
}

func (ps *PeerSet) getPeer(pid peer.ID) *Peer {
	if peer, ok := ps.peers[pid]; ok {
		return peer
	}
	return nil
}

func (ps *PeerSet) mustGetPeer(pid peer.ID) *Peer {
	p := ps.getPeer(pid)
	if p == nil {
		p = NewPeer(pid)
		ps.peers[pid] = p
	}
	return p
}

func (ps *PeerSet) UpdatePeerInfo(
	pid peer.ID,
	status StatusCode,
	moniker string,
	agent string,
	publicKey *bls.PublicKey,
	nodeNetwork bool) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Status = status
	p.Moniker = moniker
	p.Agent = agent
	p.PublicKey = *publicKey
	p.SetNodeNetworkFlag(nodeNetwork)
}

func (ps *PeerSet) UpdateHeight(pid peer.ID, height int32) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Height = height
	ps.maxClaimedHeight = util.Max32(ps.maxClaimedHeight, height)
}

func (ps *PeerSet) UpdateStatus(pid peer.ID, status StatusCode) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Status = status
}

func (ps *PeerSet) UpdateLastSeen(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.LastSeen = time.Now()
}

func (ps *PeerSet) IncreaseReceivedBundlesCounter(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.ReceivedBundles++
}

func (ps *PeerSet) IncreaseInvalidBundlesCounter(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.InvalidBundles++
}

func (ps *PeerSet) IncreaseReceivedBytesCounter(pid peer.ID, c int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.ReceivedBytes += c
}
