package sync

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
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
	pid1 := network.TestRandomPeerID()
	pid2 := network.TestRandomPeerID()
	ps.UpdatePeerInfo(
		pid1,
		peerset.StatusCodeKnown,
		"test-1",
		version.Agent(),
		pub1,
		true)
	ps.UpdateHeight(pid1, util.RandInt32(100000))

	ps.UpdatePeerInfo(
		pid2,
		peerset.StatusCodeBanned,
		"test-1",
		version.Agent(),
		pub2,
		false)
	ps.UpdateHeight(pid1, util.RandInt32(100000))

	return &MockSync{
		ID:      network.TestRandomPeerID(),
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
