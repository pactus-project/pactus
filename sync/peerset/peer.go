package peerset

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
)

// TODO: write tests for me

const (
	PeerFlagNodeNetwork = 0x01
)

type Peer struct {
	Status          StatusCode
	Moniker         string
	Agent           string
	PeerID          peer.ID
	ConsensusKeys   map[bls.PublicKey]bool
	Flags           int
	LastSeen        time.Time
	Height          uint32
	ReceivedBundles int
	InvalidBundles  int
	ReceivedBytes   int
	SendSuccess     int
	SendFailed      int
}

func NewPeer(peerID peer.ID) *Peer {
	return &Peer{
		ConsensusKeys: make(map[bls.PublicKey]bool),
		Status:        StatusCodeUnknown,
		PeerID:        peerID,
	}
}

func (p *Peer) IsKnownOrTrusty() bool {
	return p.Status == StatusCodeKnown || p.Status == StatusCodeTrusty
}

func (p *Peer) IsBanned() bool {
	return p.Status == StatusCodeBanned
}

func (p *Peer) IsNodeNetwork() bool {
	return util.IsFlagSet(p.Flags, PeerFlagNodeNetwork)
}

func (p *Peer) LastSeenAt() time.Time {
	return p.LastSeen
}
