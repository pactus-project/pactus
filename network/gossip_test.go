package network

import (
	"context"
	"testing"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2pcore "github.com/libp2p/go-libp2p/core"
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

func TestTopicValidator(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)

	selfID := net.host.ID()
	propagate := false
	validator := net.gossip.createValidator(TopicIDConsensus,
		func(_ *GossipMessage) bool { return propagate })

	tests := []struct {
		name           string
		peerID         lp2pcore.PeerID
		propagate      bool
		expectedResult lp2pps.ValidationResult
	}{
		{
			name:           "Message from self",
			propagate:      false,
			peerID:         selfID,
			expectedResult: lp2pps.ValidationAccept,
		},
		{
			name:           "Message from other peer, should not propagate",
			propagate:      false,
			peerID:         "other-peerID",
			expectedResult: lp2pps.ValidationIgnore,
		},
		{
			name:           "Message from other peer, should propagate",
			propagate:      true,
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
			propagate = tt.propagate
			result := validator(context.Background(), tt.peerID, msg)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
