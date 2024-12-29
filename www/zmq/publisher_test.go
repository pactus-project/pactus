package zmq

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasePublisher(t *testing.T) {
	suite := setup(t)

	topic := BlockInfo
	addr, err := url.Parse(fmt.Sprintf("tcp://127.0.0.1:%d", suite.FindFreePort()))
	require.NoError(t, err)

	base := &basePub{
		topic:     topic,
		zmqSocket: zmq4.NewPub(context.TODO()),
	}

	require.NoError(t, base.zmqSocket.Listen(addr.String()))
	assert.Equal(t, base.Address(), addr.Host)
	assert.Equal(t, base.TopicName(), topic.String())
}
