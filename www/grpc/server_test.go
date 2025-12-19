package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	walletMgr "github.com/pactus-project/pactus/wallet/manager"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/zmq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testData struct {
	*testsuite.TestSuite

	mockState     *state.MockState
	mockSync      *sync.MockSync
	consMocks     []*consensus.MockConsensus
	mockConsMgr   manager.Manager
	mockWalletMgr *walletMgr.MockIManager
	listener      *bufconn.Listener
	server        *Server
}

func testConfig() *Config {
	conf := DefaultConfig()
	conf.WalletsDir = util.TempDirPath()

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
	valKeys := []*bls.ValidatorKey{ts.RandValKey(), ts.RandValKey()}
	mockState := state.MockingState(ts)
	mockNet := network.MockingNetwork(ts, ts.RandPeerID())
	mockSync := sync.MockingSync(ts)
	mockConsMgr, consMocks := manager.MockingManager(ts, mockState, valKeys)

	mockState.CommitTestBlocks(10)
	mockWalletMgr := walletMgr.NewMockIManager(ts.MockingController())

	zmqPublishers := []zmq.Publisher{
		zmq.MockingPublisher("zmq_address", "zmq_topic", 100),
	}

	server := NewServer(context.Background(), conf,
		mockState, mockSync, mockNet, mockConsMgr,
		mockWalletMgr, zmqPublishers,
	)
	err := server.startListening(listener)
	assert.NoError(t, err)

	return &testData{
		TestSuite:     ts,
		mockState:     mockState,
		mockSync:      mockSync,
		consMocks:     consMocks,
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
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.Nil(t, conn.Close())
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
