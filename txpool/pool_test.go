package txpool

import (
	"testing"

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
func TestValidationWhileSyncing(t *testing.T) {
	invalTrx, _ := tx.GenerateTestSendTx()
	assert.NoError(t, pool.appendTx(*invalTrx))
	pool.isSyncing = false
	assert.Error(t, pool.appendTx(*invalTrx))
}
