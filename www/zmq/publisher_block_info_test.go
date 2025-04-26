package zmq

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestBlockInfoPublisher(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr

	td := setup(t, conf)
	defer td.closeServer()

	td.server.Publishers()

	sub := zmq4.NewSub(context.TODO(), zmq4.WithAutomaticReconnect(false))

	err := sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicBlockInfo.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())
	td.pipe.Send(blk)

	received, err := sub.Recv()
	require.NoError(t, err)

	require.NotNil(t, received.Frames)
	require.GreaterOrEqual(t, len(received.Frames), 1)

	msg := received.Frames[0]
	require.Len(t, msg, 37)

	topic := msg[:2]
	proposerBytes := msg[2:23]
	timestamp := binary.BigEndian.Uint32(msg[23:27])
	txCount := binary.BigEndian.Uint16(msg[27:29])
	height := binary.BigEndian.Uint32(msg[29:33])
	seqNo := binary.BigEndian.Uint32(msg[33:])

	require.Equal(t, TopicBlockInfo.Bytes(), topic)
	require.Equal(t, blk.Header().ProposerAddress().Bytes(), proposerBytes)
	require.Equal(t, blk.Header().UnixTime(), timestamp)
	require.Equal(t, uint16(len(blk.Transactions())), txCount)
	require.Equal(t, blk.Height(), height)
	require.Equal(t, uint32(0), seqNo)

	require.NoError(t, sub.Close())
}
