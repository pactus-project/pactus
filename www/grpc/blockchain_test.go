package grpc

import (
	"context"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	height := uint32(100)
	b := td.mockState.TestStore.AddTestBlock(height)
	data, _ := b.Bytes()

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{
				Height: height + 1, Verbosity: pactus.BlockVerbosity_BLOCK_DATA,
			})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block data (verbosity: 0)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_DATA})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return object with  (verbosity: 1)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_INFO})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Empty(t, res.Data)
		assert.NotEmpty(t, res.Header)
		assert.Equal(t, res.PrevCert.Committers, b.PrevCertificate().Committers())
		assert.Equal(t, res.PrevCert.Absentees, b.PrevCertificate().Absentees())
		for i, trx := range res.Txs {
			blockTrx := b.Transactions()[i]
			data, _ := blockTrx.Bytes()

			assert.Equal(t, blockTrx.ID().Bytes(), trx.Id)
			assert.Equal(t, data, trx.Data)
			assert.Zero(t, trx.LockTime)
			assert.Empty(t, trx.Signature)
			assert.Empty(t, trx.PublicKey)
		}
	})

	t.Run("Should return object with  (verbosity: 2)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Empty(t, res.Data)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
		for i, trx := range res.Txs {
			blockTrx := b.Transactions()[i]

			assert.Equal(t, blockTrx.ID().Bytes(), trx.Id)
			assert.Empty(t, trx.Data)
			assert.Equal(t, blockTrx.LockTime(), trx.LockTime)
			assert.Equal(t, blockTrx.Signature().Bytes(), trx.Signature)
			assert.Equal(t, blockTrx.PublicKey().String(), trx.PublicKey)
		}
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetBlockHash(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	b := td.mockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(context.Background(),
			&pactus.GetBlockHashRequest{Height: 0})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(context.Background(),
			&pactus.GetBlockHashRequest{Height: 100})

		assert.NoError(t, err)
		assert.Equal(t, b.Hash().Bytes(), res.Hash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetBlockHeight(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	b := td.mockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for invalid hash", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: nil})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: td.RandHash().Bytes()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: b.Hash().Bytes()})

		assert.NoError(t, err)
		assert.Equal(t, uint32(100), res.Height)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetBlockchainInfo(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	t.Run("Should return the last block height", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(context.Background(),
			&pactus.GetBlockchainInfoRequest{})

		assert.NoError(t, err)
		assert.Equal(t, td.mockState.TestStore.LastHeight, res.LastBlockHeight)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetAccount(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	acc, addr := td.mockState.TestStore.AddTestAccount()

	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetAccount(context.Background(),
			&pactus.GetAccountRequest{Address: ""})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(context.Background(),
			&pactus.GetAccountRequest{Address: td.RandAccAddress().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccount(context.Background(),
			&pactus.GetAccountRequest{Address: addr.String()})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Balance, acc.Balance().ToNanoPAC())
		assert.Equal(t, res.Account.Number, acc.Number())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetValidator(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	val1 := td.mockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(context.Background(),
			&pactus.GetValidatorRequest{Address: ""})

		assert.Error(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(context.Background(),
			&pactus.GetValidatorRequest{Address: td.RandAccAddress().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(context.Background(),
			&pactus.GetValidatorRequest{Address: val1.Address().String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetValidatorByNumber(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	val1 := td.mockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(context.Background(),
			&pactus.GetValidatorByNumberRequest{Number: -1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(context.Background(),
			&pactus.GetValidatorByNumberRequest{Number: val1.Number() + 1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator matching with public key and number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(context.Background(),
			&pactus.GetValidatorByNumberRequest{Number: val1.Number()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
		assert.Equal(t, val1.Number(), res.GetValidator().GetNumber())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetValidatorAddresses(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	t.Run("should return list of validator addresses", func(t *testing.T) {
		td.mockState.TestStore.AddTestValidator()
		td.mockState.TestStore.AddTestValidator()

		res, err := client.GetValidatorAddresses(context.Background(),
			&pactus.GetValidatorAddressesRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res.GetAddresses()))
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetPublicKey(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	val := td.mockState.TestStore.AddTestValidator()

	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetPublicKey(context.Background(),
			&pactus.GetPublicKeyRequest{Address: ""})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing public key ", func(t *testing.T) {
		res, err := client.GetPublicKey(context.Background(),
			&pactus.GetPublicKeyRequest{Address: td.RandAccAddress().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return the public key", func(t *testing.T) {
		res, err := client.GetPublicKey(context.Background(),
			&pactus.GetPublicKeyRequest{Address: val.Address().String()})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.PublicKey, val.PublicKey().String())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestConsensusInfo(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.blockchainClient(t)

	v1, _ := td.GenerateTestPrepareVote(100, 2)
	v2, _ := td.GenerateTestPrepareVote(100, 2)
	td.consMocks[1].Active = true
	td.consMocks[1].Height = 100
	td.consMocks[0].AddVote(v1)
	td.consMocks[1].AddVote(v2)

	t.Run("Should return the consensus info", func(t *testing.T) {
		res, err := client.GetConsensusInfo(context.Background(), &pactus.GetConsensusInfoRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.False(t, res.Instances[0].Active, true)
		assert.True(t, res.Instances[1].Active, true)
		assert.Equal(t, res.Instances[1].Height, uint32(100))
		assert.Equal(t, res.Instances[0].Votes[0].Type, pactus.VoteType_VOTE_PREPARE)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
