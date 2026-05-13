package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	pactusgrpc "github.com/pactus-project/pactus/www/grpc"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/grpc/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testData struct {
	*testsuite.TestSuite

	listener *bufconn.Listener
	server   *mock.MockGRPCServer
}

func testConfig() *pactusgrpc.Config {
	conf := pactusgrpc.DefaultConfig()
	conf.Listen = ""

	return conf
}

func setup(t *testing.T, conf *pactusgrpc.Config) *testData {
	t.Helper()

	if conf == nil {
		conf = testConfig()
	}

	ts := testsuite.NewTestSuite(t)
	gRPCServer := mock.SetupServer(t, ts, conf)

	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)

	err := gRPCServer.Server.StartListening(listener)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, listener.Close())
	})

	return &testData{
		TestSuite: ts,
		listener:  listener,
		server:    gRPCServer,
	}
}

func (td *testData) bufDialer(context.Context, string) (net.Conn, error) {
	return td.listener.Dial()
}

func (td *testData) newClient(t *testing.T) *grpc.ClientConn {
	t.Helper()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(td.bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	return conn
}

func (td *testData) blockchainClient(t *testing.T) pactus.BlockchainClient {
	t.Helper()

	return pactus.NewBlockchainClient(td.newClient(t))
}

func (td *testData) networkClient(t *testing.T) pactus.NetworkClient {
	t.Helper()

	return pactus.NewNetworkClient(td.newClient(t))
}

func (td *testData) transactionClient(t *testing.T) pactus.TransactionClient {
	t.Helper()

	return pactus.NewTransactionClient(td.newClient(t))
}

func (td *testData) walletClient(t *testing.T) pactus.WalletClient {
	t.Helper()

	return pactus.NewWalletClient(td.newClient(t))
}

func (td *testData) utilClient(t *testing.T) pactus.UtilsClient {
	t.Helper()

	return pactus.NewUtilsClient(td.newClient(t))
}
