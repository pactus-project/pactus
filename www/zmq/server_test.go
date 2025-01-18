package zmq

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	mockState *state.MockState
	server    *Server
	eventCh   chan any
}

func setup(t *testing.T, conf *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockState := state.MockingState(ts)
	eventCh := make(chan any)
	server, err := New(context.TODO(), conf, eventCh)
	require.NoError(t, err)

	return &testData{
		TestSuite: ts,
		mockState: mockState,
		server:    server,
		eventCh:   eventCh,
	}
}

func (ts *testData) closeServer() {
	ts.server.Close()
}

func TestTopicsWithSameSocket(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr
	conf.ZmqPubTxInfo = addr
	conf.ZmqPubRawBlock = addr
	conf.ZmqPubRawTx = addr

	td := setup(t, conf)
	defer td.closeServer()

	require.Len(t, td.server.Publishers(), 4)
	require.Len(t, td.server.sockets, 1)

	expectedAddr := td.server.Publishers()[0].Address()

	for _, pub := range td.server.Publishers() {
		require.Equal(t, expectedAddr, pub.Address(), "All publishers must have the same address")
		require.Equal(t, conf.ZmqPubHWM, pub.HWM(), "All publishers must have the same HWM")
	}
}

func TestTopicsWithDifferentSockets(t *testing.T) {
	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubTxInfo = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubRawBlock = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubRawTx = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())

	td := setup(t, conf)
	defer td.closeServer()

	require.Len(t, td.server.Publishers(), 4)
	require.Len(t, td.server.sockets, 4)

	for _, pub := range td.server.Publishers() {
		require.Equal(t, conf.ZmqPubHWM, pub.HWM(), "All publishers must have the same HWM")
	}
}

func TestPublishMessageOnSameSockets(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://localhost:%d", port)
	conf := DefaultConfig()
	conf.ZmqPubRawTx = addr
	conf.ZmqPubTxInfo = addr
	conf.ZmqPubRawBlock = addr
	conf.ZmqPubBlockInfo = addr

	td := setup(t, conf)
	defer td.closeServer()

	td.server.Publishers()

	sub := zmq4.NewSub(context.TODO(), zmq4.WithAutomaticReconnect(false))

	err := sub.Dial(addr)
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicTransactionInfo.Bytes()))
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicRawTransaction.Bytes()))
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicBlockInfo.Bytes()))
	require.NoError(t, err)

	err = sub.SetOption(zmq4.OptionSubscribe, string(TopicRawBlock.Bytes()))
	require.NoError(t, err)

	blk, _ := td.TestSuite.GenerateTestBlock(td.RandHeight())

	td.eventCh <- blk

	for i := 0; i < (len(blk.Transactions())*2)+2; i++ {
		received, err := sub.Recv()
		require.NoError(t, err)

		require.NotNil(t, received.Frames)
		require.GreaterOrEqual(t, len(received.Frames), 1)

		msg := received.Frames[0]

		topic := TopicFromBytes(msg[:2])
		blockNumberOffset := len(msg) - 8
		height := binary.BigEndian.Uint32(msg[blockNumberOffset : blockNumberOffset+4])
		seqNo := binary.BigEndian.Uint32(msg[len(msg)-4:])
		t.Logf("[%s] %d", topic, seqNo)

		require.Equal(t, height, blk.Height())

		switch topic {
		case TopicRawTransaction:
			rawTx := msg[2 : len(msg)-8]

			txn, err := tx.FromBytes(rawTx)

			require.NoError(t, err)
			require.NotNil(t, txn)
			require.Equal(t, TopicRawTransaction, topic)
			require.NotEqual(t, 0, txn.SerializeSize())
		case TopicTransactionInfo:
			txHash := msg[2:34]
			id, err := hash.FromBytes(txHash)

			require.NoError(t, err)
			require.NotNil(t, id)
			require.Equal(t, TopicTransactionInfo, topic)

		case TopicRawBlock:
			blockHeader := msg[2 : len(msg)-8]
			buf := bytes.NewBuffer(blockHeader)
			header := new(block.Header)

			require.NoError(t, header.Decode(buf))
			require.NotNil(t, header)
			require.Equal(t, TopicRawBlock, topic)
			require.Equal(t, header.PrevBlockHash(), blk.Header().PrevBlockHash())
			require.Equal(t, header.StateRoot(), blk.Header().StateRoot())
		case TopicBlockInfo:
			proposerBytes := msg[2:23]
			timestamp := binary.BigEndian.Uint32(msg[23:27])
			txCount := binary.BigEndian.Uint16(msg[27:29])

			require.Equal(t, TopicBlockInfo, topic)
			require.Equal(t, blk.Header().ProposerAddress().Bytes(), proposerBytes)
			require.Equal(t, blk.Header().UnixTime(), timestamp)
			require.Equal(t, uint16(len(blk.Transactions())), txCount)
		}
	}

	require.NoError(t, sub.Close())
}
