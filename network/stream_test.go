package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	size := 6
	nets := setup(t, size)
	msg := []byte("test-stream")
	require.NoError(t, nets[0].SendTo(msg, nets[1].SelfID()))

	e := shouldReceiveEvent(t, nets[1]).(*StreamMessage)
	buf := make([]byte, len(msg))
	_, err := e.Reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, e.Source, nets[0].SelfID())
	assert.Equal(t, buf, msg)
}
