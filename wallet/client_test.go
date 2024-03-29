package wallet

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type mockService struct {
	lastBlockHash     hash.Hash
	lastBlockHeight   uint32
	testAccountAddr   crypto.Address
	testValidatorAddr crypto.Address
}

func (s *mockService) GetBlockchainInfo(_ context.Context,
	_ *pactus.GetBlockchainInfoRequest,
) (*pactus.GetBlockchainInfoResponse, error) {
	return &pactus.GetBlockchainInfoResponse{
		LastBlockHeight: s.lastBlockHeight,
		LastBlockHash:   s.lastBlockHash.Bytes(),
	}, nil
}

func (s *mockService) GetConsensusInfo(_ context.Context,
	_ *pactus.GetConsensusInfoRequest,
) (*pactus.GetConsensusInfoResponse, error) {
	return &pactus.GetConsensusInfoResponse{}, nil
}

func (s *mockService) GetBlockHash(_ context.Context,
	_ *pactus.GetBlockHashRequest,
) (*pactus.GetBlockHashResponse, error) {
	return &pactus.GetBlockHashResponse{}, nil
}

func (s *mockService) GetBlockHeight(_ context.Context,
	_ *pactus.GetBlockHeightRequest,
) (*pactus.GetBlockHeightResponse, error) {
	return &pactus.GetBlockHeightResponse{}, nil
}

func (s *mockService) GetBlock(_ context.Context,
	_ *pactus.GetBlockRequest,
) (*pactus.GetBlockResponse, error) {
	return &pactus.GetBlockResponse{}, nil
}

func (s *mockService) GetAccount(_ context.Context,
	req *pactus.GetAccountRequest,
) (*pactus.GetAccountResponse, error) {
	if s.testAccountAddr.String() == req.Address {
		return &pactus.GetAccountResponse{
			Account: &pactus.AccountInfo{Balance: 1},
		}, nil
	}

	return nil, fmt.Errorf("not found")
}

func (s *mockService) GetValidatorAddresses(_ context.Context,
	_ *pactus.GetValidatorAddressesRequest,
) (*pactus.GetValidatorAddressesResponse, error) {
	return &pactus.GetValidatorAddressesResponse{}, nil
}

func (s *mockService) GetValidatorByNumber(_ context.Context,
	_ *pactus.GetValidatorByNumberRequest,
) (*pactus.GetValidatorResponse, error) {
	return &pactus.GetValidatorResponse{}, nil
}

func (s *mockService) GetValidator(_ context.Context,
	req *pactus.GetValidatorRequest,
) (*pactus.GetValidatorResponse, error) {
	if s.testAccountAddr.String() == req.Address {
		return &pactus.GetValidatorResponse{
			Validator: &pactus.ValidatorInfo{Stake: 2},
		}, nil
	}

	return nil, fmt.Errorf("not found")
}

func (s *mockService) GetPublicKey(_ context.Context,
	_ *pactus.GetPublicKeyRequest,
) (*pactus.GetPublicKeyResponse, error) {
	return &pactus.GetPublicKeyResponse{}, nil
}

func (s *mockService) GetTransaction(_ context.Context,
	_ *pactus.GetTransactionRequest,
) (*pactus.GetTransactionResponse, error) {
	return &pactus.GetTransactionResponse{}, nil
}

func (s *mockService) CalculateFee(_ context.Context,
	_ *pactus.CalculateFeeRequest,
) (*pactus.CalculateFeeResponse, error) {
	return &pactus.CalculateFeeResponse{Fee: 0}, nil
}

func (s *mockService) BroadcastTransaction(_ context.Context,
	req *pactus.BroadcastTransactionRequest,
) (*pactus.BroadcastTransactionResponse, error) {
	trx, _ := tx.FromBytes(req.SignedRawTransaction)

	return &pactus.BroadcastTransactionResponse{
		Id: trx.ID().Bytes(),
	}, nil
}

func (s *mockService) GetRawBondTransaction(_ context.Context,
	_ *pactus.GetRawBondTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	return &pactus.GetRawTransactionResponse{
		RawTransaction: make([]byte, 0),
	}, nil
}

func (s *mockService) GetRawTransferTransaction(_ context.Context,
	_ *pactus.GetRawTransferTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	return &pactus.GetRawTransactionResponse{
		RawTransaction: make([]byte, 0),
	}, nil
}

func (s *mockService) GetRawUnbondTransaction(_ context.Context,
	_ *pactus.GetRawUnbondTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	return &pactus.GetRawTransactionResponse{
		RawTransaction: make([]byte, 0),
	}, nil
}

func (s *mockService) GetRawWithdrawTransaction(_ context.Context,
	_ *pactus.GetRawWithdrawTransactionRequest,
) (*pactus.GetRawTransactionResponse, error) {
	return &pactus.GetRawTransactionResponse{
		RawTransaction: make([]byte, 0),
	}, nil
}
