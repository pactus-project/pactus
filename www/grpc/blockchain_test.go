package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
	conn, client := callBlockchainServer(t)

	height := uint32(100)
	b := tMockState.TestStore.AddTestBlock(height)
	data, _ := b.Bytes()

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &pactus.BlockRequest{Hash: hash.GenerateTestHash().Bytes(), Verbosity: pactus.BlockVerbosity_BLOCK_DATA})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block data", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &pactus.BlockRequest{Hash: b.Hash().Bytes(), Verbosity: pactus.BlockVerbosity_BLOCK_DATA})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return object with verbosity 1 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &pactus.BlockRequest{Hash: b.Hash().Bytes(), Verbosity: pactus.BlockVerbosity_BLOCK_INFO})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
		assert.Equal(t, res.PrevCert.Committers, b.PrevCertificate().Committers())
		assert.Equal(t, res.PrevCert.Absentees, b.PrevCertificate().Absentees())
	})

	t.Run("Should return object with verbosity 2 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &pactus.BlockRequest{Hash: b.Hash().Bytes(), Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockHash(t *testing.T) {
	conn, client := callBlockchainServer(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for non existing block ", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx, &pactus.BlockHashRequest{Height: 101})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx, &pactus.BlockHashRequest{Height: 100})
		assert.NoError(t, err)
		assert.Equal(t, b.Hash().Bytes(), res.Hash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockchainInfo(t *testing.T) {
	conn, client := callBlockchainServer(t)

	t.Run("Should return the last block height", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx, &pactus.BlockchainInfoRequest{})
		assert.NoError(t, err)
		assert.Equal(t, tMockState.TestStore.LastHeight, res.LastBlockHeight)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetAccount(t *testing.T) {
	conn, client := callBlockchainServer(t)
	acc := tMockState.TestStore.AddTestAccount()

	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &pactus.AccountRequest{
			Address: "",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &pactus.AccountRequest{
			Address: crypto.GenerateTestAddress().String(),
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &pactus.AccountRequest{
			Address: acc.Address().String(),
		})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Address, acc.Address().String())
		assert.Equal(t, res.Account.Balance, acc.Balance())
		assert.Equal(t, res.Account.Number, acc.Number())
		assert.Equal(t, res.Account.Sequence, acc.Sequence())
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}
func TestGetValidator(t *testing.T) {
	conn, client := callBlockchainServer(t)

	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &pactus.ValidatorRequest{
			Address: "",
		})
		assert.Error(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &pactus.ValidatorRequest{
			Address: crypto.GenerateTestAddress().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(tCtx, &pactus.ValidatorRequest{
			Address: val1.Address().String(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := callBlockchainServer(t)

	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &pactus.ValidatorByNumberRequest{
			Number: -1,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &pactus.ValidatorByNumberRequest{
			Number: val1.Number() + 1,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator matching with public key and number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx, &pactus.ValidatorByNumberRequest{
			Number: val1.Number(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
		assert.Equal(t, val1.Number(), res.GetValidator().GetNumber())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidators(t *testing.T) {
	conn, client := callBlockchainServer(t)

	t.Run("should return list of validators", func(t *testing.T) {
		res, err := client.GetValidators(tCtx, &pactus.ValidatorsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 21, len(res.GetValidators()))
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
