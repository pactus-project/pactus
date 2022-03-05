package grpc

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetNetworkInfo(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return node PeerID", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, tMockSync.SelfID().String(), res.SelfId)
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
				assert.Equal(t, p.Height, int32(pp.Height))
				assert.Equal(t, p.PublicKey, pp.PublicKey.String())
				break
			}
		}
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}

func TestGetBlockchainInfo(t *testing.T) {
	conn, client := callServer(t)
	tMockState.Store.Blocks = make(map[int]*block.Block)

	t.Run("Should return 0,for no block yet", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), res.Height)
		assert.Equal(t, hash.UndefHash.String(), res.LastBlockHash)
	})

	tMockState.CommitTestBlocks(10)
	t.Run("Should return 1, for first block", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Height)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	t.Run("Should return 5", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), res.Height)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}
