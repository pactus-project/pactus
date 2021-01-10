package peerset

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/util"
)

// TODO:
// - Add tests for peerset
// - Implementing garbage collection for peerset

type PeerSet struct {
	lk deadlock.RWMutex

	peers            map[peer.ID]*Peer
	maxClaimedHeight int
}

func NewPeerSet() *PeerSet {
	return &PeerSet{
		peers: make(map[peer.ID]*Peer),
	}
}

func (ps *PeerSet) GetPeer(peerID peer.ID) *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	return ps.getPeer(peerID)
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

func (ps *PeerSet) FindHighestPeer() *Peer {
	ps.lk.RLock()
	defer ps.lk.RUnlock()

	for _, p := range ps.peers {
		if p.Height() >= ps.maxClaimedHeight {
			return p
		}
	}
	return nil
}
