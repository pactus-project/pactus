//go:build gtk

package model

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// NetworkModel holds network gRPC client and provides network and peer data
// for the network widget.
type NetworkModel struct {
	ctx           context.Context
	networkClient pactus.NetworkClient
}

// NewNetworkModel creates a NetworkModel that uses gRPC to fetch network data.
func NewNetworkModel(
	ctx context.Context,
	networkClient pactus.NetworkClient,
) *NetworkModel {
	return &NetworkModel{
		ctx:           ctx,
		networkClient: networkClient,
	}
}

// GetNetworkInfo returns overall network info from gRPC.
func (m *NetworkModel) GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error) {
	return m.networkClient.GetNetworkInfo(m.ctx, &pactus.GetNetworkInfoRequest{})
}

// ListPeers returns active (connected) peers only. Set includeDisconnected true to include disconnected.
func (m *NetworkModel) ListPeers(includeDisconnected bool) (*pactus.ListPeersResponse, error) {
	return m.networkClient.ListPeers(m.ctx, &pactus.ListPeersRequest{
		IncludeDisconnected: includeDisconnected,
	})
}
