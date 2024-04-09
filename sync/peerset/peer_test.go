package peerset

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/stretchr/testify/assert"
)

func TestPeerStatus(t *testing.T) {
	testCases := []struct {
		status      StatusCode
		isBanned    bool
		isConnected bool
		isKnown     bool
	}{
		{StatusCodeBanned, true, false, false},
		{StatusCodeDisconnected, false, false, false},
		{StatusCodeConnected, false, true, false},
		{StatusCodeKnown, false, true, true},
		{StatusCodeTrusty, false, true, true},
	}

	for _, testCase := range testCases {
		p := NewPeer(peer.ID("test"))
		p.Status = testCase.status

		assert.Equal(t, testCase.isBanned, p.IsBanned())
		assert.Equal(t, testCase.isConnected, p.IsConnected())
		assert.Equal(t, testCase.isKnown, p.IsKnownOrTrusty())
	}
}

func TestHasNetworkService(t *testing.T) {
	p1 := NewPeer(peer.ID("peer-1"))
	p2 := NewPeer(peer.ID("peer-1"))
	p2.Services = service.New(service.Network)

	assert.False(t, p1.HasNetworkService())
	assert.True(t, p2.HasNetworkService())
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
		p := NewPeer(peer.ID("test"))
		p.TotalSessions = testCase.totalSession
		p.CompletedSessions = testCase.completedSession

		score := p.DownloadScore()
		assert.Equal(t, testCase.expectedScore, score,
			"Test %v failed. expected score %d, got %d",
			i+1, testCase.expectedScore, score)
	}
}
