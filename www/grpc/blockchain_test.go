package grpc

import (
	"context"
	"encoding/hex"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	height := uint32(100)
	blk := td.mockState.TestStore.AddTestBlock(height)
	data, _ := blk.Bytes()

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{
				Height: height + 1, Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_DATA,
			})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return an existing block data (verbosity: 0)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_DATA})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, height, res.Height)
		assert.Equal(t, blk.Hash().String(), res.Hash)
		assert.Equal(t, hex.EncodeToString(data), res.Data)
		assert.Empty(t, res.Header)
		assert.Empty(t, res.Txs)
	})

	t.Run("Should return object with (verbosity: 1)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_INFO})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, height, res.Height)
		assert.Equal(t, blk.Hash().String(), res.Hash)
		assert.Empty(t, res.Data)
		assert.NotEmpty(t, res.Header)
		assert.Equal(t, blk.PrevCertificate().Committers(), res.PrevCert.Committers)
		assert.Equal(t, blk.PrevCertificate().Absentees(), res.PrevCert.Absentees)
		for i, trx := range res.Txs {
			blockTrx := blk.Transactions()[i]
			blk, err := blockTrx.Bytes()
			assert.NoError(t, err)
			assert.Equal(t, trx.Id, blockTrx.ID().String())
			assert.Equal(t, trx.Data, hex.EncodeToString(blk))
			assert.Zero(t, trx.LockTime)
			assert.Empty(t, trx.Signature)
			assert.Empty(t, trx.PublicKey)
		}
	})

	t.Run("Should return object with (verbosity: 2)", func(t *testing.T) {
		res, err := client.GetBlock(context.Background(),
			&pactus.GetBlockRequest{Height: height, Verbosity: pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, height, res.Height)
		assert.Equal(t, blk.Hash().String(), res.Hash)
		assert.Empty(t, res.Data)
		assert.NotEmpty(t, res.Header)
		assert.NotEmpty(t, res.Txs)
		for i, trx := range res.Txs {
			blockTrx := blk.Transactions()[i]

			assert.Equal(t, trx.Id, blockTrx.ID().String())
			assert.Empty(t, trx.Data)
			assert.Equal(t, trx.LockTime, blockTrx.LockTime())
			if blockTrx.IsSubsidyTx() {
				assert.Empty(t, trx.Signature)
				assert.Empty(t, trx.PublicKey)
			} else {
				assert.Equal(t, trx.Signature, blockTrx.Signature().String())
				assert.Equal(t, trx.PublicKey, blockTrx.PublicKey().String())
			}
		}
	})
}

func TestGetBlockHash(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	height := td.RandHeight()
	blk := td.mockState.TestStore.AddTestBlock(height)

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(context.Background(),
			&pactus.GetBlockHashRequest{Height: 0})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHash(context.Background(),
			&pactus.GetBlockHashRequest{Height: height})

		assert.NoError(t, err)
		assert.Equal(t, blk.Hash().String(), res.Hash)
	})
}

func TestGetBlockHeight(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	height := td.RandHeight()
	blk := td.mockState.TestStore.AddTestBlock(height)

	t.Run("Should return error for invalid hash", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: ""})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return error for non existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: td.RandHash().String()})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return height of existing block", func(t *testing.T) {
		res, err := client.GetBlockHeight(context.Background(),
			&pactus.GetBlockHeightRequest{Hash: blk.Hash().String()})

		assert.NoError(t, err)
		assert.Equal(t, height, res.Height)
	})
}

func TestGetBlockchainInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	t.Run("Should return the last block height", func(t *testing.T) {
		res, err := client.GetBlockchainInfo(context.Background(),
			&pactus.GetBlockchainInfoRequest{})

		assert.NoError(t, err)
		assert.Equal(t, td.mockState.TestStore.LastHeight, res.LastBlockHeight)
		assert.NotEmpty(t, res.LastBlockHash)
		assert.Zero(t, res.PruningHeight)
		assert.False(t, res.IsPruned)
	})
}

func TestGetCommitteeInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	t.Run("Should return committee info", func(t *testing.T) {
		res, err := client.GetCommitteeInfo(context.Background(),
			&pactus.GetCommitteeInfoRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, res.CommitteePower, int64(0))
		assert.NotNil(t, res.Validators)
		assert.NotNil(t, res.ProtocolVersions)
	})
}

func TestGetAccount(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	addr, acc := td.mockState.TestStore.AddTestAccount()

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
		assert.Equal(t, acc.Balance().ToNanoPAC(), res.Account.Balance)
		assert.Equal(t, acc.Number(), res.Account.Number)
	})
}

func TestGetValidator(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

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
}

func TestGetValidatorByNumber(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

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
}

func TestGetValidatorAddresses(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	t.Run("should return list of validator addresses", func(t *testing.T) {
		td.mockState.TestStore.AddTestValidator()
		td.mockState.TestStore.AddTestValidator()

		res, err := client.GetValidatorAddresses(context.Background(),
			&pactus.GetValidatorAddressesRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res.GetAddresses()))
	})
}

func TestGetPublicKey(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

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
		assert.Equal(t, val.PublicKey().String(), res.PublicKey)
	})
}

func TestConsensusInfo(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	consHeight := td.RandHeight()
	consRound := td.RandRound()
	vote1, _ := td.GenerateTestPrepareVote(consHeight, consRound)
	vote2, _ := td.GenerateTestPrecommitVote(consHeight, consRound)
	prop := td.GenerateTestProposal(consHeight, consRound)

	td.consMocks[0].Active = true
	td.consMocks[0].Height = consHeight
	td.consMocks[0].Round = consRound
	td.consMocks[0].AddVote(vote1)
	td.consMocks[0].AddVote(vote2)
	td.consMocks[0].SetProposal(prop)

	td.consMocks[1].Active = false
	td.consMocks[1].Height = consHeight
	td.consMocks[1].Round = consRound

	t.Run("Should return the consensus info", func(t *testing.T) {
		res, err := client.GetConsensusInfo(context.Background(), &pactus.GetConsensusInfoRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.True(t, res.Instances[0].Active)
		assert.Equal(t, consHeight, res.Instances[0].Height)
		assert.Equal(t, int32(consRound), res.Instances[0].Round)
		assert.Len(t, res.Instances[0].Votes, 2)
		assert.Equal(t, pactus.VoteType_VOTE_TYPE_PREPARE, res.Instances[0].Votes[0].Type)
		assert.Equal(t, pactus.VoteType_VOTE_TYPE_PRECOMMIT, res.Instances[0].Votes[1].Type)

		assert.False(t, res.Instances[1].Active)
		assert.Equal(t, consHeight, res.Instances[1].Height)
		assert.Equal(t, int32(consRound), res.Instances[1].Round)

		assert.Equal(t, consHeight, res.Proposal.Height)
		assert.Equal(t, int32(consRound), res.Proposal.Round)
		assert.Equal(t, prop.Signature().String(), res.Proposal.Signature)
	})
}

func TestGetTxPoolContent(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	_ = td.mockState.AddPendingTx(td.GenerateTestBondTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestBondTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestTransferTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestUnbondTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestSortitionTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestSortitionTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestTransferTx())
	_ = td.mockState.AddPendingTx(td.GenerateTestWithdrawTx())

	t.Run("Should return all transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_UNSPECIFIED,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		assert.Equal(t, 8, len(resp.Txs))
	})

	t.Run("Should return all Bond transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_BOND,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, len(resp.Txs), 0)

		for _, tx := range resp.Txs {
			assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_BOND, tx.PayloadType)
		}
	})

	t.Run("Should return all Transfer transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, len(resp.Txs), 0)

		for _, tx := range resp.Txs {
			assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_TRANSFER, tx.PayloadType)
		}
	})

	t.Run("Should return all Unbond transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_UNBOND,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, len(resp.Txs), 0)

		for _, tx := range resp.Txs {
			assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_UNBOND, tx.PayloadType)
		}
	})

	t.Run("Should return all Sortition transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_SORTITION,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, len(resp.Txs), 0)

		for _, tx := range resp.Txs {
			assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_SORTITION, tx.PayloadType)
		}
	})

	t.Run("Should return all Withdraw transactions", func(t *testing.T) {
		in := &pactus.GetTxPoolContentRequest{
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_WITHDRAW,
		}
		resp, err := client.GetTxPoolContent(context.Background(), in)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, len(resp.Txs), 0)

		for _, tx := range resp.Txs {
			assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_WITHDRAW, tx.PayloadType)
		}
	})
}
