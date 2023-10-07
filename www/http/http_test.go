package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/stretchr/testify/assert"
)

var tTestData *testData

type testData struct {
	*testsuite.TestSuite

	mockState   *state.MockState
	mockSync    *sync.MockSync
	mockConsMgr consensus.Manager
	gRPCServer  *grpc.Server
	httpServer  *Server
}

func setup(t *testing.T) *testData {
	t.Helper()

	if tTestData != nil {
		return tTestData
	}

	ts := testsuite.NewTestSuite(t)

	mockState := state.MockingState(ts)
	mockSync := sync.MockingSync(ts)
	mockConsMgr, _ := consensus.MockingManager(ts, []*bls.ValidatorKey{
		ts.RandValKey(), ts.RandValKey(),
	})

	mockConsMgr.MoveToNewHeight()

	grpcConf := &grpc.Config{
		Enable: true,
		Listen: "[::]:0",
	}
	httpConf := &Config{
		Enable: true,
		Listen: "[::]:0",
	}

	gRPCServer := grpc.NewServer(grpcConf, mockState, mockSync, mockConsMgr)
	assert.NoError(t, gRPCServer.StartServer())

	httpServer := NewServer(httpConf)
	assert.NoError(t, httpServer.StartServer(gRPCServer.Address()))

	tTestData = &testData{
		TestSuite:   ts,
		mockState:   mockState,
		mockSync:    mockSync,
		mockConsMgr: mockConsMgr,
		gRPCServer:  gRPCServer,
		httpServer:  httpServer,
	}
	return tTestData
}

func TestRootHandler(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)
	td.httpServer.RootHandler(w, r)
	assert.Equal(t, w.Code, 200)
	fmt.Println(w.Body)
}
