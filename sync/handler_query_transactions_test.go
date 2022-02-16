package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingQueryTransactionsMessages(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()

	tSync.cache.AddTransaction(trx1)
	tState.Store.SaveTransaction(trx2)
	pid := util.RandomPeerID()
	pld := payload.NewQueryTransactionsPayload([]hash.Hash{trx1.ID(), trx2.ID(), trx3.ID()})

	t.Run("Not in the committee, should not respond to the query transaction message", func(t *testing.T) {
		assert.Error(t, testReceiveingNewMessage(t, tSync, pld, pid))
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeerToCommittee(t, pub, pid)

	t.Run("In the committee, should respond to the query transaction message", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeTransactions)
		assert.Len(t, msg.Payload.(*payload.TransactionsPayload).Transactions, 2)
	})

	t.Run("In the committee, but doesn't have the transaction", func(t *testing.T) {
		pld := payload.NewQueryTransactionsPayload([]hash.Hash{hash.GenerateTestHash()})
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeTransactions)
	})
}

func TestBroadcastingQueryTransactionsMessages(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestBondTx()
	tSync.cache.AddTransaction(trx1)
	pld := payload.NewQueryTransactionsPayload([]hash.Hash{trx1.ID(), trx2.ID()})

	t.Run("Not in the committee, should not send query transaction message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryTransactions)
	})

	testAddPeerToCommittee(t, tSync.signer.PublicKey(), tSync.SelfID())

	t.Run("In the committee, should send query transaction message", func(t *testing.T) {
		tSync.broadcast(pld)

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryTransactions)
		assert.NotContains(t, msg.Payload.(*payload.QueryTransactionsPayload).IDs, trx1.ID())
	})

	t.Run("Transaction found inside the cache", func(t *testing.T) {
		tSync.cache.AddTransaction(trx2)
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryTransactions)
	})
}
