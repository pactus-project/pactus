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
	LastSeen        time.Time
	ConsensusKeys   map[bls.PublicKey]bool
	Moniker         string
	Agent           string
	PeerID          peer.ID
	Status          StatusCode
	Flags           int
	ReceivedBundles int
	InvalidBundles  int
	ReceivedBytes   int
	SendSuccess     int
	SendFailed      int
	Height          uint32
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
