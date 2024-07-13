package peer

import (
	"testing"

	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/stretchr/testify/assert"
)

func TestPeerStatus(t *testing.T) {
	testCases := []struct {
		status             status.Status
		isDisconnected     bool
		isBanned           bool
		isConnected        bool
		isKnown            bool
		isKnownOrConnected bool
	}{
		{status: status.StatusBanned, isBanned: true},
		{status: status.StatusDisconnected, isDisconnected: true},
		{status: status.StatusConnected, isConnected: true, isKnownOrConnected: true},
		{status: status.StatusKnown, isKnown: true, isKnownOrConnected: true},
	}

	for _, testCase := range testCases {
		p := NewPeer("test")
		p.Status = testCase.status

		assert.Equal(t, testCase.isDisconnected, p.Status.IsDisconnected())
		assert.Equal(t, testCase.isBanned, p.Status.IsBanned())
		assert.Equal(t, testCase.isConnected, p.Status.IsConnected())
		assert.Equal(t, testCase.isKnown, p.Status.IsKnown())
		assert.Equal(t, testCase.isKnownOrConnected, p.Status.IsConnectedOrKnown())
	}
}

func TestIsFullNode(t *testing.T) {
	p1 := NewPeer("peer-1")
	p2 := NewPeer("peer-1")
	p1.Services = service.New(service.PrunedNode)
	p2.Services = service.New(service.FullNode)

	assert.False(t, p1.IsFullNode())
	assert.True(t, p2.IsFullNode())
}

func TestDownloadScore(t *testing.T) {
	testCases := []struct {
		totalSession     int
		completedSession int
		expectedScore    int
	}{
		{0, 0, 100},
		{1, 1, 100},
		{1, 0, 50},
		{2, 0, 33},
		{3, 0, 25},
		{5, 4, 83},
		{6, 5, 85},
		{7, 6, 87},
	}

	for i, testCase := range testCases {
		p := NewPeer("test")
		p.TotalSessions = testCase.totalSession
		p.CompletedSessions = testCase.completedSession

		score := p.DownloadScore()
		assert.Equal(t, testCase.expectedScore, score,
			"Test %v failed. expected score %d, got %d",
			i+1, testCase.expectedScore, score)
	}
}
