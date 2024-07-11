package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinConsensusTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-consensus-topic")

	assert.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDConsensus),
		NotSubscribedError{
			TopicID: TopicIDConsensus,
		})
	assert.NoError(t, net.JoinTopic(TopicIDConsensus, alwaysPropagate))
	assert.NoError(t, net.gossip.Broadcast(msg, TopicIDConsensus))
}

func TestJoinInvalidTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	assert.ErrorIs(t, net.JoinTopic(TopicIDUnspecified, alwaysPropagate),
		InvalidTopicError{
			TopicID: TopicIDUnspecified,
		})

	assert.ErrorIs(t, net.JoinTopic(TopicID(-1), alwaysPropagate),
		InvalidTopicError{
			TopicID: TopicID(-1),
		})
}

func TestInvalidTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-invalid-topic")

	assert.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDUnspecified),
		InvalidTopicError{
			TopicID: TopicIDUnspecified,
		})

	assert.ErrorIs(t, net.gossip.Broadcast(msg, -1),
		InvalidTopicError{
			TopicID: TopicID(-1),
		})
}
