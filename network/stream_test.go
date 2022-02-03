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
	setup(t)

	assert.NoError(t, tNet1.Start())
	assert.NoError(t, tNet2.Start())

	received := make(chan bool)
	msg := []byte("test")
	cb := func(r io.Reader, id peer.ID) {
		buf := make([]byte, 4)
		_, err := r.Read(buf)
		assert.Equal(t, id, tNet1.SelfID())
		assert.Equal(t, buf, msg)
		assert.NoError(t, err)
		received <- true
	}
	go tNet2.SetCallback(cb)

	assert.NoError(t, tNet1.Start())
	assert.NoError(t, tNet2.Start())

	for {
		if tNet1.NumConnectedPeers() > 0 && tNet2.NumConnectedPeers() > 0 {
			break
		}
	}

	time.Sleep(1 * time.Second)
	require.NoError(t, tNet1.SendMessage(msg, tNet2.SelfID()))

	<-received

	tNet1.Stop()
	tNet2.Stop()
}
