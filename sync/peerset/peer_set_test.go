package peerset

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/service"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPeerSet(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	peerSet := NewPeerSet(time.Second)

	pk1, _ := ts.RandBLSKeyPair()
	pk2, _ := ts.RandBLSKeyPair()
	pk3, _ := ts.RandBLSKeyPair()
	pk4, _ := ts.RandBLSKeyPair()
	pk5, _ := ts.RandBLSKeyPair()
	pid1 := peer.ID("peer1")
	pid2 := peer.ID("peer2")
	pid3 := peer.ID("peer3")
	peerSet.UpdateInfo(pid1, "Moniker1", "Agent1",
		[]*bls.PublicKey{pk1, pk2}, service.New(service.Network))
	peerSet.UpdateInfo(pid2, "Moniker2", "Agent2",
		[]*bls.PublicKey{pk3}, service.New(service.None))
	peerSet.UpdateInfo(pid3, "Moniker3", "Agent3",
		[]*bls.PublicKey{pk4, pk5}, service.New(service.Network))

	t.Run("Testing Len", func(t *testing.T) {
		assert.Equal(t, 3, peerSet.Len())
	})

	t.Run("Testing GetPeerList", func(t *testing.T) {
		peerList := peerSet.GetPeerList()

		// Verify that the peer list contains the expected peers
		expectedPeerIDs := []peer.ID{pid1, pid2, pid3}
		for _, expectedID := range expectedPeerIDs {
			found := false
			for _, peer := range peerList {
				if peer.PeerID == expectedID {
					found = true
					break
				}
			}
			assert.True(t, found, "Peer with ID %s not found in the peer list", expectedID)
		}
	})

	t.Run("Testing GetPeer", func(t *testing.T) {
		p := peerSet.GetPeer(pid2)
		assert.Equal(t, pid2, p.PeerID)
		assert.Equal(t, StatusCodeUnknown, p.Status)

		p = peerSet.GetPeer(ts.RandPeerID())
		assert.Equal(t, peer.ID(""), p.PeerID)
		assert.Equal(t, StatusCodeUnknown, p.Status)
	})

	t.Run("Testing PublicKeys", func(t *testing.T) {
		p := peerSet.GetPeer(pid3)

		assert.Contains(t, p.ConsensusKeys, pk4)
		assert.Contains(t, p.ConsensusKeys, pk5)
	})

	t.Run("Testing counters", func(t *testing.T) {
		peerSet.IncreaseInvalidBundlesCounter(pid1)
		peerSet.IncreaseReceivedBundlesCounter(pid1)
		peerSet.IncreaseReceivedBytesCounter(pid1, message.TypeBlocksResponse, 100)
		peerSet.IncreaseReceivedBytesCounter(pid1, message.TypeTransactions, 150)
		peerSet.IncreaseSentBytesCounter(message.TypeBlocksRequest, 200, nil)
		peerSet.IncreaseSentBytesCounter(message.TypeBlocksRequest, 250, &pid1)

		peer1 := peerSet.getPeer(pid1)

		receivedBytes := make(map[message.Type]int64)
		receivedBytes[message.TypeBlocksResponse] = 100
		receivedBytes[message.TypeTransactions] = 150

		sentBytes := make(map[message.Type]int64)
		sentBytes[message.TypeBlocksRequest] = 450

		assert.Equal(t, peer1.InvalidBundles, 1)
		assert.Equal(t, peer1.ReceivedBundles, 1)
		assert.Equal(t, peer1.ReceivedBytes[message.TypeBlocksResponse], int64(100))
		assert.Equal(t, peer1.ReceivedBytes[message.TypeTransactions], int64(150))
		assert.Equal(t, peer1.SentBytes[message.TypeBlocksRequest], int64(250))

		assert.Equal(t, peerSet.TotalReceivedBytes(), int64(250))
		assert.Equal(t, peerSet.ReceivedBytesMessageType(message.TypeBlocksResponse), int64(100))
		assert.Equal(t, peerSet.ReceivedBytesMessageType(message.TypeTransactions), int64(150))
		assert.Equal(t, peerSet.ReceivedBytes(), receivedBytes)
		assert.Equal(t, peerSet.TotalSentBytes(), int64(450))
		assert.Equal(t, peerSet.SentBytesMessageType(message.TypeBlocksRequest), int64(450))
		assert.Equal(t, peerSet.SentBytes(), sentBytes)
	})

	t.Run("Testing UpdateStatus", func(t *testing.T) {
		peerSet.UpdateStatus(pid1, StatusCodeBanned)

		peer1 := peerSet.getPeer(pid1)
		assert.Equal(t, peer1.Status, StatusCodeBanned)
	})

	t.Run("Testing UpdateLastSent", func(t *testing.T) {
		now := time.Now()
		peerSet.UpdateLastSent(pid1)

		peer1 := peerSet.getPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastSent, now)
	})

	t.Run("Testing UpdateLastReceived", func(t *testing.T) {
		now := time.Now()
		peerSet.UpdateLastReceived(pid1)

		peer1 := peerSet.getPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastReceived, now)
	})

	t.Run("Testing StartedAt", func(t *testing.T) {
		assert.LessOrEqual(t, peerSet.StartedAt(), time.Now())
	})

	t.Run("Testing RemovePeer", func(t *testing.T) {
		peerSet.RemovePeer(peer.ID("unknown"))
		assert.Equal(t, peerSet.Len(), 3)

		peerSet.RemovePeer(peer.ID("peer2"))
		assert.Equal(t, peerSet.Len(), 2)
	})
}

func TestOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pid := peer.ID("peer1")
	session := ps.OpenSession(pid, 100, 101)

	assert.NotNil(t, session)
	from, to := session.Range()
	assert.Equal(t, uint32(100), from)
	assert.Equal(t, uint32(101), to)
	assert.Equal(t, pid, session.PeerID())
	assert.False(t, session.IsUncompleted())
	assert.LessOrEqual(t, session.StartedAt(), time.Now())
	assert.True(t, ps.HasOpenSession(pid))
	assert.False(t, ps.HasOpenSession("peer2"))
	assert.Equal(t, 1, ps.NumberOfOpenSessions())
}

func TestFindSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)
	session := ps.OpenSession("peer1", 100, 101)

	// Test finding an existing session
	foundSession := ps.FindSession(session.SessionID())

	assert.Equal(t, session, foundSession)

	// Test finding a non-existing session
	nonExistingSession := ps.FindSession(999)

	assert.Nil(t, nonExistingSession)
}

func TestNumberOfOpenSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.Equal(t, 0, ps.NumberOfOpenSessions())

	// Test when there are multiple open sessions
	ps.OpenSession("peer1", 100, 101)
	ps.OpenSession("peer2", 200, 201)
	ps.OpenSession("peer3", 300, 301)

	assert.Equal(t, 3, ps.NumberOfOpenSessions())
}

func TestHasAnyOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.False(t, ps.HasAnyOpenSession())

	// Test when there are open sessions
	s := ps.OpenSession("peer1", 100, 101)
	assert.True(t, ps.HasAnyOpenSession())

	ps.RemoveSession(s.SessionID())
	assert.False(t, ps.HasAnyOpenSession())
}

// func TestRemoveSession(t *testing.T) {
// 	ps := NewPeerSet(time.Minute)
// 	session := ps.OpenSession("peer1", 100)

// 	// Test closing an existing session
// 	ps.CloseSession(session.SessionID())

// 	assert.Equal(t, 0, ps.NumberOfOpenSessions())

// 	// Test closing a non-existing session
// 	ps.CloseSession(999)

// 	assert.Equal(t, 0, ps.NumberOfOpenSessions())
// }

func TestUncompletedSession(t *testing.T) {
	ps := NewPeerSet(time.Second)

	pid := peer.ID("peer1")
	session := ps.OpenSession(pid, 100, 101)
	assert.False(t, session.IsUncompleted())

	session.SetUncompleted()
	assert.True(t, session.IsUncompleted())
}

func TestGetRandomWeightedPeer(t *testing.T) {
	// We create 6 peers with varying success and failure counts:
	// peer_1 has 4 successful bundles and 0 failed bundles
	// peer_2 has 4 successful bundles and 1 failed attempt
	// ...
	// peer_6 has 4 successful bundles and 4 failed bundles
	// peer_7 has 4 successful bundles and 5 failed bundles
	peerSet := NewPeerSet(time.Second)
	for i := 0; i < 6; i++ {
		pid := peer.ID(fmt.Sprintf("peer_%v", i+1))
		peerSet.UpdateInfo(pid, fmt.Sprintf("Moniker_%v", i+1), "Agent1", nil, service.New())
		peerSet.UpdateStatus(pid, StatusCodeKnown)

		for s := 0; s < 4; s++ {
			peerSet.IncreaseReceivedBundlesCounter(pid)
		}
		for f := 0; f < i; f++ {
			peerSet.IncreaseInvalidBundlesCounter(pid)
		}
	}

	// Now let's run TestGetRandomPeer for 1000 times

	hits := make(map[peer.ID]int)
	for i := 0; i < 1000; i++ {
		p := peerSet.GetRandomPeer()
		hits[p.PeerID]++
	}

	assert.Greater(t, hits[peer.ID("peer_1")], hits[peer.ID("peer_3")])
	assert.Greater(t, hits[peer.ID("peer_3")], hits[peer.ID("peer_5")])
	assert.Greater(t, hits[peer.ID("peer_5")], 0)
	assert.Greater(t, hits[peer.ID("peer_6")], 0)
}

func TestGetRandomPeerConnected(t *testing.T) {
	peerSet := NewPeerSet(time.Second)

	pidBanned := peer.ID("known")
	pidConnected := peer.ID("connected")
	pidKnown := peer.ID("banned")
	peerSet.UpdateInfo(pidBanned, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidConnected, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidKnown, "moniker", "agent", nil, service.New())

	peerSet.UpdateStatus(pidBanned, StatusCodeBanned)
	peerSet.UpdateStatus(pidConnected, StatusCodeConnected)
	peerSet.UpdateStatus(pidKnown, StatusCodeKnown)

	p := peerSet.GetRandomPeer()

	assert.NotEqual(t, p.PeerID, pidBanned)
	assert.NotEqual(t, p.PeerID, pidConnected)
	assert.Equal(t, p.PeerID, pidKnown)
}

func TestGetRandomPeerOnePeer(t *testing.T) {
	peerSet := NewPeerSet(time.Second)

	pidLonely := peer.ID("banned")
	peerSet.UpdateInfo(pidLonely, "moniker", "agent", nil, service.New())
	peerSet.UpdateStatus(pidLonely, StatusCodeKnown)

	p := peerSet.GetRandomPeer()

	assert.Equal(t, p.PeerID, pidLonely)
}

func TestUpdateAddress(t *testing.T) {
	ps := NewPeerSet(time.Second)

	pid := peer.ID("peer1")
	addr := "pid-1-address"
	ps.UpdateAddress(pid, addr)

	p := ps.GetPeer(pid)
	assert.Equal(t, addr, p.Address)
}
