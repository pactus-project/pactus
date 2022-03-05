package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetBlock(t *testing.T) {
	conn, client := callServer(t)

	tMockState.CommitTestBlocks(10)

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1000, Verbosity: zarb.BlockVerbosity_BLOCK_HASH})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block hash", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 5, Verbosity: zarb.BlockVerbosity_BLOCK_HASH})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		h, err := hash.FromString(res.Hash)
		assert.NoError(t, err)
		assert.NoError(t, h.SanityCheck())
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Tranactions)
	})

	t.Run("Should return json object with verbosity 1 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1, Verbosity: zarb.BlockVerbosity_BLOCK_INFO})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		h, err := hash.FromString(res.Hash)
		assert.NoError(t, err)
		assert.NoError(t, h.SanityCheck())
		assert.NotEmpty(t, res.Header)
		assert.Empty(t, res.Tranactions)
	})

	t.Run("Should return object with verbosity 2 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1, Verbosity: zarb.BlockVerbosity_BLOCK_TRANSACTIONS})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		h, err := hash.FromString(res.Hash)
		assert.NoError(t, err)
		assert.NoError(t, h.SanityCheck())
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Tranactions)
	})

	conn.Close()
}

func TestGetBlockHieght(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return InvalidArgument for invalid hash", func(t *testing.T) {
		res, err := client.GetBlockHeight(tCtx, &zarb.BlockHeightRequest{Hash: "NOt A  valid has"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Hash provided is not Valid")
		assert.Nil(t, res)
	})

	t.Run("Should return NotFound for non existing block ", func(t *testing.T) {
		res, err := client.GetBlockHeight(tCtx, &zarb.BlockHeightRequest{Hash: hash.GenerateTestHash().String()})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "No block found with the Hash provided")
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		b5, _ := tMockState.Store.Block(5)
		res, err := client.GetBlockHeight(tCtx, &zarb.BlockHeightRequest{Hash: b5.Hash().String()})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Height)
	})

	conn.Close()
}
