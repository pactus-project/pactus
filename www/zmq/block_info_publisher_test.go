package zmq

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/go-zeromq/zmq4"
	"github.com/stretchr/testify/require"
)

func TestBlockInfoPublisher(t *testing.T) {
	td := setup(t)
	defer td.cleanup()

	port := td.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr

	err := td.initServer(ctx, conf)
	require.NoError(t, err)

	sub := zmq4.NewSub(ctx)
	defer func() {
		_ = sub.Close()
	}()

	err = sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(BlockInfo.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())

	td.eventCh <- blk

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

	require.Equal(t, BlockInfo.Bytes(), topic)
	require.Equal(t, blk.Header().ProposerAddress().Bytes(), proposerBytes)
	require.Equal(t, blk.Header().UnixTime(), timestamp)
	require.Equal(t, uint16(len(blk.Transactions())), txCount)
	require.Equal(t, blk.Height(), height)
	require.Equal(t, uint32(0), seqNo)
}
