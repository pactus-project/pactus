package remoteprovider

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcClient is a gRPC client that randomly establishes a connection to a gRPC server.
// It is used to get information such as account balance or transaction data from the server.
type grpcClient struct {
	ctx               context.Context
	servers           []string
	conn              *grpc.ClientConn
	timeout           time.Duration
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func newGrpcClient(timeout time.Duration, servers []string) *grpcClient {
	ctx := context.Background()

	cli := &grpcClient{
		ctx:               ctx,
		timeout:           timeout,
		conn:              nil,
		blockchainClient:  nil,
		transactionClient: nil,
	}

	if len(servers) > 0 {
		cli.servers = servers
	}

	return cli
}

func (c *grpcClient) connect() error {
	if c.conn != nil {
		return nil
	}

	for _, server := range c.servers {
		conn, err := grpc.NewClient(server,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
				return util.NetworkDialTimeout(ctx, "tcp", address, c.timeout)
			}))
		if err != nil {
			continue
		}

		blockchainClient := pactus.NewBlockchainClient(conn)
		transactionClient := pactus.NewTransactionClient(conn)

		// Check if client is responding
		// TODO: Use Ping API in version 1.11.0
		_, err = blockchainClient.GetBlockchainInfo(c.ctx,
			&pactus.GetBlockchainInfoRequest{})
		if err != nil {
			_ = conn.Close()

			continue
		}

		c.conn = conn
		c.blockchainClient = blockchainClient
		c.transactionClient = transactionClient

		return nil
	}

	return errors.New("unable to connect to the servers")
}

func (c *grpcClient) getBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	return c.blockchainClient.GetBlockchainInfo(ctx,
		&pactus.GetBlockchainInfoRequest{})
}

func (c *grpcClient) getAccount(addrStr string) (*pactus.GetAccountResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	return c.blockchainClient.GetAccount(ctx,
		&pactus.GetAccountRequest{Address: addrStr})
}

func (c *grpcClient) getValidator(addrStr string) (*pactus.GetValidatorResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	return c.blockchainClient.GetValidator(ctx,
		&pactus.GetValidatorRequest{Address: addrStr})
}

func (c *grpcClient) sendTx(trx *tx.Tx) (*pactus.BroadcastTransactionResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	data, err := trx.Bytes()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	return c.transactionClient.BroadcastTransaction(ctx,
		&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
}

func (c *grpcClient) getTransaction(txID tx.ID) (*pactus.GetTransactionResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	return c.transactionClient.GetTransaction(ctx,
		&pactus.GetTransactionRequest{
			Id:        txID.String(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_DATA,
		})
}
