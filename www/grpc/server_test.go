package grpc

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testData struct {
	*testsuite.TestSuite

	mockState     *state.MockState
	mockSync      *sync.MockSync
	consMocks     []*consensus.MockConsensus
	mockConsMgr   consensus.Manager
	defaultWallet *wallet.Wallet
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
	err := os.Chdir(util.TempDirPath())
	if err != nil {
		panic(err)
	}

	const bufSize = 1024 * 1024

	mockConsMgr, consMocks := consensus.MockingManager(ts, []*bls.ValidatorKey{
		ts.RandValKey(), ts.RandValKey(),
	})

	listener := bufconn.Listen(bufSize)
	mockState := state.MockingState(ts)
	mockNet := network.MockingNetwork(ts, ts.RandPeerID())
	mockSync := sync.MockingSync(ts)

	mockState.CommitTestBlocks(10)

	wltPath := filepath.Join(conf.WalletsDir, "default_wallet")
	mnemonic, _ := wallet.GenerateMnemonic(128)
	defaultWallet, err := wallet.Create(wltPath, mnemonic, "", genesis.Mainnet)
	require.NoError(t, err)
	require.NoError(t, defaultWallet.Save())

	server := NewServer(conf, mockState, mockSync, mockNet, mockConsMgr)
	err = server.startListening(listener)
	assert.NoError(t, err)

	return &testData{
		TestSuite:     ts,
		mockState:     mockState,
		mockSync:      mockSync,
		consMocks:     consMocks,
		mockConsMgr:   mockConsMgr,
		defaultWallet: defaultWallet,
		server:        server,
		listener:      listener,
	}
}

func (td *testData) StopServer() {
	td.server.StopServer()
	td.listener.Close()
}

func (td *testData) bufDialer(context.Context, string) (net.Conn, error) {
	return td.listener.Dial()
}

func (td *testData) blockchainClient(t *testing.T) (*grpc.ClientConn, pactus.BlockchainClient) {
	t.Helper()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(td.bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}

	return conn, pactus.NewBlockchainClient(conn)
}

func (td *testData) networkClient(t *testing.T) (*grpc.ClientConn, pactus.NetworkClient) {
	t.Helper()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(td.bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial network server: %v", err)
	}

	return conn, pactus.NewNetworkClient(conn)
}

func (td *testData) transactionClient(t *testing.T) (*grpc.ClientConn, pactus.TransactionClient) {
	t.Helper()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(td.bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial transaction server: %v", err)
	}

	return conn, pactus.NewTransactionClient(conn)
}

func (td *testData) walletClient(t *testing.T) (*grpc.ClientConn, pactus.WalletClient) {
	t.Helper()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(td.bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial wallet server: %v", err)
	}

	return conn, pactus.NewWalletClient(conn)
}
