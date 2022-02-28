package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
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

	t.Run("Should return newly added Peer", func(t *testing.T) {
		res, err := client.GetNetworkInfo(tCtx, &zarb.NetworkInfoRequest{})
		assert.NoError(t, err)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res.Peers))
		for _, p := range res.Peers {
			if p.Moniker == "test-1" {
				assert.NotEmpty(t, p.PeerId)
				assert.NotZero(t, p.Height)
				return
			}
		}
		t.Error("new Peer Not Found")
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
		assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", res.LastBlockHash)
	})

	b1, trxs := block.GenerateTestBlock(nil, nil)
	tMockState.AddBlock(1, b1, trxs)
	t.Run("Should return 1,for first block", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Height)
		assert.Equal(t, b1.Hash().String(), res.LastBlockHash)
	})
	b2, trxs2 := block.GenerateTestBlock(nil, nil)
	b3, trxs3 := block.GenerateTestBlock(nil, nil)
	b4, trxs4 := block.GenerateTestBlock(nil, nil)
	b5, trxs5 := block.GenerateTestBlock(nil, nil)
	tMockState.AddBlock(2, b2, trxs2)
	tMockState.AddBlock(3, b3, trxs3)
	tMockState.AddBlock(4, b4, trxs4)
	tMockState.AddBlock(5, b5, trxs5)

	t.Run("Should return 5", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &zarb.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), res.Height)
		assert.Equal(t, b5.Hash().String(), res.LastBlockHash)
	})

	err := conn.Close()

	assert.Nil(t, err, "Error closing connection")
}
