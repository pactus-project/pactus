package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	mockState   *state.MockState
	mockSync    *sync.MockSync
	mockConsMgr consensus.Manager
	gRPCServer  *grpc.Server
	httpServer  *Server
}

func (td *testData) StopServers() {
	td.httpServer.StopServer()
	td.gRPCServer.StopServer()
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	// Resetting http handlers in golang for unit testing:
	// https://stackoverflow.com/questions/40786526/resetting-http-handlers-in-golang-for-unit-testing
	//
	http.DefaultServeMux = new(http.ServeMux)

	mockState := state.MockingState(ts)
	mockSync := sync.MockingSync(ts)
	mockNet := network.MockingNetwork(ts, ts.RandPeerID())
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

	gRPCServer := grpc.NewServer(grpcConf, mockState, mockSync, mockNet, mockConsMgr)
	assert.NoError(t, gRPCServer.StartServer())

	httpServer := NewServer(httpConf, false)
	assert.NoError(t, httpServer.StartServer(gRPCServer.Address()))

	return &testData{
		TestSuite:   ts,
		mockState:   mockState,
		mockSync:    mockSync,
		mockConsMgr: mockConsMgr,
		gRPCServer:  gRPCServer,
		httpServer:  httpServer,
	}
}

func TestRootHandler(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)
	td.httpServer.RootHandler(w, r)
	assert.Equal(t, w.Code, 200)
	fmt.Println(w.Body)

	td.StopServers()
}
