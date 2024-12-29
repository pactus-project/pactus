package zmq

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	mockState *state.MockState
	server    *Server
	eventCh   chan any
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockState := state.MockingState(ts)

	return &testData{
		TestSuite: ts,
		mockState: mockState,
	}
}

func (ts *testData) initServer(ctx context.Context, conf *Config) error {
	eventCh := make(chan any)
	sv, err := New(ctx, conf, eventCh)
	if err != nil {
		return err
	}

	ts.server = sv
	ts.eventCh = eventCh

	return nil
}

func (ts *testData) resetServer() {
	ts.server = nil
	ts.eventCh = nil
}

func (ts *testData) cleanup() func() {
	return func() {
		ts.server.Close()
		ts.resetServer()
	}
}

func TestServerWithDefaultConfig(t *testing.T) {
	ts := setup(t)

	conf := DefaultConfig()

	err := ts.initServer(context.TODO(), conf)
	t.Cleanup(ts.cleanup())

	assert.NoError(t, err)
	require.NotNil(t, ts.server)
}

func TestTopicsWithSameSocket(t *testing.T) {
	ts := setup(t)
	t.Cleanup(ts.cleanup())

	port := ts.FindFreePort()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr
	conf.ZmqPubTxInfo = addr
	conf.ZmqPubRawBlock = addr
	conf.ZmqPubRawTx = addr

	err := ts.initServer(context.TODO(), conf)
	require.NoError(t, err)

	require.Len(t, ts.server.publishers, 4)

	expectedAddr := ts.server.publishers[0].Address()

	for _, pub := range ts.server.publishers {
		require.Equal(t, expectedAddr, pub.Address(), "All publishers must have the same address")
	}
}

func TestTopicsWithDifferentSockets(t *testing.T) {
	ts := setup(t)
	t.Cleanup(ts.cleanup())

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = fmt.Sprintf("tcp://127.0.0.1:%d", ts.FindFreePort())
	conf.ZmqPubTxInfo = fmt.Sprintf("tcp://127.0.0.1:%d", ts.FindFreePort())
	conf.ZmqPubRawBlock = fmt.Sprintf("tcp://127.0.0.1:%d", ts.FindFreePort())
	conf.ZmqPubRawTx = fmt.Sprintf("tcp://127.0.0.1:%d", ts.FindFreePort())

	err := ts.initServer(context.TODO(), conf)
	require.NoError(t, err)

	require.Len(t, ts.server.publishers, 4)
}
