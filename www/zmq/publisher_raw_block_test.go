package zmq

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestRawBlockPublisher(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	conf := DefaultConfig()
	conf.ZmqPubRawBlock = addr

	td := setup(t, conf)
	defer td.closeServer()

	td.server.Publishers()

	sub := zmq4.NewSub(context.TODO(), zmq4.WithAutomaticReconnect(false))

	err := sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicRawBlock.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())
	td.pipe.Send(blk)

	received, err := sub.Recv()
	require.NoError(t, err)

	require.NotNil(t, received.Frames)
	require.GreaterOrEqual(t, len(received.Frames), 1)

	msg := received.Frames[0]

	topic := msg[:2]
	blockHeader := msg[2 : len(msg)-8]
	height := binary.BigEndian.Uint32(msg[140 : len(msg)-4])
	seqNo := binary.BigEndian.Uint32(msg[len(msg)-4:])

	buf := bytes.NewBuffer(blockHeader)
	header := new(block.Header)

	require.NoError(t, header.Decode(buf))

	require.NotNil(t, header)
	require.Equal(t, uint32(0), seqNo)
	require.Equal(t, blk.Height(), height)
	require.Equal(t, TopicRawBlock.Bytes(), topic)
	require.Equal(t, header.PrevBlockHash(), blk.Header().PrevBlockHash())
	require.Equal(t, header.StateRoot(), blk.Header().StateRoot())

	require.NoError(t, sub.Close())
}
