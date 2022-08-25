package grpc

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetNetworkInfo(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return node PeerID", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, []byte(tMockSync.SelfID()), res.SelfId)
		assert.Equal(t, 2, len(res.Peers))
	})

	t.Run("Should return peer info", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
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

func TestGetBlockchainInfo(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return the last block height", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, tMockState.TestStore.LastHeight, res.LastBlockHeight)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
