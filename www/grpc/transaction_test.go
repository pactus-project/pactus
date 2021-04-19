package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callServer(t)
	tx1, _ := tx.GenerateTestSortitionTx()

	tMockState.Store.SaveTransaction(tx1)
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
