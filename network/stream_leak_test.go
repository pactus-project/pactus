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

	// Nobody consumes the events of network A: its streams are never read or closed.
	before := runtime.NumGoroutine()
	for range 40 {
		_, err := networkB.stream.SendTo([]byte("hello"), networkA.SelfID())
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return runtime.NumGoroutine() <= before+10
	}, 5*time.Second, 100*time.Millisecond, "the stream goroutines should end after the timeout")
}
