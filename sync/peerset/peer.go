package peerset

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
)

const (
	PeerFlagNodeNetwork = 0x01
)

type Peer struct {
	Status          StatusCode
	Moniker         string
	Agent           string
	PeerID          peer.ID
	ConsensusKeys   []*bls.PublicKey
	Flags           int
	LastSent        time.Time
	LastReceived    time.Time
	LastBlockHash   hash.Hash
	Height          uint32
	ReceivedBundles int
	InvalidBundles  int
	ReceivedBytes   map[message.Type]int64
	SentBytes       map[message.Type]int64
	SendSuccess     int
	SendFailed      int
}

func NewPeer(peerID peer.ID) *Peer {
	return &Peer{
		ConsensusKeys: make([]*bls.PublicKey, 0),
		Status:        StatusCodeConnected,
		PeerID:        peerID,
		ReceivedBytes: make(map[message.Type]int64),
		SentBytes:     make(map[message.Type]int64),
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
