package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var tMockState *state.MockState
var tMockSync *sync.MockSync
var tListener *bufconn.Listener
var tCtx context.Context

func init() {
	const bufSize = 1024 * 1024

	tListener = bufconn.Listen(bufSize)
	tMockState = state.MockingState()
	tMockSync = sync.MockingSync()
	tCtx = context.Background()

	tMockState.CommitTestBlocks(10)
	logger := logger.NewLogger("_grpc", nil)

	s := grpc.NewServer()
	blockchainServer := &blockchainServer{
		state:  tMockState,
		logger: logger,
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

func callBlockchainServer(t *testing.T) (*grpc.ClientConn, pactus.BlockchainClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}
	return conn, pactus.NewBlockchainClient(conn)
}

func callNetworkServer(t *testing.T) (*grpc.ClientConn, pactus.NetworkClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial network server: %v", err)
	}
	return conn, pactus.NewNetworkClient(conn)
}

func callTransactionServer(t *testing.T) (*grpc.ClientConn, pactus.TransactionClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial transaction server: %v", err)
	}
	return conn, pactus.NewTransactionClient(conn)
}

func callWalletSerer(t *testing.T) (*grpc.ClientConn, zarb.WalletClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial wallet server: %v", err)
	}
	return conn, zarb.NewWalletClient(conn)
}
