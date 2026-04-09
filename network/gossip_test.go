package network

import (
	"testing"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinBlockTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-block-topic")

	require.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDBlock),
		NotSubscribedError{
			TopicID: TopicIDBlock,
		})
	require.NoError(t, net.JoinTopic(TopicIDBlock, alwaysPropagate))
	require.NoError(t, net.gossip.Broadcast(msg, TopicIDBlock))

	require.Error(t, net.JoinTopic(TopicIDBlock, alwaysPropagate), "already joined")
}

func TestJoinConsensusTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-consensus-topic")

	require.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDConsensus),
		NotSubscribedError{
			TopicID: TopicIDConsensus,
		})
	require.NoError(t, net.JoinTopic(TopicIDConsensus, alwaysPropagate))
	require.NoError(t, net.gossip.Broadcast(msg, TopicIDConsensus))

	require.Error(t, net.JoinTopic(TopicIDConsensus, alwaysPropagate), "already joined")
}

func TestJoinTransactionTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-transaction-topic")

	require.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDTransaction),
		NotSubscribedError{
			TopicID: TopicIDTransaction,
		})
	require.NoError(t, net.JoinTopic(TopicIDTransaction, alwaysPropagate))
	require.NoError(t, net.gossip.Broadcast(msg, TopicIDTransaction))

	require.Error(t, net.JoinTopic(TopicIDTransaction, alwaysPropagate), "already joined")
}

func TestJoinInvalidTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	require.ErrorIs(t, net.JoinTopic(TopicIDUnspecified, alwaysPropagate),
		InvalidTopicError{
			TopicID: TopicIDUnspecified,
		})

	require.ErrorIs(t, net.JoinTopic(TopicID(-1), alwaysPropagate),
		InvalidTopicError{
			TopicID: TopicID(-1),
		})
}

func TestInvalidTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-invalid-topic")

	require.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDUnspecified),
		InvalidTopicError{
			TopicID: TopicIDUnspecified,
		})

	require.ErrorIs(t, net.gossip.Broadcast(msg, -1),
		InvalidTopicError{
			TopicID: TopicID(-1),
		})
}

func TestTopicValidator(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	selfID := net.host.ID()
	propagate := Drop
	validator := net.gossip.createValidator(TopicIDConsensus,
		func(_ *GossipMessage) PropagationPolicy { return propagate })

	tests := []struct {
		name           string
		peerID         lp2pcore.PeerID
		policy         PropagationPolicy
		expectedResult lp2pps.ValidationResult
	}{
		{
			name:           "Message from self",
			policy:         Drop,
			peerID:         selfID,
			expectedResult: lp2pps.ValidationAccept,
		},
		{
			name:           "Message from self",
			policy:         DropButConsume,
			peerID:         selfID,
			expectedResult: lp2pps.ValidationAccept,
		},
		{
			name:           "Message from self",
			policy:         propagate,
			peerID:         selfID,
			expectedResult: lp2pps.ValidationAccept,
		},
		{
			name:           "Message from other peer, should not propagate",
			policy:         Drop,
			peerID:         "other-peerID",
			expectedResult: lp2pps.ValidationIgnore,
		},
		{
			name:           "Message from other peer, should not propagate",
			policy:         DropButConsume,
			peerID:         "other-peerID",
			expectedResult: lp2pps.ValidationIgnore,
		},
		{
			name:           "Message from other peer, should propagate",
			policy:         Propagate,
			peerID:         "other-peerID",
			expectedResult: lp2pps.ValidationAccept,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &lp2pps.Message{
				Message: &pubsubpb.Message{
					Data: []byte("some-data"),
				},
			}
			propagate = tt.policy
			result := validator(t.Context(), tt.peerID, msg)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
