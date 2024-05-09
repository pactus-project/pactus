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
		for _, p := range res.ConnectedPeers {
			assert.NotEmpty(t, p.PeerId)
			pid, _ := lp2ppeer.IDFromBytes(p.PeerId)
			pp := td.mockSync.PeerSet().GetPeer(pid)
			assert.Equal(t, p.Agent, pp.Agent)
			assert.Equal(t, p.Moniker, pp.Moniker)
			assert.Equal(t, p.Height, pp.Height)
			assert.NotEmpty(t, pp.ConsensusKeys)
			for _, key := range pp.ConsensusKeys {
				assert.Contains(t, p.ConsensusKeys, key.String())
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
	assert.Equal(t, []byte(td.mockSync.SelfID()), res.PeerId)
	assert.Equal(t, "test-moniker", res.Moniker)

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
