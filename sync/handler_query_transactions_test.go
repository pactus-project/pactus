package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

func TestParsingQueryTransactionsMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()
	trx4, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in her cache
	tAliceSync.cache.AddTransaction(trx1)
	tAliceSync.cache.AddTransaction(trx2)
	tBobSync.cache.AddTransaction(trx3)
	tBobSync.cache.AddTransaction(trx4)
	pld := payload.NewQueryTransactionsPayload([]hash.Hash{trx2.ID(), trx3.ID(), trx4.ID()})

	t.Run("Alice should not send query transaction message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryTransactions)
		assert.NotNil(t, tAliceState.PendingTx(trx2.ID()))
	})

	t.Run("Bob should not process Alice's message because she is not an active validator", func(t *testing.T) {
		simulatingReceiveingNewMessage(t, tBobSync, pld, tAlicePeerID)
		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeTransactions)
	})

	joinAliceToCommittee(t)

	t.Run("Alice sends query transaction message", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Bob processes Alice's request", func(t *testing.T) {
		simulatingReceiveingNewMessage(t, tBobSync, pld, tAlicePeerID)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeTransactions)
	})

	t.Run("Alice queries for a transaction, but she has it in her cache", func(t *testing.T) {
		tAliceBroadcastCh <- payload.NewQueryTransactionsPayload([]hash.Hash{trx1.ID()})
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryTransactions)
	})
}
