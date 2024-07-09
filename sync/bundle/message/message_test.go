package message

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	testCases := []struct {
		msgType         Type
		typeName        string
		topicID         network.TopicID
		shouldBroadcast bool
	}{
		{TypeHello, "hello", network.TopicIDUnspecified, false},
		{TypeHelloAck, "hello-ack", network.TopicIDUnspecified, false},
		{TypeTransaction, "transaction", network.TopicIDTransaction, true},
		{TypeQueryProposal, "query-proposal", network.TopicIDConsensus, true},
		{TypeProposal, "proposal", network.TopicIDConsensus, true},
		{TypeQueryVote, "query-vote", network.TopicIDConsensus, true},
		{TypeVote, "vote", network.TopicIDConsensus, true},
		{TypeBlockAnnounce, "block-announce", network.TopicIDBlock, true},
		{TypeBlocksRequest, "blocks-request", network.TopicIDUnspecified, false},
		{TypeBlocksResponse, "blocks-response", network.TopicIDUnspecified, false},
	}

	for _, tc := range testCases {
		msg := MakeMessage(tc.msgType)
		assert.Equal(t, tc.typeName, msg.Type().String())
		assert.Equal(t, tc.topicID, msg.TopicID())
		assert.Equal(t, tc.shouldBroadcast, msg.ShouldBroadcast())
	}
}
