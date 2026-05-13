package html_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/grpc/mock"
	"github.com/pactus-project/pactus/www/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	gRPCServer *mock.MockGRPCServer
	httpServer *html.Server
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	// Resetting http handlers in golang for unit testing:
	// https://stackoverflow.com/questions/40786526/resetting-http-handlers-in-golang-for-unit-testing
	//
	http.DefaultServeMux = new(http.ServeMux)

	grpcConf := &grpc.Config{
		Enable: true,
		Listen: "[::]:0",
	}
	httpConf := &html.Config{
		Enable: true,
		Listen: "[::]:0",
	}
	gRPCServer := mock.SetupServer(t, ts, grpcConf)
	require.NoError(t, gRPCServer.Server.StartServer())

	httpServer := html.NewServer(t.Context(), httpConf, false)
	require.NoError(t, httpServer.StartServer(gRPCServer.Server.Address()))

	t.Cleanup(func() {
		httpServer.StopServer()
	})

	return &testData{
		TestSuite: ts,

		gRPCServer: gRPCServer,
		httpServer: httpServer,
	}
}

func TestRootHandler(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)
	td.httpServer.RootHandler(w, r)
	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body)
}
