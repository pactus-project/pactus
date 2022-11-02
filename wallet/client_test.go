package wallet

import (
	"context"
	"fmt"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainServer struct{}
type networkServer struct{}

var tBlockchainInfoResponse *pactus.BlockchainInfoResponse
var tAccountRequest *pactus.AccountRequest
var tAccountResponse *pactus.AccountResponse
var tValidatorRequest *pactus.ValidatorRequest
var tValidatorResponse *pactus.ValidatorResponse

func (s *blockchainServer) GetBlockchainInfo(_ context.Context,
	req *pactus.BlockchainInfoRequest) (*pactus.BlockchainInfoResponse, error) {
	return tBlockchainInfoResponse, nil
}

func (s *blockchainServer) GetBlockHash(_ context.Context,
	req *pactus.BlockHashRequest) (*pactus.BlockHashResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetBlockHeight(_ context.Context,
	req *pactus.BlockHeightRequest) (*pactus.BlockHeightResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetBlock(_ context.Context,
	req *pactus.BlockRequest) (*pactus.BlockResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetAccount(_ context.Context,
	req *pactus.AccountRequest) (*pactus.AccountResponse, error) {
	if req.Address == tAccountRequest.Address {
		return tAccountResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *blockchainServer) GetValidatorByNumber(_ context.Context,
	req *pactus.ValidatorByNumberRequest) (*pactus.ValidatorResponse, error) {
	return nil, nil
}

func (s *blockchainServer) GetValidator(_ context.Context,
	req *pactus.ValidatorRequest) (*pactus.ValidatorResponse, error) {
	if req.Address == tValidatorRequest.Address {
		return tValidatorResponse, nil
	}
	return nil, fmt.Errorf("unknown request")
}

func (s *blockchainServer) GetValidators(_ context.Context,
	req *pactus.ValidatorsRequest) (*pactus.ValidatorsResponse, error) {
	return nil, nil
}

func (s *networkServer) GetNetworkInfo(_ context.Context,
	req *pactus.NetworkInfoRequest) (*pactus.NetworkInfoResponse, error) {
	return nil, nil
}

func (s *networkServer) GetPeerInfo(_ context.Context,
	req *pactus.PeerInfoRequest) (*pactus.PeerInfoResponse, error) {
	return nil, nil
}
