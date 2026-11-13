package network

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The host accepts connections as soon as it is created, before the other modules
// register their receivers on the pipe. The events of these connections are sent
// when the service starts, in their order, and the next events are sent at once.
func TestNotifeeKeepsTheEventsUntilStart(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	notifee := notifeeFactory(t)()

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	ma, err := multiaddr.NewMultiaddr("/ip4/2.2.2.2/tcp/21888")
	require.NoError(t, err)

	notifee.Connected(nil, &conn{pid: pid1, ma: ma})
	notifee.Disconnected(nil, &conn{pid: pid1, ma: ma})

	events := notifee.networkPipe.UnsafeGetChannel()
	select {
	case evt := <-events:
		require.FailNow(t, "an event was sent before the start", "type %v", evt.Type())
	default:
	}

	notifee.Start()
	t.Cleanup(notifee.Stop)

	first := <-events
	second := <-events
	require.Equal(t, EventTypeConnect, first.Type())
	require.Equal(t, EventTypeDisconnect, second.Type())
	assert.Equal(t, pid1, first.(*ConnectEvent).PeerID)
	assert.Equal(t, pid1, second.(*DisconnectEvent).PeerID)

	notifee.Connected(nil, &conn{pid: pid2, ma: ma})
	third := <-events
	require.Equal(t, EventTypeConnect, third.Type())
	assert.Equal(t, pid2, third.(*ConnectEvent).PeerID)
}
