package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetBlock(t *testing.T) {
	conn, client := callServer(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Hash: hash.GenerateTestHash().RawBytes(), Verbosity: zarb.BlockVerbosity_BLOCK_HASH})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block hash", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Hash: b.Hash().RawBytes(), Verbosity: zarb.BlockVerbosity_BLOCK_HASH})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		//assert.Equal(t, res.Height, 100)
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return json object with verbosity 1 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Hash: b.Hash().RawBytes(), Verbosity: zarb.BlockVerbosity_BLOCK_INFO})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		//assert.Equal(t, res.Height, 100)
		assert.NotEmpty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return object with verbosity 2 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Hash: b.Hash().RawBytes(), Verbosity: zarb.BlockVerbosity_BLOCK_TRANSACTIONS})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		//assert.Equal(t, res.Height, 100)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockHash(t *testing.T) {
	conn, client := callServer(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for invalid height", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx, &zarb.BlockHashRequest{Height: -1})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return error for non existing block ", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx, &zarb.BlockHashRequest{Height: 101})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx, &zarb.BlockHashRequest{Height: 100})
		assert.NoError(t, err)
		assert.Equal(t, b.Hash().RawBytes(), res.Hash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
