package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSub(t *testing.T) {
	msg := []byte("test-general-topic")

	require.NoError(t, tNetwork1.Broadcast(msg, TopicIDGeneral))

	e := shouldReceiveEvent(t, tNetwork2).(*GossipMessage)
	assert.Equal(t, e.Source, tNetwork1.SelfID())
	assert.Equal(t, e.Data, msg)
}

func TestJoinTopic(t *testing.T) {
	msg := []byte("test-consensus-topic")

	require.Error(t, tNetwork1.Broadcast(msg, TopicIDConsensus))
	require.NoError(t, tNetwork1.JoinConsensusTopic())
	require.NoError(t, tNetwork1.Broadcast(msg, TopicIDConsensus))
}

func TestInvalidTopic(t *testing.T) {
	msg := []byte("test-invalid-topic")

	require.Error(t, tNetwork1.Broadcast(msg, -1))
}
