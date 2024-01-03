package grpc

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	conn, client := testNetworkClient(t)

	t.Run("Should return node PeerID", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &pactus.GetNetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.ConnectedPeers))
	})

	t.Run("Should return peer info", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &pactus.GetNetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.ConnectedPeers))
		for _, p := range res.ConnectedPeers {
			assert.NotEmpty(t, p.PeerId)
			pid, _ := peer.IDFromBytes(p.PeerId)
			pp := tMockSync.PeerSet().GetPeer(pid)
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
}

func TestGetNodeInfo(t *testing.T) {
	conn, client := testNetworkClient(t)

	res, err := client.GetNodeInfo(tCtx, &pactus.GetNodeInfoRequest{})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, version.Agent(), res.Agent)
	assert.Equal(t, []byte(tMockSync.SelfID()), res.PeerId)
	assert.Equal(t, "test-moniker", res.Moniker)

	assert.Nil(t, conn.Close(), "Error closing connection")
}
