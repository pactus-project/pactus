package peerset

import (
	"sync"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/metric"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util"
)

type PeerSet struct {
	lk sync.RWMutex

	peers          map[peer.ID]*peer.Peer
	sessionManager *session.Manager
	startedAt      time.Time
	metric         metric.Metric
}

// NewPeerSet constructs a new PeerSet for managing peer information.
func NewPeerSet(sessionTimeout time.Duration) *PeerSet {
	return &PeerSet{
		peers:          make(map[peer.ID]*peer.Peer),
		sessionManager: session.NewManager(sessionTimeout),
		metric:         metric.NewMetric(),
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

// HasPeer checks if a peer with the given id exists in the set.
func (ps *PeerSet) HasPeer(pid peer.ID) bool {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	_, ok := ps.peers[pid]

	return ok
}

func (ps *PeerSet) findPeer(pid peer.ID) *peer.Peer {
	if p, ok := ps.peers[pid]; ok {
		return p
	}

	return nil
}

// FindOrCreatePeer tries to find a peer with the given pid.
// If not found, it creates a new peer and assigns the pid to it.
func (ps *PeerSet) findOrCreatePeer(pid peer.ID) *peer.Peer {
	per := ps.findPeer(pid)
	if per == nil {
		per = peer.NewPeer(pid)
		ps.peers[pid] = per
	}

	return per
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

func (ps *PeerSet) UpdateAddress(pid peer.ID, addr string, direction lp2pnetwork.Direction) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.Address = addr
	p.Direction = direction
}

func (ps *PeerSet) UpdateOutboundHelloSent(pid peer.ID, sent bool) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	p := ps.findOrCreatePeer(pid)
	p.OutboundHelloSent = sent
}

func (ps *PeerSet) UpdateStatus(pid peer.ID, status status.Status) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	peer := ps.findOrCreatePeer(pid)

	// Don't update the status if peer is banned (unless status is disconnected).
	// This helps banned peers to recover, like fixing their system time, etc.
	if peer.Status.IsBanned() && !status.IsDisconnected() {
		return
	}

	// Don't change status to connected if peer is known already
	if peer.Status.IsKnown() && status.IsConnected() {
		return
	}

	peer.Status = status

	if status.IsDisconnected() {
		peer.OutboundHelloSent = false

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

func (ps *PeerSet) UpdateInvalidMetric(pid peer.ID, bytes int64) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.metric.UpdateInvalidMetric(bytes)

	p := ps.findOrCreatePeer(pid)
	p.Metric.UpdateInvalidMetric(bytes)
}

func (ps *PeerSet) UpdateReceivedMetric(pid peer.ID, msgType message.Type, bytes int64) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.metric.UpdateReceivedMetric(msgType, bytes)

	p := ps.findOrCreatePeer(pid)
	p.Metric.UpdateReceivedMetric(msgType, bytes)
}

func (ps *PeerSet) UpdateSentMetric(pid *peer.ID, msgType message.Type, bytes int64) {
	ps.lk.Lock()
	defer ps.lk.Unlock()

	ps.metric.UpdateSentMetric(msgType, bytes)

	if pid != nil {
		p := ps.findOrCreatePeer(*pid)
		p.Metric.UpdateSentMetric(msgType, bytes)
	}
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

func (ps *PeerSet) Metric() metric.Metric {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.metric
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
	for _, peer := range ps.peers {
		if !peer.Status.IsConnectedOrKnown() {
			continue
		}

		score := peer.DownloadScore()
		totalScore += score
		peers = append(peers, scoredPeer{
			peer:  peer,
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
