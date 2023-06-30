package sync

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
)

var _ Synchronizer = &MockSync{}

type MockSync struct {
	TestID      peer.ID
	TestPeerSet *peerset.PeerSet
}

func MockingSync(ts *testsuite.TestSuite) *MockSync {
	ps := peerset.NewPeerSet(1 * time.Second)
	pub1, _ := ts.RandomBLSKeyPair()
	pub2, _ := ts.RandomBLSKeyPair()
	pid1 := ts.RandomPeerID()
	pid2 := ts.RandomPeerID()
	ps.UpdatePeerInfo(
		pid1,
		peerset.StatusCodeKnown,
		"test-peer-1",
		version.Agent(),
		pub1,
		true)
	ps.UpdateHeight(pid1, ts.RandUint32(100000))

	ps.UpdatePeerInfo(
		pid2,
		peerset.StatusCodeBanned,
		"test-peer-2",
		version.Agent(),
		pub2,
		false)
	ps.UpdateHeight(pid1, ts.RandUint32(100000))

	return &MockSync{
		TestID:      ts.RandomPeerID(),
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
