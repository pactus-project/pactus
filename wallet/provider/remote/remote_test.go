package remote

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// untrustedServer answers like a server that can't be trusted: its answers are empty
// or don't match the request.
type untrustedServer struct {
	pactus.UnimplementedBlockchainServer
	pactus.UnimplementedTransactionServer
}

func (untrustedServer) GetBlockchainInfo(context.Context, *pactus.GetBlockchainInfoRequest,
) (*pactus.GetBlockchainInfoResponse, error) {
	return &pactus.GetBlockchainInfoResponse{LastBlockHeight: 1}, nil
}

func (untrustedServer) GetAccount(context.Context, *pactus.GetAccountRequest,
) (*pactus.GetAccountResponse, error) {
	return &pactus.GetAccountResponse{}, nil
}

func (untrustedServer) GetValidator(context.Context, *pactus.GetValidatorRequest,
) (*pactus.GetValidatorResponse, error) {
	return &pactus.GetValidatorResponse{}, nil
}

func (untrustedServer) BroadcastTransaction(context.Context, *pactus.BroadcastTransactionRequest,
) (*pactus.BroadcastTransactionResponse, error) {
	return &pactus.BroadcastTransactionResponse{Id: "not the id of the transaction"}, nil
}

func startUntrustedServer(t *testing.T) string {
	t.Helper()

	listener, err := new(net.ListenConfig).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	server := grpc.NewServer()
	pactus.RegisterBlockchainServer(server, untrustedServer{})
	pactus.RegisterTransactionServer(server, untrustedServer{})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	return listener.Addr().String()
}

func TestUntrustedServerAnswers(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	provider, err := NewRemoteBlockchainProvider(t.Context(), genesis.Localnet,
		WithCustomServers([]string{startUntrustedServer(t)}), WithTimeout(time.Second))
	require.NoError(t, err)
	defer func() { _ = provider.Close() }()

	t.Run("an empty answer is an error, not a panic", func(t *testing.T) {
		_, err := provider.GetAccount(ts.RandAccAddress().String())
		require.Error(t, err)

		_, err = provider.GetValidator(ts.RandValAddress().String())
		require.Error(t, err)
	})

	t.Run("the ID of a sent transaction is computed locally", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		id, err := provider.SendTx(trx)
		require.NoError(t, err)
		assert.Equal(t, trx.ID().String(), id)
	})
}

// A server that accepts the connection but never answers doesn't block the wallet.
func TestSilentServer(t *testing.T) {
	listener, err := new(net.ListenConfig).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { _ = conn.Close() })
		}
	}()

	provider, err := NewRemoteBlockchainProvider(t.Context(), genesis.Localnet,
		WithCustomServers([]string{listener.Addr().String()}), WithTimeout(500*time.Millisecond))
	require.NoError(t, err)
	defer func() { _ = provider.Close() }()

	start := time.Now()
	_, err = provider.LastBlockHeight()
	require.Error(t, err)
	assert.Less(t, time.Since(start), 10*time.Second)
}
