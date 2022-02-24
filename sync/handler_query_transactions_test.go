package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
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
	msg := message.NewQueryTransactionsMessage([]hash.Hash{trx1.ID(), trx2.ID(), trx3.ID()})

	t.Run("Not in the committee, should not respond to the query transaction message", func(t *testing.T) {
		assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))
	})

	testAddPeerToCommittee(t, pid, nil)

	t.Run("In the committee, should respond to the query transaction message", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeTransactions)
		assert.Len(t, bnd.Message.(*message.TransactionsMessage).Transactions, 2)
	})

	t.Run("In the committee, but doesn't have the transaction", func(t *testing.T) {
		msg := message.NewQueryTransactionsMessage([]hash.Hash{hash.GenerateTestHash()})
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeTransactions)
	})
}

func TestBroadcastingQueryTransactionsMessages(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestBondTx()
	tSync.cache.AddTransaction(trx1)
	msg := message.NewQueryTransactionsMessage([]hash.Hash{trx1.ID(), trx2.ID()})

	t.Run("Not in the committee, should not send query transaction message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryTransactions)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signer.PublicKey())

	t.Run("In the committee, should send query transaction message", func(t *testing.T) {
		tSync.broadcast(msg)

		bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryTransactions)
		assert.NotContains(t, bnd.Message.(*message.QueryTransactionsMessage).IDs, trx1.ID())
	})

	t.Run("Transaction found inside the cache", func(t *testing.T) {
		tSync.cache.AddTransaction(trx2)
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeQueryTransactions)
	})
}
