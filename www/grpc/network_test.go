package grpc

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/version"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetNetworkInfo(t *testing.T) {
	conn, client := callNetworkServer(t)

	t.Run("Should return node PeerID", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &pactus.NetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, []byte(tMockSync.SelfID()), res.SelfId)
		assert.Equal(t, 2, len(res.Peers))
	})

	t.Run("Should return peer info", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &pactus.NetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.Peers))
		for _, p := range res.Peers {
			if p.Moniker == "test-1" {
				assert.NotEmpty(t, p.PeerId)
				pid, _ := peer.IDFromBytes(p.PeerId)
				pp := tMockSync.PeerSet.GetPeer(pid)
				assert.Equal(t, p.Agent, pp.Agent)
				assert.Equal(t, p.Moniker, pp.Moniker)
				assert.Equal(t, p.Height, pp.Height)
				assert.Equal(t, p.PublicKey, pp.PublicKey.String())
				break
			}
		}
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetPeerInfo(t *testing.T) {
	conn, client := callNetworkServer(t)

	res, err := client.GetPeerInfo(tCtx, &pactus.PeerInfoRequest{})
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, version.Agent(), res.Agent)
	assert.Equal(t, []byte(tMockSync.SelfID()), res.PeerId)
	assert.Equal(t, tMockSync.PublicKey().String(), res.PublicKey)
	assert.Equal(t, "test-moniker", res.Moniker)

	assert.Nil(t, conn.Close(), "Error closing connection")
}
