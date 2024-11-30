package network

import (
	"context"
	"testing"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	"github.com/stretchr/testify/assert"
)

func TestJoinBlockTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-block-topic")

	assert.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDBlock),
		NotSubscribedError{
			TopicID: TopicIDBlock,
		})
	assert.NoError(t, net.JoinTopic(TopicIDBlock, alwaysPropagate))
	assert.NoError(t, net.gossip.Broadcast(msg, TopicIDBlock))

	assert.Error(t, net.JoinTopic(TopicIDBlock, alwaysPropagate), "already joined")
}

func TestJoinConsensusTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-consensus-topic")

	assert.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDConsensus),
		NotSubscribedError{
			TopicID: TopicIDConsensus,
		})
	assert.NoError(t, net.JoinTopic(TopicIDConsensus, alwaysPropagate))
	assert.NoError(t, net.gossip.Broadcast(msg, TopicIDConsensus))

	assert.Error(t, net.JoinTopic(TopicIDConsensus, alwaysPropagate), "already joined")
}

func TestJoinTransactionTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	msg := []byte("test-transaction-topic")

	assert.ErrorIs(t, net.gossip.Broadcast(msg, TopicIDTransaction),
		NotSubscribedError{
			TopicID: TopicIDTransaction,
		})
	assert.NoError(t, net.JoinTopic(TopicIDTransaction, alwaysPropagate))
	assert.NoError(t, net.gossip.Broadcast(msg, TopicIDTransaction))

	assert.Error(t, net.JoinTopic(TopicIDTransaction, alwaysPropagate), "already joined")
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
			result := validator(context.Background(), tt.peerID, msg)
			assert.Equal(t, result, tt.expectedResult)
		})
	}
}
