package sync

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
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
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	p1 := ps.MustGetPeer(util.RandomPeerID())
	p2 := ps.MustGetPeer(util.RandomPeerID())
	p1.UpdateMoniker("test-1")
	p2.UpdateMoniker("test-2")
	p1.UpdatePublicKey(pub1)
	p2.UpdatePublicKey(pub2)
	p1.IncreaseInvalidMessage()
	p1.IncreaseReceivedBytes(100)
	p1.IncreaseReceivedMessage()
	p1.UpdateNodeVersion(version.NodeVersion)
	p2.UpdateNodeVersion(version.NodeVersion)
	p1.UpdateInitialBlockDownload(true)
	p1.UpdateHeight(100)
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

func (m *MockSync) PeerID() peer.ID {
	return m.ID
}

func (m *MockSync) Peers() []*peerset.Peer {
	return m.PeerSet.GetPeerList()
}

// AddPeer will add new peer to mocked PeerSet
func (m *MockSync) AddPeer(name string, height int) *peerset.Peer {
	newPeer := m.PeerSet.MustGetPeer(util.RandomPeerID())
	_, pub1, _ := crypto.GenerateTestKeyPair()
	newPeer.UpdateMoniker(name)
	newPeer.UpdatePublicKey(pub1)
	newPeer.IncreaseInvalidMessage()
	newPeer.IncreaseReceivedBytes(height * 8)
	newPeer.IncreaseReceivedMessage()
	newPeer.UpdateNodeVersion(version.NodeVersion)
	newPeer.UpdateHeight(height)
	return newPeer
}
