package peerset

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/util"
)

// TODO:
// - Add tests for peerset
// - Is it thread safe (GetPeer and IteratePeers) ??

type PeerSet struct {
	lk sync.RWMutex

	peers              map[peer.ID]*Peer
	sessions           map[int]*Session
	nextSessionID      int
	sessionTimeout     time.Duration
	totalSentBytes     int64
	totalReceivedBytes int64
	sentBytes          map[message.Type]int64
	receivedBytes      map[message.Type]int64
	startedAt          time.Time
}

func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*Peer),
		sessions:       make(map[int]*Session),
		sessionTimeout: sessionTimeout,
		sentBytes:      make(map[message.Type]int64),
		receivedBytes:  make(map[message.Type]int64),
		startedAt:      time.Now(),
	}
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

func (ps *PeerSet) HasOpenSession(pid peer.ID) bool {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	for _, s := range ps.sessions {
		if s.PeerID() == pid {
			return true
		}
	}

	return false
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

// GetPeer finds a peer by id and returns a copy of the peer object.
func (ps *PeerSet) GetPeer(pid peer.ID) Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	p := ps.getPeer(pid)
	if p != nil {
		return *p
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

func (ps *PeerSet) UpdateInfo(
	pid peer.ID,
	moniker string,
	agent string,
	consKeys []*bls.PublicKey,
	services services.Services,
) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Moniker = moniker
	p.Agent = agent
	p.ConsensusKeys = consKeys
	p.Services = services
}

func (ps *PeerSet) UpdateHeight(pid peer.ID, height uint32, lastBlockHash hash.Hash) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Height = height
	p.LastBlockHash = lastBlockHash
}

func (ps *PeerSet) UpdateStatus(pid peer.ID, status StatusCode) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Status = status
}

func (ps *PeerSet) UpdateLastSent(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.LastSent = time.Now()
}

func (ps *PeerSet) UpdateLastReceived(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.LastReceived = time.Now()
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

func (ps *PeerSet) IncreaseReceivedBytesCounter(pid peer.ID, msgType message.Type, c int64) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.ReceivedBytes[msgType] += c

	ps.totalReceivedBytes += c
	ps.receivedBytes[msgType] += c
}

func (ps *PeerSet) IncreaseSentBytesCounter(msgType message.Type, c int64, pid *peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.totalSentBytes += c
	ps.sentBytes[msgType] += c

	if pid != nil {
		p := ps.mustGetPeer(*pid)
		p.SentBytes[msgType] += c
	}
}

func (ps *PeerSet) TotalSentBytes() int64 {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.totalSentBytes
}

func (ps *PeerSet) TotalReceivedBytes() int64 {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.totalReceivedBytes
}

func (ps *PeerSet) SentBytesMessageType(msgType message.Type) int64 {
	if sentBytes, ok := ps.sentBytes[msgType]; ok {
		return sentBytes
	}
	return 0
}

func (ps *PeerSet) ReceivedBytesMessageType(msgType message.Type) int64 {
	if receivedBytes, ok := ps.receivedBytes[msgType]; ok {
		return receivedBytes
	}
	return 0
}

func (ps *PeerSet) SentBytes() map[message.Type]int64 {
	return ps.sentBytes
}

func (ps *PeerSet) ReceivedBytes() map[message.Type]int64 {
	return ps.receivedBytes
}

func (ps *PeerSet) StartedAt() time.Time {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.startedAt
}

func (ps *PeerSet) IteratePeers(consumer func(peer *Peer)) {
	for _, p := range ps.peers {
		consumer(p)
	}
}
