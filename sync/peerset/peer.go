package peerset

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

// TODO: write tests for me

const (
	PeerFlagNodeNetwork = 0x1
)

type Peer struct {
	Status          StatusCode
	Moniker         string
	Agent           string
	PeerID          peer.ID
	PublicKey       bls.PublicKey
	Flags           int
	LastSeen        time.Time
	Height          int
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

func (p *Peer) Fingerprint() string {
	return fmt.Sprintf("{%v %v}",
		util.FingerprintPeerID(p.PeerID),
		p.Height)
}

func (p *Peer) IsKnownOrTrusted() bool {
	return p.Status == StatusCodeKnown || p.Status == StatusCodeTrusted
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
