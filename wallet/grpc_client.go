package wallet

import (
	"context"
	"encoding/hex"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"google.golang.org/grpc"
)

type grpcClient struct {
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func gewGRPCClient(rpcEndpoint string) (*grpcClient, error) {
	conn, err := grpc.Dial(rpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
	}, nil
}

func (c *grpcClient) getStamp() (hash.Stamp, error) {
	info, err := c.blockchainClient.GetBlockchainInfo(context.Background(), &pactus.BlockchainInfoRequest{})
	if err != nil {
		return hash.Stamp{}, err
	}
	h, _ := hash.FromBytes(info.LastBlockHash)
	return h.Stamp(), nil
}

func (c *grpcClient) getAccountBalance(addr crypto.Address) (int64, error) {
	acc, err := c.blockchainClient.GetAccount(context.Background(), &pactus.AccountRequest{Address: addr.String()})
	if err != nil {
		return 0, err
	}

	return acc.Account.Balance, nil
}

func (c *grpcClient) getAccountSequence(addr crypto.Address) (int32, error) {
	acc, err := c.blockchainClient.GetAccount(context.Background(), &pactus.AccountRequest{Address: addr.String()})
	if err != nil {
		return 0, err
	}

	return acc.Account.Sequence + 1, nil
}

func (c *grpcClient) GetValidatorSequence(addr crypto.Address) (int32, error) {
	val, err := c.blockchainClient.GetValidator(context.Background(), &pactus.ValidatorRequest{Address: addr.String()})
	if err != nil {
		return 0, err
	}

	return val.Validator.Sequence + 1, nil
}

func (c *grpcClient) getValidatorStake(addr crypto.Address) (int64, error) {
	val, err := c.blockchainClient.GetValidator(context.Background(), &pactus.ValidatorRequest{Address: addr.String()})
	if err != nil {
		return 0, err
	}

	return val.Validator.Stake, nil
}

func (c *grpcClient) sendTx(payload []byte) (string, error) {
	res, err := c.transactionClient.SendRawTransaction(context.Background(), &pactus.SendRawTransactionRequest{
		Data: hex.EncodeToString(payload),
	})

	if err != nil {
		return "", err
	}

	return res.Id, nil
}
