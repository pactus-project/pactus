package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinConsensusTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-consensus-topic")

	require.Error(t, net.Broadcast(msg, TopicIDConsensus))
	require.NoError(t, net.JoinConsensusTopic())
	require.NoError(t, net.Broadcast(msg, TopicIDConsensus))
}

func TestInvalidTopic(t *testing.T) {
	net, err := NewNetwork(testConfig())
	assert.NoError(t, err)

	msg := []byte("test-invalid-topic")

	require.Error(t, net.Broadcast(msg, -1))
}
