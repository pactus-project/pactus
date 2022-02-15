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
	StatusCodeBanned  = StatusCode(-1)
	StatusCodeUnknown = StatusCode(0)
	StatusCodeKnown   = StatusCode(1)
	StatusCodeTrusted = StatusCode(2)
)

func (code StatusCode) String() string {
	switch code {
	case StatusCodeBanned:
		return "banned"
	case StatusCodeUnknown:
		return "unknown"
	case StatusCodeKnown:
		return "known"
	case StatusCodeTrusted:
		return "trusted"
	}
	return "invalid"
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
	Agent                string
	PeerID               peer.ID
	Address              *crypto.Address
	PublicKey            *bls.PublicKey
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

func (p *Peer) IsKnownOrTrusted() bool {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Status == StatusCodeKnown || p.data.Status == StatusCodeTrusted
}

func (p *Peer) Moniker() string {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Moniker
}

func (p *Peer) Agent() string {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.Agent
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

func (p *Peer) UpdateAgent(version string) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Agent = version
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
		p.data.PublicKey = pub.(*bls.PublicKey)
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
