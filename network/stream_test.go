package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloseStream(t *testing.T) {
	confA := testConfig()
	confA.StreamTimeout = 1 * time.Second // Reduce timeout for testing
	confA.EnableUDP = true
	confA.EnableMdns = true
	networkA := makeTestNetwork(t, confA, nil)

	confB := testConfig()
	confB.EnableUDP = true
	confB.EnableMdns = true
	confB.StreamTimeout = 1 * time.Second
	confB.BootstrapAddrStrings = []string{
		fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%v", confA.DefaultPort, networkA.SelfID().String()),
		fmt.Sprintf("/ip4/127.0.0.1/udp/%v/quic-v1/p2p/%v", confA.DefaultPort, networkA.SelfID().String()),
	}
	networkB := makeTestNetwork(t, confB, nil)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		e := <-networkA.networkPipe.UnsafeGetChannel()
		assert.Equal(c, EventTypeConnect, e.Type())
	}, 5*time.Second, 100*time.Millisecond)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		e := <-networkB.networkPipe.UnsafeGetChannel()
		assert.Equal(c, EventTypeConnect, e.Type())
	}, 5*time.Second, 100*time.Millisecond)

	t.Run("Stream timeout", func(t *testing.T) {
		stream, err := networkA.stream.SendTo([]byte("test-1"), networkB.SelfID())
		require.NoError(t, err)

		// NetworkB doesn't close the stream.
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			_, ok := e.(*StreamMessage)
			assert.True(c, ok)
		}, 5*time.Second, 100*time.Millisecond)

		// Wait fot the steam timeout.
		time.Sleep(2 * confA.StreamTimeout)

		_, err = stream.Write([]byte("should-be-closed"))
		// The error can be either "stream closed" (from LibP2P)
		// or "write on closed stream" (from QUIC-UDP).
		assert.ErrorContains(t, err, "closed")
	})

	t.Run("Stream closed", func(t *testing.T) {
		stream, err := networkA.stream.SendTo([]byte("test-2"), networkB.SelfID())
		require.NoError(t, err)

		// NetworkB close the stream.
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			s, ok := e.(*StreamMessage)
			assert.True(c, ok)

			if ok {
				err := s.Reader.Close()
				assert.NoError(t, err)
			}
		}, 5*time.Second, 100*time.Millisecond)

		_, err = stream.Write([]byte("should-be-closed"))
		assert.ErrorContains(t, err, "closed")
	})

	// TODO: test for stream reset
	// network.ErrReset
}
