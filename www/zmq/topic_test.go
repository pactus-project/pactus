package zmq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTopicFromBytes(t *testing.T) {
	validRawTopic := TopicRawTransaction.Bytes()
	invalidRawTopic := make([]byte, 0)

	topic := TopicFromBytes(validRawTopic)
	require.Equal(t, TopicRawTransaction, topic)

	topic = TopicFromBytes(invalidRawTopic)
	require.Equal(t, 0, int(topic))
}
