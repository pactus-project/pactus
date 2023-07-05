package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
	conn, client := testBlockchainClient(t)

	height := uint32(100)
	b := tMockState.TestStore.AddTestBlock(height)
	data, _ := b.Bytes()

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &pactus.GetBlockRequest{
			Height: height + 1, Verbosity: pactus.BlockVerbosity_BLOCK_DATA})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block data", func(t *testing.T) {
		data, _ := b.Bytes()
		res, err := client.GetBlock(tCtx,
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_DATA})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return object with verbosity 1 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx,
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_INFO})

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
		res, err := client.GetBlock(tCtx,
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_TRANSACTIONS})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Height, height)
		assert.Equal(t, res.Hash, b.Hash().Bytes())
		assert.Equal(t, res.Data, data)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
		for i, trx := range res.Txs {
			data, _ := b.Transactions()[i].Bytes()
			assert.Equal(t, b.Transactions()[i].ID().Bytes(), trx.Id)
			assert.Equal(t, b.Transactions()[i].Signature().Bytes(), trx.Signature)
			assert.Equal(t, data, trx.Data)
		}
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockHash(t *testing.T) {
	conn, client := testBlockchainClient(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx,
			&pactus.GetBlockHashRequest{Height: 0})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(tCtx,
			&pactus.GetBlockHashRequest{Height: 100})

		assert.NoError(t, err)
		assert.Equal(t, b.Hash().Bytes(), res.Hash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockHeight(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testBlockchainClient(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Should return error for invalid hash", func(t *testing.T) {
		res, err := client.GetBlockHeight(tCtx,
			&pactus.GetBlockHeightRequest{Hash: nil})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(tCtx,
			&pactus.GetBlockHeightRequest{Hash: ts.RandomHash().Bytes()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(tCtx,
			&pactus.GetBlockHeightRequest{Hash: b.Hash().Bytes()})

		assert.NoError(t, err)
		assert.Equal(t, uint32(100), res.Height)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetBlockchainInfo(t *testing.T) {
	conn, client := testBlockchainClient(t)

	t.Run("Should return the last block height", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(tCtx,
			&pactus.GetBlockchainInfoRequest{})

		assert.NoError(t, err)
		assert.Equal(t, tMockState.TestStore.LastHeight, res.LastBlockHeight)
		assert.NotEmpty(t, res.LastBlockHash)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetAccount(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testBlockchainClient(t)
	acc, signer := tMockState.TestStore.AddTestAccount()

	t.Run("Should return error for non-parsable address ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx,
			&pactus.GetAccountRequest{Address: ""})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx,
			&pactus.GetAccountRequest{Address: ts.RandomAddress().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccount(tCtx,
			&pactus.GetAccountRequest{Address: signer.Address().String()})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Balance, acc.Balance())
		assert.Equal(t, res.Account.Number, acc.Number())
		assert.Equal(t, res.Account.Sequence, acc.Sequence())
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetAccountByNumber(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testBlockchainClient(t)
	acc, _ := tMockState.TestStore.AddTestAccount()

	t.Run("Should return nil value due to invalid number ", func(t *testing.T) {
		res, err := client.GetAccountByNumber(tCtx,
			&pactus.GetAccountByNumberRequest{Number: -1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccountByNumber(tCtx,
			&pactus.GetAccountByNumberRequest{Number: ts.RandInt32(1000)})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return account details", func(t *testing.T) {
		res, err := client.GetAccountByNumber(tCtx,
			&pactus.GetAccountByNumberRequest{Number: acc.Number()})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Account.Balance, acc.Balance())
		assert.Equal(t, res.Account.Number, acc.Number())
		assert.Equal(t, res.Account.Sequence, acc.Sequence())
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testBlockchainClient(t)
	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid address", func(t *testing.T) {
		res, err := client.GetValidator(tCtx,
			&pactus.GetValidatorRequest{Address: ""})

		assert.Error(t, err, "Error should be returned")
		assert.Nil(t, res, "Response should be empty")
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidator(tCtx,
			&pactus.GetValidatorRequest{Address: ts.RandomAddress().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator, and the public keys should match", func(t *testing.T) {
		res, err := client.GetValidator(tCtx,
			&pactus.GetValidatorRequest{Address: val1.Address().String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidatorByNumber(t *testing.T) {
	conn, client := testBlockchainClient(t)

	val1 := tMockState.TestStore.AddTestValidator()

	t.Run("Should return nil value due to invalid number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx,
			&pactus.GetValidatorByNumberRequest{Number: -1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("should return Not Found", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx,
			&pactus.GetValidatorByNumberRequest{Number: val1.Number() + 1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return validator matching with public key and number", func(t *testing.T) {
		res, err := client.GetValidatorByNumber(tCtx,
			&pactus.GetValidatorByNumberRequest{Number: val1.Number()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, val1.PublicKey().String(), res.GetValidator().PublicKey)
		assert.Equal(t, val1.Number(), res.GetValidator().GetNumber())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestGetValidatorAddresses(t *testing.T) {
	conn, client := testBlockchainClient(t)

	t.Run("should return list of validator addresses", func(t *testing.T) {
		res, err := client.GetValidatorAddresses(tCtx,
			&pactus.GetValidatorAddressesRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res.GetAddresses()))
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestConsensusInfo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testBlockchainClient(t)

	v1, _ := ts.GenerateTestPrepareVote(100, 2)
	v2, _ := ts.GenerateTestChangeProposerVote(100, 2)
	tConsMocks[1].Active = true
	tConsMocks[1].Height = 100
	tConsMocks[0].AddVote(v1)
	tConsMocks[1].AddVote(v2)

	t.Run("Should return the consensus info", func(t *testing.T) {
		res, err := client.GetConsensusInfo(tCtx, &pactus.GetConsensusInfoRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.False(t, res.Instances[0].Active, true)
		assert.True(t, res.Instances[1].Active, true)
		assert.Equal(t, res.Instances[1].Height, uint32(100))
		assert.Equal(t, res.Instances[0].Votes[0].Type, pactus.VoteType_VOTE_PREPARE)
		assert.Equal(t, res.Instances[1].Votes[0].Type, pactus.VoteType_VOTE_CHANGE_PROPOSER)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
