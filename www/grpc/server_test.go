package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var tMockState *state.MockState
var tConsMocks []*consensus.MockConsensus
var tMockSync *sync.MockSync
var tListener *bufconn.Listener
var tCtx context.Context

func init() {
	// for saving test wallets in temp directory
	err := os.Chdir(util.TempDirPath())
	if err != nil {
		panic(err)
	}

	const bufSize = 1024 * 1024

	consMgr, consMocks := consensus.MockingManager([]crypto.Signer{
		bls.GenerateTestSigner(), bls.GenerateTestSigner(),
	})

	tListener = bufconn.Listen(bufSize)
	tConsMocks = consMocks
	tMockState = state.MockingState()
	tMockSync = sync.MockingSync()
	tCtx = context.Background()

	tMockState.CommitTestBlocks(10)
	logger := logger.NewLogger("_grpc", nil)

	s := grpc.NewServer()
	blockchainServer := &blockchainServer{
		state:   tMockState,
		consMgr: consMgr,
		logger:  logger,
	}
	networkServer := &networkServer{
		sync:   tMockSync,
		logger: logger,
	}
	transactionServer := &transactionServer{
		state:  tMockState,
		logger: logger,
	}
	walletServer := &walletServer{
		logger: logger,
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
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}
	return conn, pactus.NewBlockchainClient(conn)
}

func testNetworkClient(t *testing.T) (*grpc.ClientConn, pactus.NetworkClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial network server: %v", err)
	}
	return conn, pactus.NewNetworkClient(conn)
}

func testTransactionClient(t *testing.T) (*grpc.ClientConn, pactus.TransactionClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial transaction server: %v", err)
	}
	return conn, pactus.NewTransactionClient(conn)
}

func testWalletClient(t *testing.T) (*grpc.ClientConn, pactus.WalletClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial wallet server: %v", err)
	}
	return conn, pactus.NewWalletClient(conn)
}
