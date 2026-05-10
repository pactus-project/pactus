package grpc

import (
	"context"
	"encoding/base64"
	"net"
	"testing"

	consmgr "github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/zmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type testData struct {
	*testsuite.TestSuite

	mockState     *state.MockState
	mockSync      *sync.MockSync
	mockCons      *consmgr.MockReader
	mockConsMgr   *consmgr.MockManagerReader
	mockWalletMgr *wltmgr.MockIManager
	listener      *bufconn.Listener
	server        *Server
}

func testConfig() *Config {
	conf := DefaultConfig()

	return conf
}

func setup(t *testing.T, conf *Config) *testData {
	t.Helper()

	if conf == nil {
		conf = testConfig()
	}

	ts := testsuite.NewTestSuite(t)

	// for saving test wallets in temp directory
	t.Chdir(util.TempDirPath())

	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	mockState := state.MockingState(ts)
	mockNet := network.MockingNetwork(ts, ts.RandPeerID())
	mockSync := sync.MockingSync(ts)
	mockConsMgr := consmgr.NewMockManagerReader(ts.Ctrl)
	mockCons := consmgr.NewMockReader(ts.Ctrl)

	pub, _ := ts.RandBLSKeyPair()
	mockCons.EXPECT().ConsensusKey().Return(pub).AnyTimes()

	mockConsMgr.EXPECT().Instances().Return([]consmgr.Reader{mockCons}).AnyTimes()

	mockState.CommitTestBlocks(10)
	mockWalletMgr := wltmgr.NewMockIManager(ts.MockingController())

	zmqPublishers := []zmq.Publisher{
		zmq.MockingPublisher("zmq_address", "zmq_topic", 100),
	}

	server := NewServer(t.Context(), conf,
		mockState, mockSync, mockNet, mockConsMgr,
		mockWalletMgr, zmqPublishers,
	)
	err := server.startListening(listener)
	require.NoError(t, err)

	return &testData{
		TestSuite:     ts,
		mockState:     mockState,
		mockSync:      mockSync,
		mockCons:      mockCons,
		mockConsMgr:   mockConsMgr,
		mockWalletMgr: mockWalletMgr,
		server:        server,
		listener:      listener,
	}
}

func (td *testData) StopServer() {
	td.server.StopServer()
	_ = td.listener.Close()
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
		td.StopServer()
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

func (td *testData) healthClient(t *testing.T) healthpb.HealthClient {
	t.Helper()

	return healthpb.NewHealthClient(td.newClient(t))
}

func TestHealthCheck(t *testing.T) {
	td := setup(t, nil)
	client := td.healthClient(t)

	services := []struct {
		name string
		id   string
	}{
		{name: "Aggregate"},
		{name: "Blockchain", id: pactus.Blockchain_ServiceDesc.ServiceName},
		{name: "Transaction", id: pactus.Transaction_ServiceDesc.ServiceName},
		{name: "Network", id: pactus.Network_ServiceDesc.ServiceName},
		{name: "Utils", id: pactus.Utils_ServiceDesc.ServiceName},
	}

	for _, service := range services {
		t.Run(service.name, func(t *testing.T) {
			response, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
				Service: service.id,
			})
			require.NoError(t, err)
			assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.Status)
		})
	}

	_, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Wallet_ServiceDesc.ServiceName,
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestHealthWatch(t *testing.T) {
	td := setup(t, nil)
	client := td.healthClient(t)

	stream, err := client.Watch(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Blockchain_ServiceDesc.ServiceName,
	})
	require.NoError(t, err)

	response, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.Status)
}

func TestHealthCheckWithWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true
	td := setup(t, conf)
	client := td.healthClient(t)

	response, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Wallet_ServiceDesc.ServiceName,
	})
	require.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.Status)
}

func TestHealthCheckBasicAuth(t *testing.T) {
	conf := testConfig()
	conf.BasicAuth = "user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2"
	td := setup(t, conf)
	client := td.healthClient(t)

	_, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	ctx := contextWithBasicAuth(t.Context(), "user", "password")
	response, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.Status)
}

func TestHealthWatchBasicAuth(t *testing.T) {
	conf := testConfig()
	conf.BasicAuth = "user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2"
	td := setup(t, conf)
	client := td.healthClient(t)

	stream, err := client.Watch(t.Context(), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	_, err = stream.Recv()
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	ctx := contextWithBasicAuth(t.Context(), "user", "password")
	stream, err = client.Watch(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)

	response, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.Status)
}

func contextWithBasicAuth(ctx context.Context, user, password string) context.Context {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", auth))
}
