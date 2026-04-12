package grpc

import (
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.networkClient(t)

	res, err := client.GetNetworkInfo(t.Context(),
		&pactus.GetNetworkInfoRequest{})
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.NetworkName)
}

func TestListPeers(t *testing.T) {
	td := setup(t, nil)
	client := td.networkClient(t)

	res, err := client.ListPeers(t.Context(),
		&pactus.ListPeersRequest{IncludeDisconnected: true})
	require.NoError(t, err)
	assert.Len(t, res.Peers, 2)
	for _, peer := range res.Peers {
		assert.NotEmpty(t, peer.PeerId)
		require.NoError(t, err)
		pid, _ := lp2ppeer.Decode(peer.PeerId)
		pp := td.mockSync.PeerSet().GetPeer(pid)
		assert.Equal(t, peer.Agent, pp.Agent)
		assert.Equal(t, peer.Moniker, pp.Moniker)
		assert.Equal(t, peer.Height, pp.Height)
		assert.NotEmpty(t, pp.ConsensusKeys)
		for _, key := range pp.ConsensusKeys {
			assert.Contains(t, peer.ConsensusKeys, key.String())
		}
	}
}

func TestGetNodeInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.networkClient(t)

	res, err := client.GetNodeInfo(t.Context(),
		&pactus.GetNodeInfoRequest{})
	require.NoError(t, err)
	assert.Equal(t, version.NodeAgent.String(), res.Agent)
	assert.Equal(t, td.mockSync.SelfID().String(), res.PeerId)
	assert.Equal(t, "test-moniker", res.Moniker)
	assert.Equal(t, "zmq_address", res.ZmqPublishers[0].Address)
	assert.Equal(t, "zmq_topic", res.ZmqPublishers[0].Topic)
	assert.Equal(t, int32(100), res.ZmqPublishers[0].Hwm)
}

func TestPing(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)
	client := td.networkClient(t)

	t.Run("Should return empty response for ping", func(t *testing.T) {
		res, err := client.Ping(t.Context(), &pactus.PingRequest{})

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.IsType(t, &pactus.PingResponse{}, res)
	})

	t.Run("Should handle multiple ping requests", func(t *testing.T) {
		// Test multiple consecutive pings to ensure consistency
		for i := 0; i < 5; i++ {
			res, err := client.Ping(t.Context(), &pactus.PingRequest{})
			require.NoError(t, err)
			assert.NotNil(t, res)
		}
	})
}
