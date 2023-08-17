package wallet

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/types/tx"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainServer struct{}

type transactionServer struct{}

var (
	tBlockchainInfoResponse *pactus.GetBlockchainInfoResponse
	tAccountResponse        *pactus.GetAccountResponse
	tValidatorResponse      *pactus.GetValidatorResponse
)

func (s *blockchainServer) GetBlockchainInfo(_ context.Context,
	_ *pactus.GetBlockchainInfoRequest,
) (*pactus.GetBlockchainInfoResponse, error) {
	return tBlockchainInfoResponse, nil
}

func (s *blockchainServer) GetConsensusInfo(_ context.Context,
	_ *pactus.GetConsensusInfoRequest,
) (*pactus.GetConsensusInfoResponse, error) {
	return &pactus.GetConsensusInfoResponse{}, nil
}

func (s *blockchainServer) GetBlockHash(_ context.Context,
	_ *pactus.GetBlockHashRequest,
) (*pactus.GetBlockHashResponse, error) {
	return &pactus.GetBlockHashResponse{}, nil
}

func (s *blockchainServer) GetBlockHeight(_ context.Context,
	_ *pactus.GetBlockHeightRequest,
) (*pactus.GetBlockHeightResponse, error) {
	return &pactus.GetBlockHeightResponse{}, nil
}

func (s *blockchainServer) GetBlock(_ context.Context,
	_ *pactus.GetBlockRequest,
) (*pactus.GetBlockResponse, error) {
	return &pactus.GetBlockResponse{}, nil
}

func (s *blockchainServer) GetAccount(_ context.Context,
	_ *pactus.GetAccountRequest,
) (*pactus.GetAccountResponse, error) {
	if tAccountResponse != nil {
		return tAccountResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *blockchainServer) GetValidatorAddresses(_ context.Context,
	_ *pactus.GetValidatorAddressesRequest,
) (*pactus.GetValidatorAddressesResponse, error) {
	return &pactus.GetValidatorAddressesResponse{}, nil
}

func (s *blockchainServer) GetAccountByNumber(_ context.Context,
	_ *pactus.GetAccountByNumberRequest,
) (*pactus.GetAccountResponse, error) {
	return &pactus.GetAccountResponse{}, nil
}

func (s *blockchainServer) GetValidatorByNumber(_ context.Context,
	_ *pactus.GetValidatorByNumberRequest,
) (*pactus.GetValidatorResponse, error) {
	return &pactus.GetValidatorResponse{}, nil
}

func (s *blockchainServer) GetValidator(_ context.Context,
	_ *pactus.GetValidatorRequest,
) (*pactus.GetValidatorResponse, error) {
	if tValidatorResponse != nil {
		return tValidatorResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *transactionServer) GetTransaction(_ context.Context,
	_ *pactus.GetTransactionRequest,
) (*pactus.GetTransactionResponse, error) {
	return &pactus.GetTransactionResponse{}, nil
}

func (s *transactionServer) CalculateFee(_ context.Context,
	_ *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	return &pactus.CalculateFeeResponse{Fee: 0}, nil
}

func (s *transactionServer) SendRawTransaction(_ context.Context,
	req *pactus.SendRawTransactionRequest,
) (*pactus.SendRawTransactionResponse, error) {
	trx, _ := tx.FromBytes(req.Data)
	return &pactus.SendRawTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}
