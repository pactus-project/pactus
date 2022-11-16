package wallet

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/types/tx"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainServer struct{}
type transactionServer struct{}

var tBlockchainInfoResponse *pactus.GetBlockchainInfoResponse
var tAccountRequest *pactus.GetAccountRequest
var tAccountResponse *pactus.GetAccountResponse
var tValidatorRequest *pactus.GetValidatorRequest
var tValidatorResponse *pactus.GetValidatorResponse

func (s *blockchainServer) GetBlockchainInfo(_ context.Context,
	req *pactus.GetBlockchainInfoRequest) (*pactus.GetBlockchainInfoResponse, error) {
	return tBlockchainInfoResponse, nil
}

func (s *blockchainServer) GetConsensusInfo(_ context.Context,
	req *pactus.GetConsensusInfoRequest) (*pactus.GetConsensusInfoResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetBlockHash(_ context.Context,
	req *pactus.GetBlockHashRequest) (*pactus.GetBlockHashResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetBlockHeight(_ context.Context,
	req *pactus.GetBlockHeightRequest) (*pactus.GetBlockHeightResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetBlock(_ context.Context,
	req *pactus.GetBlockRequest) (*pactus.GetBlockResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetAccount(_ context.Context,
	req *pactus.GetAccountRequest) (*pactus.GetAccountResponse, error) {
	if req.Address == tAccountRequest.Address {
		return tAccountResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *blockchainServer) GetValidatorByNumber(_ context.Context,
	req *pactus.GetValidatorByNumberRequest) (*pactus.GetValidatorResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetValidator(_ context.Context,
	req *pactus.GetValidatorRequest) (*pactus.GetValidatorResponse, error) {
	if req.Address == tValidatorRequest.Address {
		return tValidatorResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *blockchainServer) GetValidators(_ context.Context,
	req *pactus.GetValidatorsRequest) (*pactus.GetValidatorsResponse, error) {
	return nil, nil
}

func (s *transactionServer) GetTransaction(ctx context.Context,
	req *pactus.GetTransactionRequest) (*pactus.GetTransactionResponse, error) {
	return nil, nil
}

func (s *transactionServer) SendRawTransaction(ctx context.Context,
	req *pactus.SendRawTransactionRequest) (*pactus.SendRawTransactionResponse, error) {
	trx, _ := tx.FromBytes(req.Data)
	return &pactus.SendRawTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}
