package grpc

import (
	"context"
	"encoding/hex"
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.networkClient(t)

	t.Run("Should return node PeerID", func(t *testing.T) {
		res, err := client.GetNetworkInfo(context.Background(),
			&pactus.GetNetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.ConnectedPeers))
	})

	t.Run("Should return peer info", func(t *testing.T) {
		res, err := client.GetNetworkInfo(context.Background(),
			&pactus.GetNetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.ConnectedPeers))
		for _, peer := range res.ConnectedPeers {
			assert.NotEmpty(t, peer.PeerId)
			b, err := hex.DecodeString(peer.PeerId)
			assert.NoError(t, err)
			pid, _ := lp2ppeer.IDFromBytes(b)
			pp := td.mockSync.PeerSet().GetPeer(pid)
			assert.Equal(t, peer.Agent, pp.Agent)
			assert.Equal(t, peer.Moniker, pp.Moniker)
			assert.Equal(t, peer.Height, pp.Height)
			assert.NotEmpty(t, pp.ConsensusKeys)
			for _, key := range pp.ConsensusKeys {
				assert.Contains(t, peer.ConsensusKeys, key.String())
			}
		}
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetNodeInfo(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.networkClient(t)

	res, err := client.GetNodeInfo(context.Background(),
		&pactus.GetNodeInfoRequest{})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, version.NodeAgent.String(), res.Agent)
	assert.Equal(t, hex.EncodeToString([]byte(td.mockSync.SelfID())), res.PeerId)
	assert.Equal(t, "test-moniker", res.Moniker)
	assert.Equal(t, res.ZmqPublishers[0].Address, "zmq_address")
	assert.Equal(t, res.ZmqPublishers[0].Topic, "zmq_topic")
	assert.Equal(t, res.ZmqPublishers[0].Hwm, int32(100))

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
func TestPing(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)
	conn, client := td.networkClient(t)

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

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
