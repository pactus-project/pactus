package grpc

import (
	"context"
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.networkClient(t)

	res, err := client.GetNetworkInfo(context.Background(),
		&pactus.GetNetworkInfoRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.NetworkName)
}

func TestListPeers(t *testing.T) {
	td := setup(t, nil)
	client := td.networkClient(t)

	res, err := client.ListPeers(context.Background(),
		&pactus.ListPeersRequest{IncludeDisconnected: true})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res.Peers))
	for _, peer := range res.Peers {
		assert.NotEmpty(t, peer.PeerId)
		assert.NoError(t, err)
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

	res, err := client.GetNodeInfo(context.Background(),
		&pactus.GetNodeInfoRequest{})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, version.NodeAgent.String(), res.Agent)
	assert.Equal(t, td.mockSync.SelfID().String(), res.PeerId)
	assert.Equal(t, "test-moniker", res.Moniker)
	assert.Equal(t, res.ZmqPublishers[0].Address, "zmq_address")
	assert.Equal(t, res.ZmqPublishers[0].Topic, "zmq_topic")
	assert.Equal(t, res.ZmqPublishers[0].Hwm, int32(100))
}

func TestPing(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)
	client := td.networkClient(t)

	t.Run("Should return empty response for ping", func(t *testing.T) {
		res, err := client.Ping(context.Background(), &pactus.PingRequest{})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.IsType(t, &pactus.PingResponse{}, res)
	})

	t.Run("Should handle multiple ping requests", func(t *testing.T) {
		// Test multiple consecutive pings to ensure consistency
		for i := 0; i < 5; i++ {
			res, err := client.Ping(context.Background(), &pactus.PingRequest{})
			assert.Nil(t, err)
			assert.NotNil(t, res)
		}
	})
}
