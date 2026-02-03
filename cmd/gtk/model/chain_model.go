//go:build gtk

package model

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// ChainModel holds blockchain gRPC client and provides chain data for the node widget.
type ChainModel struct {
	ctx              context.Context
	blockchainClient pactus.BlockchainClient
}

// NewChainModel creates a ChainModel that uses gRPC to fetch chain data.
func NewChainModel(ctx context.Context, blockchainClient pactus.BlockchainClient) *ChainModel {
	return &ChainModel{
		ctx:              ctx,
		blockchainClient: blockchainClient,
	}
}

// GetBlockchainInfo returns current blockchain info (height, block hash, etc.).
func (m *ChainModel) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	return m.blockchainClient.GetBlockchainInfo(m.ctx, &pactus.GetBlockchainInfoRequest{})
}
