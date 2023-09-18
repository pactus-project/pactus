package sync

import (
	"time"

	"github.com/pactus-project/pactus/crypto/bls"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/services"
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
	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	ps.UpdateInfo(
		pid1,
		"test-peer-1",
		version.Agent(),
		[]*bls.PublicKey{pub1},
		services.New(services.Network))
	ps.UpdateHeight(pid1, ts.RandHeight(), ts.RandHash())

	ps.UpdateInfo(
		pid2,
		"test-peer-2",
		version.Agent(),
		[]*bls.PublicKey{pub2},
		services.New(services.None))
	ps.UpdateHeight(pid1, ts.RandHeight(), ts.RandHash())

	return &MockSync{
		TestID:      ts.RandPeerID(),
		TestPeerSet: ps,
	}
}

func (m *MockSync) Start() error {
	return nil
}

func (m *MockSync) Stop() {
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
