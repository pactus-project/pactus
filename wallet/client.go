package wallet

import (
	"context"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AddressTypeBLSAccount string = "bls_account"
	AddressTypeValidator  string = "validator"
)

type grpcClient struct {
	ctx    context.Context
	cancel func() // TODO: call me!

	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func newGRPCClient(rpcEndpoint string) (*grpcClient, error) {
	conn, err := grpc.Dial(rpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &grpcClient{
		ctx:               ctx,
		cancel:            cancel,
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
	}, nil
}

func (c *grpcClient) getBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	info, err := c.blockchainClient.GetBlockchainInfo(c.ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *grpcClient) getAccount(addr crypto.Address) (*pactus.AccountInfo, error) {
	res, err := c.blockchainClient.GetAccount(c.ctx,
		&pactus.GetAccountRequest{Address: addr.String()})
	if err != nil {
		return nil, err
	}

	return res.Account, nil
}

func (c *grpcClient) getValidator(addr crypto.Address) (*pactus.ValidatorInfo, error) {
	res, err := c.blockchainClient.GetValidator(c.ctx,
		&pactus.GetValidatorRequest{Address: addr.String()})
	if err != nil {
		return nil, err
	}

	return res.Validator, nil
}

func (c *grpcClient) sendTx(trx *tx.Tx) (tx.ID, error) {
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

func (c *grpcClient) getFee(amount int64, payloadType payload.Type) (int64, error) {
	res, err := c.transactionClient.CalculateFee(c.ctx,
		&pactus.CalculateFeeRequest{Amount: amount, PayloadType: pactus.PayloadType(payloadType)})
	if err != nil {
		return 0, err
	}

	return res.Fee, nil
}
