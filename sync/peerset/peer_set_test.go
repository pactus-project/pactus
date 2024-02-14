package peerset

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPeerSet(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	peerSet := NewPeerSet(time.Minute)

	pk1, _ := ts.RandBLSKeyPair()
	pk2, _ := ts.RandBLSKeyPair()
	pk3, _ := ts.RandBLSKeyPair()
	pk4, _ := ts.RandBLSKeyPair()
	pk5, _ := ts.RandBLSKeyPair()
	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()
	peerSet.UpdateInfo(pid1, "Moniker1", "Agent1",
		[]*bls.PublicKey{pk1, pk2}, service.New(service.Network))
	peerSet.UpdateInfo(pid2, "Moniker2", "Agent2",
		[]*bls.PublicKey{pk3}, service.New(service.None))
	peerSet.UpdateInfo(pid3, "Moniker3", "Agent3",
		[]*bls.PublicKey{pk4, pk5}, service.New(service.Network))

	t.Run("Testing Len", func(t *testing.T) {
		assert.Equal(t, 3, peerSet.Len())
	})

	t.Run("Testing Iterate peers", func(t *testing.T) {
		// Verify that the peer list contains the expected peers
		found := false
		peerSet.IteratePeers(func(p *Peer) bool {
			if p.PeerID == pid2 {
				found = true

				return true
			}

			return false
		})

		assert.True(t, found, "Peer with ID %s not found in the peer list", pid2)
	})

	t.Run("Testing GetPeer", func(t *testing.T) {
		p := peerSet.GetPeer(pid2)
		assert.Equal(t, pid2, p.PeerID)
		assert.Equal(t, StatusCodeUnknown, p.Status)

		p = peerSet.GetPeer(ts.RandPeerID())
		assert.Nil(t, p)
	})

	t.Run("Testing ConsensusKeys", func(t *testing.T) {
		p := peerSet.GetPeer(pid3)

		assert.Contains(t, p.ConsensusKeys, pk4)
		assert.Contains(t, p.ConsensusKeys, pk5)
	})

	t.Run("Testing counters", func(t *testing.T) {
		peerSet.IncreaseInvalidBundlesCounter(pid1)
		peerSet.IncreaseReceivedBundlesCounter(pid1)
		peerSet.IncreaseReceivedBytesCounter(pid1, message.TypeBlocksResponse, 100)
		peerSet.IncreaseReceivedBytesCounter(pid1, message.TypeTransactions, 150)
		peerSet.IncreaseSentCounters(message.TypeBlocksRequest, 200, nil)
		peerSet.IncreaseSentCounters(message.TypeBlocksRequest, 250, &pid1)

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
		assert.Equal(t, peerSet.TotalSentBundles(), 2)
	})

	t.Run("Testing UpdateHeight", func(t *testing.T) {
		height := ts.RandHeight()
		h := ts.RandHash()
		peerSet.UpdateHeight(pid1, height, h)

		peer1 := peerSet.getPeer(pid1)
		assert.Equal(t, height, peer1.Height)
		assert.Equal(t, h, peer1.LastBlockHash)
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
		peerSet.RemovePeer(ts.RandPeerID())
		assert.Equal(t, peerSet.Len(), 3)

		peerSet.RemovePeer(pid2)
		assert.Equal(t, peerSet.Len(), 2)
	})
}

func TestOpenSession(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	ps := NewPeerSet(time.Minute)

	pid := ts.RandPeerID()
	ssn := ps.OpenSession(pid, 100, 1)

	assert.NotNil(t, ssn)
	assert.Equal(t, uint32(100), ssn.From)
	assert.Equal(t, uint32(1), ssn.Count)
	assert.Equal(t, pid, ssn.PeerID)
	assert.Equal(t, session.Open, ssn.Status)
	assert.LessOrEqual(t, ssn.LastActivity, time.Now())
	assert.True(t, ps.HasOpenSession(pid))
	assert.False(t, ps.HasOpenSession(ts.RandPeerID()))
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.Contains(t, ps.Sessions(), ssn)
}

func TestFindSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)
	ssn := ps.OpenSession("peer1", 100, 101)

	// Test finding an existing session
	foundSsn := ps.FindSession(ssn.SessionID)

	assert.Equal(t, ssn, foundSsn)

	// Test finding a non-existing session
	nonExistingSsn := ps.FindSession(999)

	assert.Nil(t, nonExistingSsn)
}

func TestNumberOfSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.Equal(t, 0, ps.NumberOfSessions())

	// Test when there are multiple open sessions
	ps.OpenSession("peer1", 100, 101)
	ps.OpenSession("peer2", 200, 201)
	ps.OpenSession("peer3", 300, 301)

	assert.Equal(t, 3, ps.NumberOfSessions())
}

func TestHasAnyOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.False(t, ps.HasAnyOpenSession())

	ssn := ps.OpenSession("peer1", 100, 101)
	assert.True(t, ps.HasAnyOpenSession())

	ps.SetSessionCompleted(ssn.SessionID)
	assert.False(t, ps.HasAnyOpenSession())
}

func TestRemoveAllSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	_ = ps.OpenSession("peer1", 100, 101)
	_ = ps.OpenSession("peer2", 100, 101)
	_ = ps.OpenSession("peer3", 100, 101)

	ps.RemoveAllSessions()
	assert.Zero(t, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
}

func TestCompletedSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	ssn := ps.OpenSession("peer1", 100, 101)
	assert.Equal(t, session.Open, ssn.Status)

	ps.SetSessionCompleted(ssn.SessionID)
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
	assert.Equal(t, session.Completed, ssn.Status)
}

func TestUncompletedSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	ssn := ps.OpenSession("peer1", 100, 101)
	assert.Equal(t, session.Open, ssn.Status)

	ps.SetSessionUncompleted(ssn.SessionID)
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
	assert.Equal(t, session.Uncompleted, ssn.Status)
}

func TestExpireSessions(t *testing.T) {
	timeout := 100 * time.Millisecond
	ps := NewPeerSet(timeout)

	ssn := ps.OpenSession("peer1", 100, 101)
	time.Sleep(timeout)

	ps.SetExpiredSessionsAsUncompleted()
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
	assert.Equal(t, session.Uncompleted, ssn.Status)
}

func TestGetRandomWeightedPeer(t *testing.T) {
	// We create 6 peers with varying success and failure counts:
	// peer_1 has score 6/6 (completed sessions / total sessions)
	// peer_2 has score 5/6
	// ...
	// peer_6 has score 1/5
	peerSet := NewPeerSet(time.Minute)
	for i := 0; i < 6; i++ {
		pid := peer.ID(fmt.Sprintf("peer_%v", i+1))
		peerSet.UpdateInfo(pid, fmt.Sprintf("Moniker_%v", i+1), "Agent1", nil, service.New())
		peerSet.UpdateStatus(pid, StatusCodeKnown)

		for r := 0; r < 6; r++ {
			ssn := peerSet.OpenSession(pid, 0, 0)

			if r < 6-i {
				peerSet.SetSessionCompleted(ssn.SessionID)
			}
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
	assert.GreaterOrEqual(t, hits[peer.ID("peer_5")], 0)
	assert.GreaterOrEqual(t, hits[peer.ID("peer_6")], 0)
}

func TestGetRandomPeerConnected(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pidBanned := peer.ID("known")
	pidConnected := peer.ID("connected")
	pidDisconnected := peer.ID("disconnected")
	pidKnown := peer.ID("banned")
	peerSet.UpdateInfo(pidBanned, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidConnected, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidDisconnected, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidKnown, "moniker", "agent", nil, service.New())

	peerSet.UpdateStatus(pidBanned, StatusCodeBanned)
	peerSet.UpdateStatus(pidConnected, StatusCodeConnected)
	peerSet.UpdateStatus(pidDisconnected, StatusCodeDisconnected)
	peerSet.UpdateStatus(pidKnown, StatusCodeKnown)

	p := peerSet.GetRandomPeer()

	assert.NotEqual(t, p.PeerID, pidBanned)
	assert.NotEqual(t, p.PeerID, pidConnected)
	assert.NotEqual(t, p.PeerID, pidDisconnected)
	assert.Equal(t, p.PeerID, pidKnown)
}

func TestGetRandomPeerOnePeer(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pidLonely := peer.ID("banned")
	peerSet.UpdateInfo(pidLonely, "moniker", "agent", nil, service.New())
	peerSet.UpdateStatus(pidLonely, StatusCodeKnown)

	p := peerSet.GetRandomPeer()

	assert.Equal(t, p.PeerID, pidLonely)
}

func TestUpdateAddress(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pid := peer.ID("peer1")
	addr := "pid-1-address"
	dir := "Inbound"
	ps.UpdateAddress(pid, addr, dir)

	p := ps.GetPeer(pid)
	assert.Equal(t, addr, p.Address)
	assert.Equal(t, dir, p.Direction)
}

func TestUpdateSessionLastActivity(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	ssn := ps.OpenSession("peer1", 100, 101)
	activity1 := ssn.LastActivity
	time.Sleep(10 * time.Millisecond)
	ps.UpdateSessionLastActivity(ssn.SessionID)
	assert.Greater(t, ssn.LastActivity, activity1)
}

func TestUpdateProtocols(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pid := peer.ID("peer-1")
	protocols := []string{"protocol-1"}
	ps.UpdateProtocols(pid, protocols)

	p := ps.GetPeer(pid)
	assert.Equal(t, p.Protocols, protocols)
}
