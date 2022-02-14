package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	net1, net2 := setup(t, TestConfig(), TestConfig())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	for {
		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
			break
		}
	}

	msg := []byte("test")
	require.NoError(t, net1.SendTo(msg, net2.SelfID()))

	e := shouldReceiveEvent(t, net2).(*StreamMessage)
	buf := make([]byte, 4)
	_, err := e.Reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, e.Source, net1.SelfID())
	assert.Equal(t, buf, msg)

	net1.Stop()
	net2.Stop()
}
