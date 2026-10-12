package session

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenSession(t *testing.T) {
	mgr := NewManager(time.Minute)
	assert.Zero(t, mgr.NumberOfSessions())
	assert.False(t, mgr.HasAnyOpenSession())

	ssn1 := mgr.OpenSession("peer1", 100, 10)
	ssn2 := mgr.OpenSession("peer2", 110, 20)
	ssn3 := mgr.OpenSession("peer1", 130, 30)

	// Session IDs are unique, and they never repeat.
	assert.NotEqual(t, ssn1.SessionID, ssn2.SessionID)
	assert.NotEqual(t, ssn2.SessionID, ssn3.SessionID)
	assert.NotEqual(t, ssn1.SessionID, ssn3.SessionID)

	assert.Equal(t, Open, ssn1.Status)
	assert.Equal(t, peer.ID("peer1"), ssn1.PeerID)
	assert.Equal(t, uint32(100), uint32(ssn1.From))
	assert.Equal(t, uint32(10), ssn1.Count)
	assert.WithinDuration(t, time.Now(), ssn1.LastActivity, time.Second)

	assert.Equal(t, 3, mgr.NumberOfSessions())
	assert.Len(t, mgr.Sessions(), 3)
	assert.True(t, mgr.HasAnyOpenSession())
	assert.True(t, mgr.HasOpenSession("peer1"))
	assert.True(t, mgr.HasOpenSession("peer2"))
	assert.False(t, mgr.HasOpenSession("peer3"))
}

func TestSessionStatus(t *testing.T) {
	mgr := NewManager(time.Minute)
	ssn1 := mgr.OpenSession("peer1", 100, 10)
	ssn2 := mgr.OpenSession("peer2", 110, 10)
	ssn3 := mgr.OpenSession("peer3", 120, 10)

	assert.Equal(t, Stats{Total: 3, Open: 3}, mgr.Stats())

	require.Same(t, ssn1, mgr.SetSessionCompleted(ssn1.SessionID))
	mgr.SetSessionUncompleted(ssn2.SessionID)

	assert.Equal(t, Completed, ssn1.Status)
	assert.Equal(t, Uncompleted, ssn2.Status)
	assert.Equal(t, Open, ssn3.Status)
	assert.Equal(t, Stats{Total: 3, Open: 1, Completed: 1, Uncompleted: 1}, mgr.Stats())

	// Only open sessions are counted as open.
	assert.False(t, mgr.HasOpenSession("peer1"))
	assert.False(t, mgr.HasOpenSession("peer2"))
	assert.True(t, mgr.HasOpenSession("peer3"))

	mgr.SetSessionCompleted(ssn3.SessionID)
	assert.False(t, mgr.HasAnyOpenSession())
}

func TestUnknownSession(t *testing.T) {
	mgr := NewManager(time.Minute)
	mgr.OpenSession("peer1", 100, 10)

	require.NotPanics(t, func() {
		assert.Nil(t, mgr.SetSessionCompleted(12345))
		mgr.SetSessionUncompleted(12345)
		mgr.UpdateSessionLastActivity(12345)
	})

	assert.Equal(t, Stats{Total: 1, Open: 1}, mgr.Stats())
}

func TestRemoveAllSessions(t *testing.T) {
	mgr := NewManager(time.Minute)
	ssn1 := mgr.OpenSession("peer1", 100, 10)
	mgr.OpenSession("peer2", 110, 10)

	mgr.RemoveAllSessions()

	assert.Zero(t, mgr.NumberOfSessions())
	assert.Empty(t, mgr.Sessions())
	assert.False(t, mgr.HasAnyOpenSession())

	// A removed session never comes back and its ID is not reused.
	ssn3 := mgr.OpenSession("peer1", 100, 10)
	assert.NotEqual(t, ssn1.SessionID, ssn3.SessionID)
	assert.Nil(t, mgr.SetSessionCompleted(ssn1.SessionID))
}

func TestExpiredSessions(t *testing.T) {
	timeout := 50 * time.Millisecond
	mgr := NewManager(timeout)

	open := mgr.OpenSession("peer1", 100, 10)
	completed := mgr.OpenSession("peer2", 110, 10)
	uncompleted := mgr.OpenSession("peer3", 120, 10)
	refreshed := mgr.OpenSession("peer4", 130, 10)
	mgr.SetSessionCompleted(completed.SessionID)
	mgr.SetSessionUncompleted(uncompleted.SessionID)

	time.Sleep(2 * timeout)
	mgr.UpdateSessionLastActivity(refreshed.SessionID)
	mgr.SetExpiredSessionsAsUncompleted()

	assert.Equal(t, Uncompleted, open.Status, "an open session without activity is expired")
	assert.Equal(t, Open, refreshed.Status, "a session with a recent activity is not expired")
	assert.Equal(t, Uncompleted, uncompleted.Status)

	// The blocks of a completed session are downloaded. It must not be requested again,
	// even if the session is old.
	assert.Equal(t, Completed, completed.Status, "a completed session should not become uncompleted")
}

func TestStatsString(t *testing.T) {
	stats := Stats{Total: 4, Open: 1, Completed: 2, Uncompleted: 1}

	assert.Equal(t, "total: 4, open: 1, completed: 2, uncompleted: 1", stats.String())
}
