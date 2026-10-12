package message

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		msg, err := MakeMessage(tt.msgType)
		require.NoError(t, err)

		assert.Equal(t, tt.typeName, msg.Type().String())
		assert.Equal(t, tt.topicID, msg.TopicID())
		assert.Equal(t, tt.shouldBroadcast, msg.ShouldBroadcast())
	}
}

func TestInvalidMessageType(t *testing.T) {
	_, err := MakeMessage(66)
	require.ErrorIs(t, err, InvalidMessageTypeError{Type: 66})
}

func TestMessageTypeString(t *testing.T) {
	tests := map[Type]string{
		TypeHello:          "hello",
		TypeHelloAck:       "hello-ack",
		TypeTransaction:    "transaction",
		TypeQueryProposal:  "query-proposal",
		TypeProposal:       "proposal",
		TypeQueryVote:      "query-vote",
		TypeVote:           "vote",
		TypeBlockAnnounce:  "block-announce",
		TypeBlocksRequest:  "blocks-request",
		TypeBlocksResponse: "blocks-response",
		Type(1234):         "1234",
	}

	for msgType, name := range tests {
		assert.Equal(t, name, msgType.String())
	}
}

func TestResponseCodeString(t *testing.T) {
	tests := map[ResponseCode]string{
		ResponseCodeOK:           "ok",
		ResponseCodeRejected:     "rejected",
		ResponseCodeMoreBlocks:   "more-blocks",
		ResponseCodeNoMoreBlocks: "no-more-blocks",
		ResponseCodeSynced:       "synced",
		ResponseCode(1234):       "1234",
	}

	for code, name := range tests {
		assert.Equal(t, name, code.String())
	}
}

func TestErrorMessages(t *testing.T) {
	assert.Equal(t, "bad message", BasicCheckError{Reason: "bad message"}.Error())
	assert.Equal(t, "invalid message type: 99", InvalidMessageTypeError{Type: 99}.Error())
}
