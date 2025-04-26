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

func TestTxInfoPublisher(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	conf := DefaultConfig()
	conf.ZmqPubTxInfo = addr

	td := setup(t, conf)
	defer td.closeServer()

	td.server.Publishers()

	sub := zmq4.NewSub(context.TODO(), zmq4.WithAutomaticReconnect(false))

	err := sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicTransactionInfo.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())
	td.pipe.Send(blk)

	for i := 0; i < len(blk.Transactions()); i++ {
		received, err := sub.Recv()
		require.NoError(t, err)

		require.NotNil(t, received.Frames)
		require.GreaterOrEqual(t, len(received.Frames), 1)

		msg := received.Frames[0]
		require.Len(t, msg, 42)

		topic := msg[:2]
		txHash := msg[2:34]
		height := binary.BigEndian.Uint32(msg[34:38])
		seqNo := binary.BigEndian.Uint32(msg[38:])

		require.Equal(t, TopicTransactionInfo.Bytes(), topic)
		require.Equal(t, blk.Transactions()[i].ID().Bytes(), txHash)
		require.Equal(t, blk.Height(), height)
		require.Equal(t, uint32(i), seqNo)
	}

	require.NoError(t, sub.Close())
}
