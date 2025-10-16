package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloseStream(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	makeNetworks := func(streamTimeout time.Duration) (networkA, networkB *network) {
		confA := testConfig()
		confA.StreamTimeout = streamTimeout
		confA.CheckConnectivityInterval = 60 * time.Second
		networkA = makeTestNetwork(t, confA, nil)

		confB := testConfig()
		confB.StreamTimeout = streamTimeout
		confB.CheckConnectivityInterval = 60 * time.Second
		_ = util.WriteFile(confB.PeerStorePath,
			[]byte(fmt.Sprintf("[\"/ip4/127.0.0.1/tcp/%v/p2p/%v\"]",
				confA.DefaultPort, networkA.SelfID().String())))
		networkB = makeTestNetwork(t, confB, nil)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			e := <-networkA.networkPipe.UnsafeGetChannel()
			assert.Equal(collect, EventTypeConnect, e.Type())
		}, 5*time.Second, 100*time.Millisecond)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			assert.Equal(collect, EventTypeConnect, e.Type())
		}, 5*time.Second, 100*time.Millisecond)

		return networkA, networkB
	}

	// streamClosed check if a stream is fully closed for both ends
	streamClosed := func(networkA, networkB *network, streamID string) bool {
		connsAtoB := networkA.host.Network().ConnsToPeer(networkB.SelfID())
		streamsAtoB := connsAtoB[0].GetStreams()

		hasStream := false
		for _, s := range streamsAtoB {
			if s.ID() == streamID {
				hasStream = true
			}
		}

		connsBtoA := networkB.host.Network().ConnsToPeer(networkA.SelfID())
		streamsBtoA := connsBtoA[0].GetStreams()
		for _, s := range streamsBtoA {
			if s.ID() == streamID {
				hasStream = true
			}
		}

		return !hasStream
	}

	t.Run("Normal Case", func(t *testing.T) {
		networkA, networkB := makeNetworks(10 * time.Second)
		sentMsg := ts.RandBytes(32)
		stream, err := networkA.stream.SendTo(sentMsg, networkB.SelfID())
		require.NoError(t, err)

		// NetworkB close the stream after reading the data
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			streamMsg, ok := e.(*StreamMessage)
			assert.True(collect, ok)
			if ok {
				receivedMsg := make([]byte, len(sentMsg))
				_, _ = streamMsg.Reader.Read(receivedMsg)
				assert.Equal(collect, sentMsg, receivedMsg)

				_ = streamMsg.Reader.Close()
			}
		}, 5*time.Second, 100*time.Millisecond)

		// NetworkA should receive EOF and close/remove the stream
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			assert.True(collect, streamClosed(networkA, networkB, stream.ID()))
		}, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("Receiver (NetworkB) doesn't close the stream", func(t *testing.T) {
		networkA, networkB := makeNetworks(1 * time.Second)
		sentMsg := ts.RandBytes(32)
		stream, err := networkA.stream.SendTo(sentMsg, networkB.SelfID())
		require.NoError(t, err)

		// NetworkB close the stream after reading the data.
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			s, ok := e.(*StreamMessage)
			assert.True(collect, ok)
			if ok {
				receivedMsg := make([]byte, len(sentMsg))
				_, _ = s.Reader.Read(receivedMsg)
				assert.Equal(collect, sentMsg, receivedMsg)

				// NetworkB doesn't close the stream.
				// s.Reader.Close()
			}
		}, 5*time.Second, 100*time.Millisecond)

		// NetworkA should close/remove the stream after timeout.
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			assert.True(collect, streamClosed(networkA, networkB, stream.ID()))
		}, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("Receiver (NetworkB) close the stream without reading", func(t *testing.T) {
		networkA, networkB := makeNetworks(1 * time.Second)
		sentMsg := ts.RandBytes(32)
		stream, err := networkA.stream.SendTo(sentMsg, networkB.SelfID())
		require.NoError(t, err)

		// NetworkB close the stream after reading the data
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			e := <-networkB.networkPipe.UnsafeGetChannel()
			s, ok := e.(*StreamMessage)
			assert.True(collect, ok)
			if ok {
				_ = s.Reader.Close()
			}
		}, 5*time.Second, 100*time.Millisecond)

		// NetworkA should close/remove the stream as well.
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			assert.True(collect, streamClosed(networkA, networkB, stream.ID()))
		}, 5*time.Second, 100*time.Millisecond)
	})
}
