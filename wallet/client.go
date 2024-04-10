package wallet

import (
	"context"
	"errors"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AddressTypeBLSAccount string = "bls_account"
	AddressTypeValidator  string = "validator"
)

type grpcClient struct {
	ctx               context.Context
	cancel            func() // TODO: call me!
	servers           []string
	conn              *grpc.ClientConn
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func newGRPCClient(servers []string) *grpcClient {
	//  TODO: context should be passed here
	ctx, cancel := context.WithCancel(context.Background())

	return &grpcClient{
		ctx:               ctx,
		cancel:            cancel,
		servers:           servers,
		conn:              nil,
		blockchainClient:  nil,
		transactionClient: nil,
	}
}

func (c *grpcClient) connect() error {
	if c.conn != nil {
		return nil
	}

	maxTry := util.Min(3, len(c.servers))
	for i := 0; i < maxTry; i++ {
		n := util.RandInt32(int32(len(c.servers)))
		server := c.servers[n]
		conn, err := grpc.NewClient(server,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}

		blockchainClient := pactus.NewBlockchainClient(conn)
		transactionClient := pactus.NewTransactionClient(conn)

		// Check if client is responding
		_, err = blockchainClient.GetBlockchainInfo(c.ctx,
			&pactus.GetBlockchainInfoRequest{})
		if err != nil {
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
		&pactus.BroadcastTransactionRequest{SignedRawTransaction: data})
	if err != nil {
		return hash.UndefHash, err
	}

	return hash.FromBytes(res.Id)
}

// TODO: check the return value type.
func (c *grpcClient) getTransaction(id tx.ID) (*pactus.GetTransactionResponse, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	res, err := c.transactionClient.GetTransaction(c.ctx,
		&pactus.GetTransactionRequest{
			Id:        id.Bytes(),
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
