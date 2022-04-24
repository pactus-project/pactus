package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSub(t *testing.T) {
	msg := []byte("test-general-topic")

	require.NoError(t, tNetworks[0].Broadcast(msg, TopicIDGeneral))

	e := shouldReceiveEvent(t, tNetworks[1]).(*GossipMessage)
	assert.Equal(t, e.Source, tNetworks[0].SelfID())
	assert.Equal(t, e.Data, msg)
}

func TestJoinTopic(t *testing.T) {
	msg := []byte("test-consensus-topic")

	require.Error(t, tNetworks[0].Broadcast(msg, TopicIDConsensus))
	require.NoError(t, tNetworks[0].JoinConsensusTopic())
	require.NoError(t, tNetworks[0].Broadcast(msg, TopicIDConsensus))
}

func TestInvalidTopic(t *testing.T) {
	msg := []byte("test-invalid-topic")

	require.Error(t, tNetworks[0].Broadcast(msg, -1))
}
