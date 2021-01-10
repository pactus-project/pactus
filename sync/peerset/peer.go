package peerset

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/version"
)

type Peer struct {
	lk   deadlock.RWMutex
	data peerData
}

type peerData struct {
	Moniker              string
	NodeVersion          version.Version
	PeerID               peer.ID
	Address              crypto.Address
	PublicKey            crypto.PublicKey
	InitialBlockDownload bool
	Height               int
	ReceivedMsg          int
	InvalidMsg           int
	ReceivedBytes        int
}

func NewPeer(peerID peer.ID) *Peer {
	return &Peer{
		data: peerData{
			PeerID: peerID,
		},
	}
}

func (p *Peer) Moniker() string {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.Moniker
}

func (p *Peer) NodeVersion() version.Version {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.NodeVersion
}

func (p *Peer) PeerID() peer.ID {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.PeerID
}

func (p *Peer) Address() crypto.Address {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.Address
}

func (p *Peer) PublicKey() crypto.PublicKey {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.PublicKey
}

func (p *Peer) Height() int {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.Height
}

func (p *Peer) InitialBlockDownload() bool {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.InitialBlockDownload
}

func (p *Peer) ReceivedMsg() int {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.ReceivedMsg
}

func (p *Peer) InvalidMsg() int {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.InvalidMsg
}

func (p *Peer) ReceivedBytes() int {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.ReceivedBytes
}

func (p *Peer) UpdateMoniker(moniker string) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Moniker = moniker
}

func (p *Peer) UpdateInitialBlockDownload(initialBlockDownload bool) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InitialBlockDownload = initialBlockDownload
}

func (p *Peer) UpdateNodeVersion(version version.Version) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.NodeVersion = version
}

func (p *Peer) UpdatePublicKey(pub crypto.PublicKey) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.PublicKey = pub
	p.data.Address = pub.Address()
}

func (p *Peer) UpdateHeight(height int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Height = height
}

func (p *Peer) IncreaseReceivedMessage() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedMsg++
}

func (p *Peer) IncreaseReceivedBytes(len int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes += len
}

func (p *Peer) IncreaseInvalidMessage() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidMsg++
}
