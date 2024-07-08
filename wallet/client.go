package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gRPCClient is a gRPC client that randomly establishes a connection to a gRPC server.
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
			grpc.WithContextDialer(func(_ context.Context, s string) (net.Conn, error) {
				return net.DialTimeout("tcp", s, c.timeout)
			}))
		if err != nil {
			continue
		}

		blockchainClient := pactus.NewBlockchainClient(conn)
		transactionClient := pactus.NewTransactionClient(conn)

		// Check if client is responding
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

	info, err := c.blockchainClient.GetBlockchainInfo(c.ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *grpcClient) getAccount(addrStr string) (*pactus.AccountInfo, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	res, err := c.blockchainClient.GetAccount(c.ctx,
		&pactus.GetAccountRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	return res.Account, nil
}

func (c *grpcClient) getValidator(addrStr string) (*pactus.ValidatorInfo, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	res, err := c.blockchainClient.GetValidator(c.ctx,
		&pactus.GetValidatorRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	return res.Validator, nil
}

func (c *grpcClient) sendTx(trx *tx.Tx) (tx.ID, error) {
	if err := c.connect(); err != nil {
		return hash.UndefHash, err
	}

	data, err := trx.Bytes()
	if err != nil {
		return hash.UndefHash, err
	}
	res, err := c.transactionClient.BroadcastTransaction(c.ctx,
		&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
	if err != nil {
		return hash.UndefHash, err
	}

	return hash.FromString(res.Id)
}

// TODO: check the return value type.
func (c *grpcClient) getTransaction(id tx.ID) (*pactus.GetTransactionResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	res, err := c.transactionClient.GetTransaction(c.ctx,
		&pactus.GetTransactionRequest{
			Id:        id.String(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *grpcClient) getFee(amt amount.Amount, payloadType payload.Type) (amount.Amount, error) {
	if err := c.connect(); err != nil {
		return 0, err
	}

	res, err := c.transactionClient.CalculateFee(c.ctx,
		&pactus.CalculateFeeRequest{
			Amount:      amt.ToNanoPAC(),
			PayloadType: pactus.PayloadType(payloadType),
		})
	if err != nil {
		return 0, err
	}

	return amount.Amount(res.Fee), nil
}
