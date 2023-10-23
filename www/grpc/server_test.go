package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	tMockState *state.MockState
	tConsMocks []*consensus.MockConsensus
	tMockSync  *sync.MockSync
	tListener  *bufconn.Listener
	tCtx       context.Context
)

func init() {
	ts := testsuite.NewTestSuiteForSeed(0x1234)

	// for saving test wallets in temp directory
	err := os.Chdir(util.TempDirPath())
	if err != nil {
		panic(err)
	}

	const bufSize = 1024 * 1024

	consMgr, consMocks := consensus.MockingManager(ts, []*bls.ValidatorKey{
		ts.RandValKey(), ts.RandValKey(),
	})

	tListener = bufconn.Listen(bufSize)
	tConsMocks = consMocks
	tMockState = state.MockingState(ts)
	tMockSync = sync.MockingSync(ts)
	tCtx = context.Background()

	tMockState.CommitTestBlocks(10)
	subLogger := logger.NewSubLogger("_grpc", nil)

	s := grpc.NewServer()
	blockchainServer := &blockchainServer{
		state:   tMockState,
		consMgr: consMgr,
		logger:  subLogger,
	}
	networkServer := &networkServer{
		sync:   tMockSync,
		logger: subLogger,
	}
	transactionServer := &transactionServer{
		state:  tMockState,
		logger: subLogger,
	}
	walletServer := &walletServer{
		logger: subLogger,
	}

	pactus.RegisterBlockchainServer(s, blockchainServer)
	pactus.RegisterNetworkServer(s, networkServer)
	pactus.RegisterTransactionServer(s, transactionServer)
	pactus.RegisterWalletServer(s, walletServer)

	go func() {
		if err := s.Serve(tListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return tListener.Dial()
}

func testBlockchainClient(t *testing.T) (*grpc.ClientConn, pactus.BlockchainClient) {
	t.Helper()

	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}
	return conn, pactus.NewBlockchainClient(conn)
}

func testNetworkClient(t *testing.T) (*grpc.ClientConn, pactus.NetworkClient) {
	t.Helper()

	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial network server: %v", err)
	}
	return conn, pactus.NewNetworkClient(conn)
}

func testTransactionClient(t *testing.T) (*grpc.ClientConn, pactus.TransactionClient) {
	t.Helper()

	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial transaction server: %v", err)
	}
	return conn, pactus.NewTransactionClient(conn)
}

func testWalletClient(t *testing.T) (*grpc.ClientConn, pactus.WalletClient) {
	t.Helper()

	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial wallet server: %v", err)
	}
	return conn, pactus.NewWalletClient(conn)
}
