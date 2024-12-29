package zmq

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestBasePublisher(t *testing.T) {
	ts := setup(t)

	topic := BlockInfo
	addr, err := url.Parse(fmt.Sprintf("tcp://127.0.0.1:%d", ts.FindFreePort()))
	require.NoError(t, err)

	base := &basePub{
		topic:     topic,
		zmqSocket: zmq4.NewPub(context.TODO()),
	}

	require.NoError(t, base.zmqSocket.Listen(addr.String()))
	assert.Equal(t, base.Address(), addr.Host)
	assert.Equal(t, base.TopicName(), topic.String())
}
