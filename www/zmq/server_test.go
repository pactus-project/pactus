package zmq

import (
	"context"
	"fmt"
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

func setup(ctx context.Context, t *testing.T, conf *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockState := state.MockingState(ts)
	eventCh := make(chan any)
	sv, err := New(ctx, conf, eventCh)
	require.NoError(t, err)

	return &testData{
		TestSuite: ts,
		server:    sv,
		mockState: mockState,
		eventCh:   eventCh,
	}
}

func (t *testData) cleanup() func() {
	return func() {
		t.server.Close()
	}
}

func TestTopicsWithSameSocket(t *testing.T) {
	port := testsuite.FindFreePort()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = addr
	conf.ZmqPubTxInfo = addr
	conf.ZmqPubRawBlock = addr
	conf.ZmqPubRawTx = addr

	ts := setup(context.TODO(), t, conf)
	t.Cleanup(ts.cleanup())

	require.Len(t, ts.server.publishers, 4)
}

func TestTopicsWithDifferentSockets(t *testing.T) {
	conf := DefaultConfig()
	conf.ZmqPubBlockInfo = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubTxInfo = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubRawBlock = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())
	conf.ZmqPubRawTx = fmt.Sprintf("tcp://127.0.0.1:%d", testsuite.FindFreePort())

	ts := setup(context.TODO(), t, conf)
	t.Cleanup(ts.cleanup())

	require.Len(t, ts.server.publishers, 4)
}
