package peerset

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
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
	PublicKey            crypto.PublicKey
	InitialBlockDownload bool
	Height               int
	ReceivedBundles      int
	InvalidBundles       int
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

func (p *Peer) Address() crypto.Address {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.PublicKey.Address()
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

func (p *Peer) ReceivedBundles() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.ReceivedBundles
}

func (p *Peer) InvalidBundles() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	return p.data.InvalidBundles
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

func (p *Peer) UpdatePublicKey(pub crypto.PublicKey) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.PublicKey = pub
}

func (p *Peer) UpdateHeight(height int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.Height = height
}

func (p *Peer) IncreaseReceivedBundlesCounter() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBundles++
}

func (p *Peer) IncreaseInvalidBundlesCounter() {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidBundles++
}

func (p *Peer) IncreaseReceivedBytesCounter(c int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes += c
}

func (p *Peer) UpdateInvalidBundlesCounter(n int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.InvalidBundles = n
}

func (p *Peer) UpdateReceivedBundlesCounter(n int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBundles = n
}

func (p *Peer) UpdateReceivedBytesCounter(len int) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.data.ReceivedBytes = len
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
