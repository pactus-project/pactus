package integration

import (
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodInfo(t *testing.T) {
	res, err := tNetwork.GetNodeInfo(t.Context(), &pactus.GetNodeInfoRequest{})
	require.NoError(t, err)

	assert.Equal(t, "test", res.NetworkName)
}

func TestNetworkInfo(t *testing.T) {
	res, err := tNetwork.GetNetworkInfo(t.Context(), &pactus.GetNetworkInfoRequest{})
	require.NoError(t, err)

	assert.Greater(t, res.ConnectedPeersCount, uint32(1))
}

func TestListPeer(t *testing.T) {
	res, err := tNetwork.ListPeers(t.Context(), &pactus.ListPeersRequest{})
	require.NoError(t, err)

	assert.Greater(t, len(res.Peers), 1)
}
