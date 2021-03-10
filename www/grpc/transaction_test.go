package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callServer(t)
	tx1, _ := tx.GenerateTestSortitionTx()

	tMockState.Store.Transactions[tx1.ID()] = &tx.CommittedTx{Tx: tx1}
	t.Run("Should return transaction for value ", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: tx1.ID().Fingerprint()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	conn.Close()

	t.Run("Should return nil value because transcation not created", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: tx1.ID().Fingerprint()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	conn.Close()

}
