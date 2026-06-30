package jsonrpc_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/grpc/fake"
	"github.com/pactus-project/pactus/www/jsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	gRPCServer    *fake.FakeGRPCServer
	jsonrpcServer *jsonrpc.Server
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
	httpConf := &jsonrpc.Config{
		Enable: true,
		Listen: "[::]:0",
	}
	gRPCServer := fake.NewFakeGRPCServer(t, ts, grpcConf)
	require.NoError(t, gRPCServer.Server.StartServer())

	jsonrpcServer := jsonrpc.NewServer(t.Context(), httpConf)
	require.NoError(t, jsonrpcServer.StartServer(gRPCServer.Server.Address()))

	t.Cleanup(func() {
		jsonrpcServer.StopServer()
	})

	return &testData{
		TestSuite: ts,

		gRPCServer:    gRPCServer,
		jsonrpcServer: jsonrpcServer,
	}
}

func TestBlockchainInfo(t *testing.T) {
	td := setup(t)

	requestBody, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "pactus.blockchain.get_blockchain_info",
		"params":  map[string]any{},
	})
	require.NoError(t, err)

	resp, err := http.Post(
		"http://"+td.jsonrpcServer.Address(),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	result, ok := response["result"].(map[string]any)
	require.True(t, ok)

	fmt.Println(result)

	fakeState := td.gRPCServer.FakeState

	// Zero-value fields are omitted from the JSON response due to omitempty.
	_, hasLastBlockHeight := result["last_block_height"]
	assert.False(t, hasLastBlockHeight, "last_block_height should be omitted (value is 0)")

	assert.Equal(t, fakeState.LastBlockHash().String(), result["last_block_hash"])
	assert.Equal(t, float64(fakeState.FakeTime.Unix()), result["last_block_time"])

	_, hasPruningHeight := result["pruning_height"]
	assert.False(t, hasPruningHeight, "pruning_height should be omitted (value is 0)")

	_, hasIsPruned := result["is_pruned"]
	assert.False(t, hasIsPruned, "is_pruned should be omitted (value is false)")
}
