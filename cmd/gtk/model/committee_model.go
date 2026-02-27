//go:build gtk

package model

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// CommitteeModel holds blockchain gRPC client and provides committee data
// for the committee widget (size, power, total power, validators, protocol versions).
type CommitteeModel struct {
	ctx              context.Context
	blockchainClient pactus.BlockchainClient
}

// NewCommitteeModel creates a CommitteeModel that uses gRPC to fetch committee data.
func NewCommitteeModel(
	ctx context.Context,
	blockchainClient pactus.BlockchainClient,
) *CommitteeModel {
	return &CommitteeModel{
		ctx:              ctx,
		blockchainClient: blockchainClient,
	}
}

// GetCommitteeInfo returns current committee info from gRPC.
func (m *CommitteeModel) GetCommitteeInfo() (*pactus.GetCommitteeInfoResponse, error) {
	return m.blockchainClient.GetCommitteeInfo(m.ctx, &pactus.GetCommitteeInfoRequest{})
}
