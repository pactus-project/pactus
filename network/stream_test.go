package network

// func TestCloseStream(t *testing.T) {
// 	confA := testConfig()
// 	confA.StreamTimeout = 1 * time.Second // Reduce timeout for testing
// 	confA.EnableMdns = true
// 	networkA := makeTestNetwork(t, confA, nil)

// 	confB := testConfig()
// 	confB.EnableMdns = true
// 	networkB := makeTestNetwork(t, confB, nil)

// 	assert.EventuallyWithT(t, func(c *assert.CollectT) {
// 		e := <-networkA.EventChannel()
// 		assert.Equal(c, EventTypeConnect, e.Type())
// 	}, 5*time.Second, 100*time.Millisecond)

// 	assert.EventuallyWithT(t, func(c *assert.CollectT) {
// 		e := <-networkB.EventChannel()
// 		assert.Equal(c, EventTypeConnect, e.Type())
// 	}, 5*time.Second, 100*time.Millisecond)

// 	t.Run("Stream timeout", func(t *testing.T) {
// 		stream, err := networkA.stream.SendRequest([]byte("test-1"), networkB.SelfID())
// 		require.NoError(t, err)

// 		// NetworkB doesn't close the stream.
// 		assert.EventuallyWithT(t, func(c *assert.CollectT) {
// 			e := <-networkB.EventChannel()
// 			_, ok := e.(*StreamMessage)
// 			assert.True(c, ok)
// 		}, 5*time.Second, 100*time.Millisecond)

// 		// Wait fot the steam timeout.
// 		time.Sleep(2 * confA.StreamTimeout)

// 		_, err = stream.Write([]byte("should-be-closed"))
// 		assert.ErrorContains(t, err, "write on closed stream")
// 	})

// 	t.Run("Stream closed", func(t *testing.T) {
// 		stream, err := networkA.stream.SendRequest([]byte("test-2"), networkB.SelfID())
// 		require.NoError(t, err)

// 		// NetworkB close the stream.
// 		assert.EventuallyWithT(t, func(c *assert.CollectT) {
// 			e := <-networkB.EventChannel()
// 			s, ok := e.(*StreamMessage)
// 			assert.True(c, ok)

// 			if ok {
// 				err := s.Reader.Close()
// 				assert.NoError(t, err)
// 			}
// 		}, 5*time.Second, 100*time.Millisecond)

// 		_, err = stream.Write([]byte("should-be-closed"))
// 		assert.ErrorContains(t, err, "write on closed stream")
// 	})
// }
