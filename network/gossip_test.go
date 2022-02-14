package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSub(t *testing.T) {
	net1, net2 := setup(t, TestConfig(), TestConfig())

	assert.NoError(t, net1.Start())
	assert.NoError(t, net2.Start())

	assert.NoError(t, net1.JoinGeneralTopic())
	assert.NoError(t, net2.JoinGeneralTopic())
	assert.NoError(t, net1.JoinConsensusTopic())

	for {
		if net1.NumConnectedPeers() > 0 && net2.NumConnectedPeers() > 0 {
			break
		}
	}

	// TODO: Can we remove timer and run tests successfully?
	time.Sleep(1 * time.Second)
	msg := []byte("test")

	require.NoError(t, net1.Broadcast([]byte("should not broadcasted"), TopicIDConsensus))
	require.NoError(t, net1.Broadcast(msg, TopicIDGeneral))

	e := shouldReceiveEvent(t, net2).(*GossipMessage)
	assert.Equal(t, e.Source, net1.SelfID())
	assert.Equal(t, e.Data, msg)

	net1.Stop()
	net2.Stop()
}
