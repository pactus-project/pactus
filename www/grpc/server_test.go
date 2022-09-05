package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
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

	zarb.RegisterBlockchainServer(s, blockchainServer)
	zarb.RegisterNetworkServer(s, networkServer)
	zarb.RegisterTransactionServer(s, transactionServer)
	go func() {
		if err := s.Serve(tListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return tListener.Dial()
}

func callBlockchainServer(t *testing.T) (*grpc.ClientConn, zarb.BlockchainClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial blockchain server: %v", err)
	}
	return conn, zarb.NewBlockchainClient(conn)
}

func callNetworkServer(t *testing.T) (*grpc.ClientConn, zarb.NetworkClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial network server: %v", err)
	}
	return conn, zarb.NewNetworkClient(conn)
}

func callTransactionServer(t *testing.T) (*grpc.ClientConn, zarb.TransactionClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial transaction server: %v", err)
	}
	return conn, zarb.NewTransactionClient(conn)
}
