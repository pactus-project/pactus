package sync

import (
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
)

var _ Synchronizer = &MockSync{}

// MockSync is a mock implementation of the Synchronizer interface for testing.
type MockSync struct {
	TestID       peer.ID
	TestPeerSet  *peerset.PeerSet
	TestServices service.Services
}

func MockingSync(ts *testsuite.TestSuite) *MockSync {
	peerSet := peerset.NewPeerSet(1 * time.Second)
	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	peerSet.UpdateInfo(
		pid1,
		"test-peer-1",
		version.NodeAgent.String(),
		[]*bls.PublicKey{pub1},
		service.New(service.FullNode))
	peerSet.UpdateHeight(pid1, ts.RandHeight(), ts.RandHash())

	peerSet.UpdateInfo(
		pid2,
		"test-peer-2",
		version.NodeAgent.String(),
		[]*bls.PublicKey{pub2},
		service.New(service.None))
	peerSet.UpdateHeight(pid1, ts.RandHeight(), ts.RandHash())

	services := service.New()

	return &MockSync{
		TestID:       ts.RandPeerID(),
		TestPeerSet:  peerSet,
		TestServices: services,
	}
}

func (*MockSync) Start() error {
	return nil
}

func (*MockSync) Stop() {
}

func (m *MockSync) SelfID() peer.ID {
	return m.TestID
}

func (*MockSync) Moniker() string {
	return "test-moniker"
}

func (m *MockSync) PeerSet() *peerset.PeerSet {
	return m.TestPeerSet
}

func (m *MockSync) Services() service.Services {
	return m.TestServices
}

func (*MockSync) ClockOffset() (time.Duration, error) {
	return 1 * time.Second, nil
}

func (*MockSync) IsClockOutOfSync() bool {
	return false
}
