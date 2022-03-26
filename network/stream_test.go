package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
)

func TestStream(t *testing.T) {
	size := 6
	nets := setup(t, size)
	i := util.RandInt32(int32(size-1)) + 1
	msg := []byte("test-stream")
	require.NoError(t, nets[0].SendTo(msg, nets[i].SelfID()))

	e := shouldReceiveEvent(t, nets[i]).(*StreamMessage)
	buf := make([]byte, len(msg))
	_, err := e.Reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, e.Source, nets[0].SelfID())
	assert.Equal(t, buf, msg)
}
