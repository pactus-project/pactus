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
	Moniker       string
	Version       version.Version
	PeerID        peer.ID
	Address       crypto.Address
	PublicKey     crypto.PublicKey
	Height        int
	ReceivedMsg   int
	InvalidMsg    int
	ReceivedBytes int
	InvalidBytes  int
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

func (p *Peer) Version() version.Version {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.Version
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

func (p *Peer) InvalidBytes() int {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.data.InvalidBytes
}

func (p *Peer) UpdateMoniker(moniker string) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Moniker = moniker
}

func (p *Peer) UpdateVersion(version version.Version) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Version = version
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

func (p *Peer) IncreaseReceivedMsg() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedMsg++
}

func (p *Peer) IncreaseInvalidMsg() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidMsg++
}

func (p *Peer) IncreaseReceivedBytes(bytes int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes += bytes
}

func (p *Peer) IncreaseInvalidBytes(bytes int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidBytes += bytes
}
