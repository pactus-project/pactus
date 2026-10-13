package metric

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSentMetric(t *testing.T) {
	metric := NewMetric()

	testMsgType := message.Type(1)

	metric.UpdateSentMetric(testMsgType, 100)

	assert.Equal(t, int64(1), metric.TotalSent.Bundles)
	assert.Equal(t, int64(100), metric.TotalSent.Bytes)

	assert.NotNil(t, metric.MessageSent[testMsgType])
	assert.Equal(t, int64(1), metric.MessageSent[testMsgType].Bundles)
	assert.Equal(t, int64(100), metric.MessageSent[testMsgType].Bytes)
}

func TestUpdateReceivedMetric(t *testing.T) {
	metric := NewMetric()

	testMsgType := message.Type(2)

	metric.UpdateReceivedMetric(testMsgType, 200)

	assert.Equal(t, int64(1), metric.TotalReceived.Bundles)
	assert.Equal(t, int64(200), metric.TotalReceived.Bytes)

	assert.NotNil(t, metric.MessageReceived[testMsgType])
	assert.Equal(t, int64(1), metric.MessageReceived[testMsgType].Bundles)
	assert.Equal(t, int64(200), metric.MessageReceived[testMsgType].Bytes)
}

func TestUpdateInvalidMetric(t *testing.T) {
	metric := NewMetric()

	metric.UpdateInvalidMetric(123)
	assert.Equal(t, int64(1), metric.TotalInvalid.Bundles)
	assert.Equal(t, int64(123), metric.TotalInvalid.Bytes)
}

func TestClone(t *testing.T) {
	msgType1 := message.Type(1)
	msgType2 := message.Type(2)

	original := NewMetric()
	original.UpdateSentMetric(msgType1, 100)
	original.UpdateReceivedMetric(msgType2, 50)
	original.UpdateInvalidMetric(7)

	clone := original.Clone()
	assert.Equal(t, original, clone)

	// The clone doesn't share anything with the original metric.
	clone.UpdateSentMetric(msgType1, 1)
	clone.UpdateSentMetric(msgType2, 1)
	clone.UpdateReceivedMetric(msgType2, 1)
	clone.UpdateInvalidMetric(1)

	assert.Equal(t, int64(1), original.TotalSent.Bundles)
	assert.Equal(t, int64(100), original.MessageSent[msgType1].Bytes)
	assert.NotContains(t, original.MessageSent, msgType2)
	assert.Equal(t, int64(50), original.MessageReceived[msgType2].Bytes)
	assert.Equal(t, int64(7), original.TotalInvalid.Bytes)

	// And the original metric doesn't change the clone.
	original.UpdateSentMetric(msgType1, 10)
	assert.Equal(t, int64(101), clone.MessageSent[msgType1].Bytes)
}
