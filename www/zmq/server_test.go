package zmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/util/flume"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	mockState *state.MockState
	server    *Server
	pipe      flume.Pipeline[any]
}

func setup(t *testing.T, conf *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockState := state.MockingState(ts)
	pipe := flume.MockingPipeline[any]()
	server, err := New(context.TODO(), conf, pipe)
	require.NoError(t, err)

	return &testData{
		TestSuite: ts,
		mockState: mockState,
		server:    server,
		pipe:      pipe,
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
