//go:build gtk

package model

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// NodeModel holds blockchain and network gRPC clients and provides node/chain
// data for the node widget. No controller uses node.Node; they use this model.
type NodeModel struct {
	ctx              context.Context
	blockchainClient pactus.BlockchainClient
	networkClient    pactus.NetworkClient
}

// NewNodeModel creates a NodeModel that uses gRPC to fetch node and chain data.
func NewNodeModel(
	ctx context.Context,
	blockchainClient pactus.BlockchainClient,
	networkClient pactus.NetworkClient,
) *NodeModel {
	return &NodeModel{
		ctx:              ctx,
		blockchainClient: blockchainClient,
		networkClient:    networkClient,
	}
}

// GetBlockchainInfo returns current blockchain info.
func (m *NodeModel) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	return m.blockchainClient.GetBlockchainInfo(m.ctx, &pactus.GetBlockchainInfoRequest{})
}

// GetCommitteeInfo returns current committee info.
func (m *NodeModel) GetCommitteeInfo() (*pactus.GetCommitteeInfoResponse, error) {
	return m.blockchainClient.GetCommitteeInfo(m.ctx, &pactus.GetCommitteeInfoRequest{})
}

// GetConsensusInfo returns consensus instances (used to determine "in committee").
func (m *NodeModel) GetConsensusInfo() (*pactus.GetConsensusInfoResponse, error) {
	return m.blockchainClient.GetConsensusInfo(m.ctx, &pactus.GetConsensusInfoRequest{})
}

// GetNodeInfo returns this node's info (moniker, peer ID, reachability, clock offset, connections).
func (m *NodeModel) GetNodeInfo() (*pactus.GetNodeInfoResponse, error) {
	return m.networkClient.GetNodeInfo(m.ctx, &pactus.GetNodeInfoRequest{})
}
