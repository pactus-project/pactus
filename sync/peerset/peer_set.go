package peerset

import (
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util"
)

type PeerSet struct {
	lk sync.RWMutex

	peers              map[peer.ID]*Peer
	sessions           map[int]*session.Session
	nextSessionID      int
	sessionTimeout     time.Duration
	totalSentBundles   int
	totalSentBytes     int64
	totalReceivedBytes int64
	sentBytes          map[message.Type]int64
	receivedBytes      map[message.Type]int64
	startedAt          time.Time
}

func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*Peer),
		sessions:       make(map[int]*session.Session),
		sessionTimeout: sessionTimeout,
		sentBytes:      make(map[message.Type]int64),
		receivedBytes:  make(map[message.Type]int64),
		startedAt:      time.Now(),
	}
}

func (ps *PeerSet) OpenSession(pid peer.ID, from, count uint32) *session.Session {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := session.NewSession(ps.nextSessionID, pid, from, count)
	ps.sessions[ssn.SessionID] = ssn
	ps.nextSessionID++

	p := ps.mustGetPeer(pid)
	p.TotalSessions++

	return ssn
}

func (ps *PeerSet) FindSession(sid int) *session.Session {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	ssn, ok := ps.sessions[sid]
	if ok {
		return ssn
	}

	return nil
}

func (ps *PeerSet) NumberOfSessions() int {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	return len(ps.sessions)
}

func (ps *PeerSet) HasOpenSession(pid peer.ID) bool {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	for _, ssn := range ps.sessions {
		if ssn.PeerID == pid && ssn.Status == session.Open {
			return true
		}
	}

	return false
}

type SessionStats struct {
	Total       int
	Open        int
	Completed   int
	Uncompleted int
}

func (ss *SessionStats) String() string {
	return fmt.Sprintf("total: %v, open: %v, completed: %v, uncompleted: %v",
		ss.Total, ss.Open, ss.Completed, ss.Uncompleted)
}

func (ps *PeerSet) SessionStats() SessionStats {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	total := len(ps.sessions)
	open := 0
	completed := 0
	unCompleted := 0
	for _, ssn := range ps.sessions {
		switch ssn.Status {
		case session.Open:
			open++

		case session.Completed:
			completed++

		case session.Uncompleted:
			unCompleted++
		}
	}

	return SessionStats{
		Total:       total,
		Open:        open,
		Completed:   completed,
		Uncompleted: unCompleted,
	}
}

func (ps *PeerSet) HasAnyOpenSession() bool {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	for _, ssn := range ps.sessions {
		if ssn.Status == session.Open {
			return true
		}
	}

	return false
}

func (ps *PeerSet) UpdateSessionLastActivity(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := ps.sessions[sid]
	if ssn != nil {
		ssn.LastActivity = time.Now()
	}
}

func (ps *PeerSet) SetExpiredSessionsAsUncompleted() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	for _, ssn := range ps.sessions {
		if ps.sessionTimeout < util.Now().Sub(ssn.LastActivity) {
			ssn.Status = session.Uncompleted
		}
	}
}

func (ps *PeerSet) SetSessionUncompleted(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := ps.sessions[sid]
	if ssn != nil {
		ssn.Status = session.Uncompleted
	}
}

func (ps *PeerSet) SetSessionCompleted(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := ps.sessions[sid]
	if ssn != nil {
		ssn.Status = session.Completed

		p := ps.mustGetPeer(ssn.PeerID)
		p.CompletedSessions++
	}
}

func (ps *PeerSet) RemoveAllSessions() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.sessions = make(map[int]*session.Session)
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

// GetPeer finds a peer by id and returns a copy of the peer object.
func (ps *PeerSet) GetPeer(pid peer.ID) *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.getPeer(pid)
}

func (ps *PeerSet) getPeer(pid peer.ID) *Peer {
	if p, ok := ps.peers[pid]; ok {
		return p
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
	services service.Services,
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

func (ps *PeerSet) UpdateAddress(pid peer.ID, addr, direction string) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Address = addr
	p.Direction = direction
}

func (ps *PeerSet) UpdateStatus(pid peer.ID, status StatusCode) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Status = status

	if status == StatusCodeDisconnected {
		for _, ssn := range ps.sessions {
			if ssn.PeerID == pid {
				ssn.Status = session.Uncompleted
			}
		}
	}
}

func (ps *PeerSet) UpdateProtocols(pid peer.ID, protocols []string) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.mustGetPeer(pid)
	p.Protocols = protocols
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

func (ps *PeerSet) IncreaseSentCounters(msgType message.Type, c int64, pid *peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.totalSentBundles++
	ps.totalSentBytes += c
	ps.sentBytes[msgType] += c

	if pid != nil {
		p := ps.mustGetPeer(*pid)
		p.SentBytes[msgType] += c
	}
}

func (ps *PeerSet) TotalSentBundles() int {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.totalSentBundles
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

func (ps *PeerSet) IteratePeers(consumer func(peer *Peer) (stop bool)) {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	for _, p := range ps.peers {
		stopped := consumer(p)
		if stopped {
			return
		}
	}
}

func (ps *PeerSet) Sessions() []*session.Session {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	sessions := make([]*session.Session, 0, len(ps.sessions))

	for _, ssn := range ps.sessions {
		sessions = append(sessions, ssn)
	}

	return sessions
}

// GetRandomPeer selects a random peer from the peer set based on their weights.
// The weight of each peer is determined by the number of failed and total bundles.
// Peers with higher weights are more likely to be selected.
func (ps *PeerSet) GetRandomPeer() *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	type weightedPeer struct {
		peer   *Peer
		weight int
	}

	//
	totalWeight := 0
	peers := make([]weightedPeer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.IsKnownOrTrusty() {
			continue
		}

		weight := (p.CompletedSessions + 1) * 100 / (p.TotalSessions + 1)
		totalWeight += weight
		peers = append(peers, weightedPeer{
			peer:   p,
			weight: weight,
		})
	}

	if len(peers) == 0 {
		return nil
	}

	rnd := int(util.RandUint32(uint32(totalWeight)))

	// Find the index where the random number falls
	for _, p := range peers {
		totalWeight -= p.weight

		if rnd >= totalWeight {
			return p.peer
		}
	}
	panic("unreachable code")
}
