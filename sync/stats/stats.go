package stats

import (
	"encoding/hex"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

// Stats hold statistic data about peers' behaviors
type Stats struct {
	lk deadlock.RWMutex

	peers       map[peer.ID]*Peer
	genesisHash crypto.Hash
	maxHeight   int
}

func NewStats(genesisHash crypto.Hash) *Stats {
	return &Stats{
		genesisHash: genesisHash,
		peers:       make(map[peer.ID]*Peer),
	}
}

func (s *Stats) PeersCount() int {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return len(s.peers)
}

// MaxHeight returns the maximum height if the network that we know
//
// Note: This value might bot be accurate
// Imagine a bad peer reports that his height is 10000000
// Then we should wait until that height to start consensus.
//
func (s *Stats) MaxHeight() int {
	s.lk.RLock()
	defer s.lk.RUnlock()

	return s.maxHeight
}
func (s *Stats) getPeer(peerID peer.ID) *Peer {
	if peer, ok := s.peers[peerID]; ok {
		return peer
	}
	return nil
}

func (s *Stats) mustGetPeer(peerID peer.ID) *Peer {
	p := s.getPeer(peerID)
	if p == nil {
		p = NewPeer()
		s.peers[peerID] = p
	}
	return p
}

func (s *Stats) ParsMessage(data []byte, from peer.ID) *message.Message {
	s.lk.Lock()
	defer s.lk.Unlock()

	peer := s.mustGetPeer(from)
	peer.ReceivedMsg = peer.ReceivedMsg + 1

	msg := new(message.Message)
	err := msg.UnmarshalCBOR(data)
	if err != nil {
		peer.InvalidMsg = peer.InvalidMsg + 1
		logger.Debug("Error decoding message", "from", from.ShortString(), "data", hex.EncodeToString(data), "err", err)
		return nil
	}

	if err = msg.SanityCheck(); err != nil {
		peer.InvalidMsg = peer.InvalidMsg + 1
		logger.Debug("Peer sent us invalid msg", "peer", from.ShortString(), "msg", msg, "err", err)
		return nil
	}

	if s.badPeer(peer) {
		return nil
	}

	// Not from the same chain
	if !peer.BelongsToSameNetwork(s.genesisHash) {
		logger.Debug("Node doesn't belong to our network", "our_hash", s.genesisHash, "node_hash", peer.GenesisHash)
		return nil
	}

	switch msg.PayloadType() {
	case payload.PayloadTypeSalam:
		pld := msg.Payload.(*payload.SalamPayload)
		peer.Version = pld.NodeVersion
		peer.GenesisHash = pld.GenesisHash
		s.updateMaxHeight(pld.Height)

	case payload.PayloadTypeAleyk:
		pld := msg.Payload.(*payload.AleykPayload)
		peer.Version = pld.NodeVersion
		peer.GenesisHash = pld.GenesisHash
		s.updateMaxHeight(pld.Height)

	case payload.PayloadTypeHeartBeat:
		pld := msg.Payload.(*payload.HeartBeatPayload)
		peer.Height = pld.Pulse.Height()
		s.updateMaxHeight(pld.Pulse.Height() - 1)

	case payload.PayloadTypeProposal:
		pld := msg.Payload.(*payload.ProposalPayload)
		s.updateMaxHeight(pld.Proposal.Height() - 1)

	case payload.PayloadTypeVote:
		pld := msg.Payload.(*payload.VotePayload)
		s.updateMaxHeight(pld.Vote.Height() - 1)

	case payload.PayloadTypeVoteSet:
		//pld := msg.Payload.(*payload.VoteSetPayload)
	}

	return msg
}

func (s *Stats) badPeer(peer *Peer) bool {
	ratio := (peer.InvalidMsg * 100) / peer.ReceivedMsg

	return ratio > 10
}

func (s *Stats) updateMaxHeight(height int) {

	s.maxHeight = util.Max(s.maxHeight, height)
}
