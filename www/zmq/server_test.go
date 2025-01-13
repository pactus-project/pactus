package zmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
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

func (ts *testData) cleanup() {
	ts.server.Close()
	ts.resetServer()
}

func TestServerWithDefaultConfig(t *testing.T) {
	td := setup(t)

	conf := DefaultConfig()

	err := td.initServer(context.TODO(), conf)
	defer td.cleanup()

	assert.NoError(t, err)
	require.NotNil(t, td.server)
}

func TestTopicsWithSameSocket(t *testing.T) {
	td := setup(t)
	defer td.cleanup()

	port := td.FindFreePort()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr
	conf.ZmqPubTxInfo = addr
	conf.ZmqPubRawBlock = addr
	conf.ZmqPubRawTx = addr

	err := td.initServer(context.TODO(), conf)
	require.NoError(t, err)

	require.Len(t, td.server.publishers, 4)

	expectedAddr := td.server.publishers[0].Address()

	for _, pub := range td.server.publishers {
		require.Equal(t, expectedAddr, pub.Address(), "All publishers must have the same address")
	}
}

func TestTopicsWithDifferentSockets(t *testing.T) {
	td := setup(t)
	defer td.cleanup()

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = fmt.Sprintf("tcp://127.0.0.1:%d", td.FindFreePort())
	conf.ZmqPubTxInfo = fmt.Sprintf("tcp://127.0.0.1:%d", td.FindFreePort())
	conf.ZmqPubRawBlock = fmt.Sprintf("tcp://127.0.0.1:%d", td.FindFreePort())
	conf.ZmqPubRawTx = fmt.Sprintf("tcp://127.0.0.1:%d", td.FindFreePort())

	err := td.initServer(context.TODO(), conf)
	require.NoError(t, err)

	require.Len(t, td.server.publishers, 4)
}
