package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	msg := []byte("test-stream")
	require.NoError(t, tNetwork1.SendTo(msg, tNetwork2.SelfID()))

	e := shouldReceiveEvent(t, tNetwork2).(*StreamMessage)
	buf := make([]byte, len(msg))
	_, err := e.Reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, e.Source, tNetwork1.SelfID())
	assert.Equal(t, buf, msg)
}
