package network

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A receiver that never reads nor closes a stream must not keep a goroutine
// of the sender alive after the stream timeout.
func TestSendToDoesNotLeakGoroutines(t *testing.T) {
	confA := testConfig()
	confA.StreamTimeout = 300 * time.Millisecond
	confA.CheckConnectivityInterval = 60 * time.Second
	networkA := makeTestNetwork(t, confA, nil)

	confB := testConfig()
	confB.StreamTimeout = 300 * time.Millisecond
	confB.CheckConnectivityInterval = 60 * time.Second
	require.NoError(t, util.WriteFile(confB.PeerStorePath,
		[]byte(fmt.Sprintf("[\"/ip4/127.0.0.1/tcp/%v/p2p/%v\"]",
			confA.DefaultPort, networkA.SelfID().String()))))
	networkB := makeTestNetwork(t, confB, nil)

	require.Eventually(t, func() bool {
		return len(networkB.host.Network().ConnsToPeer(networkA.SelfID())) > 0
	}, 5*time.Second, 50*time.Millisecond)

	// The other tests of the package, and the two networks of this one, start and stop
	// goroutines in the background. Wait until their number is stable before counting.
	before := runtime.NumGoroutine()
	stableSince := time.Now()
	require.Eventually(t, func() bool {
		if now := runtime.NumGoroutine(); now != before {
			before = now
			stableSince = time.Now()
		}

		return time.Since(stableSince) > 500*time.Millisecond
	}, 10*time.Second, 50*time.Millisecond)

	// Nobody consumes the events of network A: its streams are never read or closed.
	// A stream that leaks keeps two goroutines.
	const streams = 40
	for range streams {
		_, err := networkB.stream.SendTo([]byte("hello"), networkA.SelfID())
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return runtime.NumGoroutine() < before+streams
	}, 5*time.Second, 100*time.Millisecond, "the stream goroutines should end after the timeout")
}
