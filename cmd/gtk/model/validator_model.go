//go:build gtk

package model

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// ValidatorModel holds blockchain gRPC client and provides validator data
// for the validator widget (this node's consensus instances).
type ValidatorModel struct {
	ctx              context.Context
	blockchainClient pactus.BlockchainClient
}

// NewValidatorModel creates a ValidatorModel that uses gRPC to fetch validator data.
func NewValidatorModel(
	ctx context.Context,
	blockchainClient pactus.BlockchainClient,
) *ValidatorModel {
	return &ValidatorModel{
		ctx:              ctx,
		blockchainClient: blockchainClient,
	}
}

// Context returns the model's context (e.g. for cancellation).
func (m *ValidatorModel) Context() context.Context {
	return m.ctx
}

// Validators returns validator info for this node's consensus instances.
// It calls GetConsensusInfo to get instance addresses, then GetValidator for each.
func (m *ValidatorModel) Validators() ([]*pactus.ValidatorInfo, error) {
	res, err := m.blockchainClient.GetConsensusInfo(m.ctx, &pactus.GetConsensusInfoRequest{})
	if err != nil {
		return nil, err
	}
	vals := make([]*pactus.ValidatorInfo, 0, len(res.Instances))
	for _, inst := range res.Instances {
		vres, err := m.blockchainClient.GetValidator(m.ctx, &pactus.GetValidatorRequest{
			Address: inst.Address,
		})
		if err != nil {
			continue // skip inactive validator
		}

		vals = append(vals, vres.Validator)
	}

	return vals, nil
}
