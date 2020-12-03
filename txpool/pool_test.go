package txpool

import (
	"testing"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"

	"github.com/zarbchain/zarb-go/logger"
)

var pool *txPool

func init() {
	logger.InitLogger(logger.DefaultConfig())
	conf := TestConfig()
	p, _ := NewTxPool(conf, nil)
	pool = p.(*txPool)
}

func TestAppendAndRemove(t *testing.T) {
	trx, _ := tx.GenerateTestSendTx()
	assert.NoError(t, pool.AppendTx(*trx))
	assert.True(t, pool.HasTx(trx.Hash()))
	assert.False(t, pool.HasTx(crypto.GenerateTestHash()))
	assert.Equal(t, pool.RemoveTx(trx.Hash()), trx)
	assert.False(t, pool.HasTx(trx.Hash()))
}

func TestValidationWhileSyncing(t *testing.T) {
	invalTrx, _ := tx.GenerateTestSendTx()
	assert.NoError(t, pool.appendTx(*invalTrx))
	pool.isSyncing = false
	assert.Error(t, pool.appendTx(*invalTrx))
}
