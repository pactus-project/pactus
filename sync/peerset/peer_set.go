package peerset

import (
	"sync"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util"
)

type PeerSet struct {
	lk sync.RWMutex

	peers              map[peer.ID]*peer.Peer
	sessionManager     *session.Manager
	totalSentBundles   int
	totalSentBytes     int64
	totalReceivedBytes int64
	sentBytes          map[message.Type]int64
	receivedBytes      map[message.Type]int64
	startedAt          time.Time
}

// NewPeerSet constructs a new PeerSet for managing peer information.
func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*peer.Peer),
		sessionManager: session.NewManager(sessionTimeout),
		sentBytes:      make(map[message.Type]int64),
		receivedBytes:  make(map[message.Type]int64),
		startedAt:      time.Now(),
	}
}

// OpenSession opens a new session for downloading blocks and returns the session ID.
func (ps *PeerSet) OpenSession(pid peer.ID, from, count uint32) int {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := ps.sessionManager.OpenSession(pid, from, count)

	p := ps.findOrCreatePeer(pid)
	p.TotalSessions++

	return ssn.SessionID
}

// NumberOfSessions returns the total number of sessions.
func (ps *PeerSet) NumberOfSessions() int {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	return ps.sessionManager.NumberOfSessions()
}

// HasOpenSession checks if the specified peer has an open session for downloading blocks.
// Note that a peer may have more than one session.
func (ps *PeerSet) HasOpenSession(pid peer.ID) bool {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.sessionManager.HasOpenSession(pid)
}

func (ps *PeerSet) SessionStats() session.Stats {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.sessionManager.Stats()
}

func (ps *PeerSet) HasAnyOpenSession() bool {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.sessionManager.HasAnyOpenSession()
}

func (ps *PeerSet) UpdateSessionLastActivity(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.sessionManager.UpdateSessionLastActivity(sid)
}

func (ps *PeerSet) SetExpiredSessionsAsUncompleted() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.sessionManager.SetExpiredSessionsAsUncompleted()
}

func (ps *PeerSet) SetSessionUncompleted(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.sessionManager.SetSessionUncompleted(sid)
}

func (ps *PeerSet) SetSessionCompleted(sid int) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ssn := ps.sessionManager.SetSessionCompleted(sid)
	if ssn != nil {
		p := ps.findOrCreatePeer(ssn.PeerID)
		p.CompletedSessions++
	}
}

func (ps *PeerSet) RemoveAllSessions() {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.sessionManager.RemoveAllSessions()
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
func (ps *PeerSet) GetPeer(pid peer.ID) *peer.Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.findPeer(pid)
}

// GetPeerStatus finds a peer by id and returns the status of the Peer.
func (ps *PeerSet) GetPeerStatus(pid peer.ID) status.Status {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	p := ps.findPeer(pid)
	if p != nil {
		return p.Status
	}

	return status.StatusUnknown
}

func (ps *PeerSet) findPeer(pid peer.ID) *peer.Peer {
	if p, ok := ps.peers[pid]; ok {
		return p
	}

	return nil
}

func (ps *PeerSet) findOrCreatePeer(pid peer.ID) *peer.Peer {
	p := ps.findPeer(pid)
	if p == nil {
		p = peer.NewPeer(pid)
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

	p := ps.findOrCreatePeer(pid)
	p.Moniker = moniker
	p.Agent = agent
	p.ConsensusKeys = consKeys
	p.Services = services
}

func (ps *PeerSet) UpdateHeight(pid peer.ID, height uint32, lastBlockHash hash.Hash) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.Height = height
	p.LastBlockHash = lastBlockHash
}

func (ps *PeerSet) UpdateAddress(pid peer.ID, addr, direction string) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.Address = addr
	p.Direction = direction
}

func (ps *PeerSet) UpdateStatus(pid peer.ID, s status.Status) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)

	if !p.Status.IsBanned() || // Don't update the status if peer is banned
		// Don't change status to connected if peer is known already
		!(p.Status.IsKnown() && s.IsConnected()) {
		p.Status = s
	}

	if s.IsDisconnected() {
		for _, ssn := range ps.sessionManager.Sessions() {
			if ssn.PeerID == pid {
				ssn.Status = session.Uncompleted
			}
		}
	}
}

func (ps *PeerSet) UpdateProtocols(pid peer.ID, protocols []string) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.Protocols = protocols
}

func (ps *PeerSet) UpdateLastSent(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.LastSent = time.Now()
}

func (ps *PeerSet) UpdateLastReceived(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.LastReceived = time.Now()
}

func (ps *PeerSet) IncreaseReceivedBundlesCounter(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.ReceivedBundles++
}

func (ps *PeerSet) IncreaseInvalidBundlesCounter(pid peer.ID) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.InvalidBundles++
}

func (ps *PeerSet) IncreaseReceivedBytesCounter(pid peer.ID, msgType message.Type, c int64) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
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
		p := ps.findOrCreatePeer(*pid)
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

func (ps *PeerSet) IteratePeers(consumer func(p *peer.Peer) (stop bool)) {
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

	return ps.sessionManager.Sessions()
}

// GetRandomPeer selects a random peer from the peer set based on their download score.
// Peers with higher score are more likely to be selected.
func (ps *PeerSet) GetRandomPeer() *peer.Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	type scoredPeer struct {
		peer  *peer.Peer
		score int
	}

	//
	totalScore := 0
	peers := make([]scoredPeer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.Status.IsConnectedOrKnown() {
			continue
		}

		score := p.DownloadScore()
		totalScore += score
		peers = append(peers, scoredPeer{
			peer:  p,
			score: score,
		})
	}

	if len(peers) == 0 {
		return nil
	}

	rnd := int(util.RandUint32(uint32(totalScore)))

	// Find the index where the random number falls
	for _, p := range peers {
		totalScore -= p.score

		if rnd >= totalScore {
			return p.peer
		}
	}
	panic("unreachable code")
}
