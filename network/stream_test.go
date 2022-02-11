package network

import (
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	net1, net2 := setup(t, TestConfig(), TestConfig())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	received := make(chan bool)
	msg := []byte("test")
	cb := func(r io.Reader, source peer.ID, from peer.ID) {
		buf := make([]byte, 4)
		_, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, source, net1.SelfID())
		assert.Equal(t, from, net1.SelfID())
		assert.Equal(t, buf, msg)

		received <- true
	}
	go net2.SetCallback(cb)

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	for {
		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
			break
		}
	}

	time.Sleep(1 * time.Second)
	require.NoError(t, net1.SendTo(msg, net2.SelfID()))

	<-received

	net1.Stop()
	net2.Stop()
}
