package zmq

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestRawTxPublisher(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	conf := DefaultConfig()
	conf.ZmqPubRawTx = addr

	td := setup(t, conf)
	defer td.closeServer()

	td.server.Publishers()

	sub := zmq4.NewSub(context.TODO(), zmq4.WithAutomaticReconnect(false))

	err := sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicRawTransaction.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())

	td.eventCh <- blk

	for i := 0; i < len(blk.Transactions()); i++ {
		received, err := sub.Recv()
		require.NoError(t, err)

		require.NotNil(t, received.Frames)
		require.GreaterOrEqual(t, len(received.Frames), 1)

		msg := received.Frames[0]

		topic := msg[:2]
		rawTx := msg[2 : len(msg)-8]

		blockNumberOffset := len(msg) - 8
		height := binary.BigEndian.Uint32(msg[blockNumberOffset : blockNumberOffset+4])
		seqNo := binary.BigEndian.Uint32(msg[len(msg)-4:])

		txn, err := tx.FromBytes(rawTx)
		require.NoError(t, err)
		require.NotNil(t, txn)

		require.Equal(t, TopicRawTransaction.Bytes(), topic)
		require.Equal(t, height, blk.Height())
		require.Equal(t, uint32(i), seqNo)
		require.NotEqual(t, 0, txn.SerializeSize())
	}

	require.NoError(t, sub.Close())
}
