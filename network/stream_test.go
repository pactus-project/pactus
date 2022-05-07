package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	msg := []byte("test-stream")
	require.NoError(t, tNetworks[0].SendTo(msg, tNetworks[2].SelfID()))

	for {
		e := shouldReceiveEvent(t, tNetworks[2])
		if e.Type() == EventTypeStream {
			stream := e.(*StreamMessage)
			buf := make([]byte, len(msg))
			_, err := stream.Reader.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, stream.Source, tNetworks[0].SelfID())
			assert.Equal(t, buf, msg)
			break
		}
	}
}

func TestCloseConnection(t *testing.T) {
	tNetworks[2].CloseConnection(tNetworks[3].SelfID())
	msg := []byte("test-stream")
	require.Error(t, tNetworks[2].SendTo(msg, tNetworks[3].SelfID()), "connection should be closed")
}
