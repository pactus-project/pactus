package txpool

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

func TestValidity(t *testing.T) {
	logger.InitLogger(logger.DefaultConfig())
	conf := DefaultConfig()
	pool, _ := NewTxPool(conf, nil)
	sb := sandbox.NewMockSandbox()
	pool.SetSandbox(sb)

	stamp := crypto.GenerateTestHash()
	senderAddr, _, _ := crypto.GenerateTestKeyPair()
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()
	bigMemo := strings.Repeat("a", 1025)

	sb.TTLInterval = 2
	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 1000, bigMemo, nil, nil)
	assert.Error(t, pool.(*txPool).validateTx(trx1))

	trx2 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 999, "", nil, nil)
	assert.Error(t, pool.(*txPool).validateTx(trx2))

	trx3 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 10000000, 1000, "", nil, nil)
	assert.Error(t, pool.(*txPool).validateTx(trx3))

	trx4 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 123456789, 123456, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx4))

	trx5 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 123456000, 123456, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx5))

	trx6 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1, 1000, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx6))

	trx7 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 1000, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx7))

	trx8 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000000, 1000, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx8))

	trx9 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 10000000, 10000, "ok", nil, nil)
	assert.NoError(t, pool.(*txPool).validateTx(trx9))

	// Test stamp validity
	sb.AppendStampAndUpdateHeight(101, crypto.GenerateTestHash())
	assert.NoError(t, pool.(*txPool).validateTx(trx9))
	sb.AppendStampAndUpdateHeight(102, crypto.GenerateTestHash())
	assert.Error(t, pool.(*txPool).validateTx(trx9))

}
