package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
)

func TestPubSub(t *testing.T) {
	size := 6
	nets := setup(t, size)
	msg := []byte("test-general-topic")

	require.NoError(t, nets[0].Broadcast(msg, TopicIDGeneral))

	i := util.RandInt(size-1) + 1
	e := shouldReceiveEvent(t, nets[i]).(*GossipMessage)
	assert.Equal(t, e.Source, nets[0].SelfID())
	assert.Equal(t, e.Data, msg)
}

func TestJoinTopic(t *testing.T) {
	nets := setup(t, 1)
	msg := []byte("test-consensus-topic")

	require.Error(t, nets[0].Broadcast(msg, TopicIDConsensus))
	require.NoError(t, nets[0].JoinConsensusTopic())
	require.NoError(t, nets[0].Broadcast(msg, TopicIDConsensus))
}

func TestInvalidTopic(t *testing.T) {
	nets := setup(t, 1)
	msg := []byte("test-invalid-topic")

	require.Error(t, nets[0].Broadcast(msg, -1))
}
