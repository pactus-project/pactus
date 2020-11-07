package stats

import (
	"encoding/hex"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/util"
)

// Stats hold statistic data about peers' behaviors
type Stats struct {
	peers       map[peer.ID]*Peer
	nodes       map[crypto.Address]*Node
	genesisHash crypto.Hash
	maxHeight   int
}

func NewStats(genesisHash crypto.Hash) *Stats {
	return &Stats{
		genesisHash: genesisHash,
		peers:       make(map[peer.ID]*Peer),
		nodes:       make(map[crypto.Address]*Node),
	}
}

func (s *Stats) PeersCount() int {
	return len(s.peers)
}

func (s *Stats) MaxHeight() int {
	return s.maxHeight
}

func (s *Stats) getPeer(peerID peer.ID) *Peer {
	if peer, ok := s.peers[peerID]; ok {
		return peer
	}
	p := NewPeer()
	s.peers[peerID] = p
	return p
}

func (s *Stats) getNode(addr crypto.Address) *Node {
	if node, ok := s.nodes[addr]; ok {
		return node
	}
	n := NewNode()
	s.nodes[addr] = n
	return n
}

func (s *Stats) IncreaseInvalidMessageCounter(peerID peer.ID) {
	peer := s.getPeer(peerID)
	peer.InvalidMsg = peer.InvalidMsg + 1
}

func (s *Stats) ParsMessage(data []byte, from peer.ID) *message.Message {
	peer := s.getPeer(from)
	peer.ReceivedMsg = peer.ReceivedMsg + 1

	msg := new(message.Message)
	err := msg.UnmarshalCBOR(data)
	if err != nil {
		peer.InvalidMsg = peer.InvalidMsg + 1
		logger.Error("Error decoding message", "from", from.ShortString(), "message", msg, "err", err)
		return nil
	}
	logger.Trace("Received a message", "from", from.ShortString(), "message", msg)

	if err = msg.SanityCheck(); err != nil {
		peer.InvalidMsg = peer.InvalidMsg + 1
		logger.Error("Peer sent us invalid msg", "from", from.ShortString(), "data", hex.EncodeToString(data), "err", err)
		return nil
	}

	node := s.getNode(msg.Initiator)

	if s.badPeer(peer) {
		return nil
	}

	if s.badNode(node) {
		return nil
	}

	//ourHeight, _ := syncer.state.LastBlockInfo()
	switch msg.PayloadType() {
	case message.PayloadTypeSalam:
		pld := msg.Payload.(*message.SalamPayload)
		node.Version = pld.Version
		node.GenesisHash = pld.GenesisHash
		s.updateMaxHeight(pld.Height)

	case message.PayloadTypeHeartBeat:
		pld := msg.Payload.(*message.HeartBeatPayload)
		node.HRS = pld.HRS
		s.updateMaxHeight(pld.HRS.Height() - 1)

	case message.PayloadTypeProposal:
		pld := msg.Payload.(*message.ProposalPayload)
		s.updateMaxHeight(pld.Proposal.Height() - 1)

	case message.PayloadTypeVote:
		pld := msg.Payload.(*message.VotePayload)
		s.updateMaxHeight(pld.Vote.Height() - 1)

	case message.PayloadTypeVoteSet:
		//pld := msg.Payload.(*message.VoteSetPayload)
	}

	return msg
}

func (s *Stats) badNode(node *Node) bool {

	return false
}

func (s *Stats) badPeer(peer *Peer) bool {
	ratio := (peer.InvalidMsg * 100) / peer.ReceivedMsg

	return ratio > 10
}

func (s *Stats) updateMaxHeight(height int) {
	s.maxHeight = util.Max(s.maxHeight, height)
}
