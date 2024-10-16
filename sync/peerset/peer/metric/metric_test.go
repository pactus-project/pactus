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
