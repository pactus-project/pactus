package network

import (
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSub(t *testing.T) {
	net1, net2 := setup(t, TestConfig(), TestConfig())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	received := make(chan bool)
	msg := []byte("test")
	cb := func(r io.Reader, id peer.ID) {
		buf := make([]byte, 4)
		_, err := r.Read(buf)
		assert.Equal(t, id, net1.SelfID())
		assert.Equal(t, buf, msg)
		assert.NoError(t, err)
		received <- true
	}
	go net2.SetCallback(cb)

	assert.NoError(t, net1.JoinGeneralTopic())
	assert.NoError(t, net2.JoinGeneralTopic())
	assert.NoError(t, net1.JoinConsensusTopic())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	for {
		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
			break
		}
	}

	time.Sleep(1 * time.Second)

	require.NoError(t, net1.BroadcastMessage([]byte("should not broadcasted"), ConsensusTopic))
	require.NoError(t, net1.BroadcastMessage(msg, GeneralTopic))

	<-received

	net1.Stop()
	net2.Stop()
}
