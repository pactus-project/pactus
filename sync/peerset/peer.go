package peerset

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
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
	PublicKey       bls.PublicKey
	Flags           int
	LastSeen        time.Time
	Height          uint32
	ReceivedBundles int
	InvalidBundles  int
	ReceivedBytes   int
}

func NewPeer(peerID peer.ID) *Peer {
	return &Peer{
		Status: StatusCodeUnknown,
		PeerID: peerID,
	}
}

func (p *Peer) IsKnownOrTrusty() bool {
	return p.Status == StatusCodeKnown || p.Status == StatusCodeTrusty
}

func (p *Peer) IsBanned() bool {
	return p.Status == StatusCodeBanned
}

func (p *Peer) Address() crypto.Address {
	return p.PublicKey.Address()
}

func (p *Peer) SetNodeNetworkFlag(nodeNetwork bool) {
	if nodeNetwork {
		p.Flags = util.SetFlag(p.Flags, PeerFlagNodeNetwork)
	} else {
		p.Flags = util.UnsetFlag(p.Flags, PeerFlagNodeNetwork)
	}
}

func (p *Peer) IsNodeNetwork() bool {
	return util.IsFlagSet(p.Flags, PeerFlagNodeNetwork)
}
