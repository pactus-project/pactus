package grpc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callServer(t)
	tx1, _ := tx.GenerateTestSortitionTx()

	tMockState.Store.SaveTransaction(&tx.CommittedTx{Tx: tx1})
	t.Run("Should return transaction, verbosity 0", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: tx1.ID().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.Json)
	})

	t.Run("Should return transaction, verbosity 1", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: tx1.ID().String(), Verbosity: 1})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Json)
	})

	t.Run("Should return nil value because transcation id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: "invalid_id"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transcation doesn't exist", func(t *testing.T) {
		id := crypto.GenerateTestHash()
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: id.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	conn.Close()

}

func TestSendRawTransaction(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should fail, Non decodable hex", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: "None Decodable hex"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, Non decodable data", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString([]byte("None Decodable data"))})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, Non verifable trx", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		_, signer := tx.GenerateTestSendTx()
		trx.SetSignature(signer.SignData(trx.SignBytes()))
		data, _ := trx.Encode()
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := tx.GenerateTestSendTx()
	data, _ := trx.Encode()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail, Not Broadcasted", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.NoError(t, conn.Close())
}
