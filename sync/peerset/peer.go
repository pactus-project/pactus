package peerset

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

type StatusCode int

const (
	StatusCodeUnknown = StatusCode(-1)
	StatusCodeOK      = StatusCode(0)
	StatusCodeBanned  = StatusCode(1)
)

func (code StatusCode) String() string {
	switch code {
	case StatusCodeUnknown:
		return "Unknown"
	case StatusCodeOK:
		return "Ok"
	case StatusCodeBanned:
		return "Banned"
	}
	return "Invalid"
}

func (code StatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.String())
}

type Peer struct {
	lk   sync.RWMutex
	data peerData
}

type peerData struct {
	Status               StatusCode
	Moniker              string
	NodeVersion          string
	PeerID               peer.ID
	Address              *crypto.Address
	PublicKey            *bls.BLSPublicKey
	InitialBlockDownload bool
	Height               int
	ReceivedMessages     int
	InvalidMessages      int
	ReceivedBytes        int
}

func NewPeer(peerID peer.ID) *Peer {
	return &Peer{
		data: peerData{
			Status: StatusCodeUnknown,
			PeerID: peerID,
		},
	}
}
func (p *Peer) Status() StatusCode {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Status
}

func (p *Peer) Moniker() string {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Moniker
}

func (p *Peer) NodeVersion() string {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.NodeVersion
}

func (p *Peer) PeerID() peer.ID {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.PeerID
}

func (p *Peer) PublicKey() crypto.PublicKey {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.PublicKey
}

func (p *Peer) Height() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Height
}

func (p *Peer) InitialBlockDownload() bool {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.InitialBlockDownload
}

func (p *Peer) ReceivedMessages() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.ReceivedMessages
}

func (p *Peer) InvalidMessages() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.InvalidMessages
}

func (p *Peer) ReceivedBytes() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.ReceivedBytes
}

func (p *Peer) UpdateStatus(status StatusCode) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Status = status
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

func (p *Peer) UpdateNodeVersion(version string) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.NodeVersion = version
}

func (p *Peer) HasPublicKey() bool {
	return p.data.PublicKey != nil
}

func (p *Peer) UpdatePublicKey(pub crypto.PublicKey) {
	p.lk.Lock()
	defer p.lk.Unlock()

	// TODO: write test for me
	if err := pub.SanityCheck(); err == nil {
		addr := pub.Address()
		p.data.PublicKey = pub.(*bls.BLSPublicKey)
		p.data.Address = &addr
	}
}

func (p *Peer) UpdateHeight(height int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Height = height
}

func (p *Peer) IncreaseReceivedMessage() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedMessages++
}

func (p *Peer) IncreaseReceivedBytes(len int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes += len
}

func (p *Peer) UpdateInvalidMessage(msg int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidMessages = msg
}

func (p *Peer) UpdateReceivedMessage(msg int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedMessages = msg
}

func (p *Peer) UpdateReceivedBytes(len int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes = len
}

func (p *Peer) IncreaseInvalidMessage() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidMessages++
}

func (p *Peer) MarshalJSON() ([]byte, error) {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return json.Marshal(p.data)
}

func (p *Peer) Fingerprint() string {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return fmt.Sprintf("{%v %v}",
		util.FingerprintPeerID(p.data.PeerID),
		p.data.Height)
}
