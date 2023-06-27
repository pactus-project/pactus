package sync

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

var _ Synchronizer = &MockSync{}

type MockSync struct {
	TestID      peer.ID
	TestPeerSet *peerset.PeerSet
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
		"test-peer-1",
		version.Agent(),
		pub1,
		true)
	ps.UpdateHeight(pid1, util.RandUint32(100000))

	ps.UpdatePeerInfo(
		pid2,
		peerset.StatusCodeBanned,
		"test-peer-2",
		version.Agent(),
		pub2,
		false)
	ps.UpdateHeight(pid1, util.RandUint32(100000))

	return &MockSync{
		TestID:      network.TestRandomPeerID(),
		TestPeerSet: ps,
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
	return m.TestID
}

func (m *MockSync) Moniker() string {
	return "test-moniker"
}
func (m *MockSync) PeerSet() *peerset.PeerSet {
	return m.TestPeerSet
}
