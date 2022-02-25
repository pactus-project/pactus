package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	msg := []byte("test-stream")
	require.NoError(t, tNetworks[0].SendTo(msg, tNetworks[1].SelfID()))

	e := shouldReceiveEvent(t, tNetworks[1]).(*StreamMessage)
	buf := make([]byte, len(msg))
	_, err := e.Reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, e.Source, tNetworks[0].SelfID())
	assert.Equal(t, buf, msg)
}
