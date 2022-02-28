package sync

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

var _ Synchronizer = &MockSync{}

type MockSync struct {
	ID      peer.ID
	PeerSet *peerset.PeerSet
}

func MockingSync() *MockSync {
	ps := peerset.NewPeerSet(1 * time.Second)
	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	ps.UpdatePeer(
		util.RandomPeerID(),
		peerset.StatusCodeKnown,
		"test-1",
		version.Agent(),
		100,
		pub1,
		true)

	ps.UpdatePeer(
		util.RandomPeerID(),
		peerset.StatusCodeBanned,
		"test-1",
		version.Agent(),
		100,
		pub2,
		false)

	return &MockSync{
		ID:      util.RandomPeerID(),
		PeerSet: ps,
	}
}

func (m *MockSync) Start() error {
	return nil
}
func (m *MockSync) Stop() {

}
func (m *MockSync) Fingerprint() string {
	return ""
}

func (m *MockSync) SelfID() peer.ID {
	return m.ID
}

func (m *MockSync) Peers() []peerset.Peer {
	return m.PeerSet.GetPeerList()
}

// AddPeer will add new peer to mocked PeerSet
func (m *MockSync) AddPeer(name string, height int) {
	pub, _ := bls.GenerateTestKeyPair()
	ps := peerset.NewPeerSet(1 * time.Second)
	ps.UpdatePeer(
		util.RandomPeerID(),
		peerset.StatusCodeBanned,
		name,
		version.Agent(),
		height,
		pub,
		false)
}
