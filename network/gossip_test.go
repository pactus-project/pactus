package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinConsensusTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-consensus-topic")

	require.ErrorIs(t, net.Broadcast(msg, TopicIDConsensus),
		NotSubscribedError{
			TopicID: TopicIDConsensus,
		})
	require.NoError(t, net.JoinConsensusTopic())
	require.NoError(t, net.Broadcast(msg, TopicIDConsensus))
}

func TestInvalidTopic(t *testing.T) {
	net, err := NewNetwork("test", testConfig())
	assert.NoError(t, err)

	msg := []byte("test-invalid-topic")

	require.ErrorIs(t, net.Broadcast(msg, -1),
		InvalidTopicError{
			TopicID: TopicID(-1),
		})
}
