package txpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

var tPool *txPool
var tSandbox *sandbox.MockSandbox
var tCh chan message.Message
var tTestTx *tx.Tx

func setup(t *testing.T) {
	logger.InitLogger(logger.TestConfig())
	tCh = make(chan message.Message, 10)
	tSandbox = sandbox.MockingSandbox()
	p, err := NewTxPool(TestConfig(), tCh)
	assert.NoError(t, err)
	p.SetNewSandboxAndRecheck(tSandbox)
	tPool = p.(*txPool)

	hash88 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	tTestTx = tx.NewMintbaseTx(hash88.Stamp(), 89, crypto.GenerateTestAddress(), 25000000, "subsidy-tx")
}

func shouldPublishTransaction(t *testing.T, id tx.ID) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-tCh:
			logger.Info("shouldPublishTransaction", "msg", msg)

			if msg.Type() == message.MessageTypeTransactions {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, m.Transactions[0].ID(), id)
				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	setup(t)

	assert.NoError(t, tPool.AppendTx(tTestTx))
	// Appending the same transaction again, should not return any error
	assert.NoError(t, tPool.AppendTx(tTestTx))
	tPool.RemoveTx(tTestTx.ID())
	assert.False(t, tPool.HasTx(tTestTx.ID()), "Transaction should be removed")
}

func TestAppendInvalidTransaction(t *testing.T) {
	setup(t)

	invalidTx, _ := tx.GenerateTestSendTx()
	assert.Error(t, tPool.AppendTx(invalidTx))
}

func TestPending(t *testing.T) {
	setup(t)

	go func(ch chan message.Message) {
		for {
			msg := <-ch
			fmt.Printf("Received a message: %v\n", msg.Fingerprint())
			m := msg.(*message.QueryTransactionsMessage)
			if m.IDs[0].EqualsTo(tTestTx.ID()) {
				assert.NoError(t, tPool.AppendTx(tTestTx))
			}
		}
	}(tCh)

	assert.Nil(t, tPool.PendingTx(tTestTx.ID()))
	assert.NotNil(t, tPool.QueryTx(tTestTx.ID()))
	assert.True(t, tPool.HasTx(tTestTx.ID()))

	invID := hash.GenerateTestHash()
	assert.Nil(t, tPool.PendingTx(invID))
}

// TestFullPool tests if the pool prunes the old transactions when it is full
func TestFullPool(t *testing.T) {
	setup(t)

	hash10000 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(10000, hash10000)
	trxs := make([]*tx.Tx, tPool.config.sendPoolSize()+1)

	signer := bls.GenerateTestSigner()
	acc1 := account.NewAccount(signer.Address(), 0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	// Make sure the pool is empty
	assert.Equal(t, tPool.Size(), 0)

	for i := 0; i < len(trxs); i++ {
		trx := tx.NewSendTx(hash10000.Stamp(), tSandbox.AccSeq(signer.Address())+1, signer.Address(), crypto.GenerateTestAddress(), 1000, 1000, "ok")
		signer.SignMsg(trx)
		assert.NoError(t, tPool.AppendTx(trx))
		trxs[i] = trx
	}

	assert.Nil(t, tPool.QueryTx(trxs[0].ID()))
	assert.NotNil(t, tPool.QueryTx(trxs[1].ID()))
	assert.Equal(t, tPool.Size(), tPool.config.sendPoolSize())
}

func TestEmptyPool(t *testing.T) {
	setup(t)

	assert.Empty(t, tPool.PrepareBlockTransactions(), "pool should be empty")
}

func TestPrepareBlockTransactions(t *testing.T) {
	setup(t)

	hash1000000 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(1000000, hash1000000)

	acc1Signer := bls.GenerateTestSigner()
	acc1 := account.NewAccount(acc1Signer.Address(), 0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	val1Signer := bls.GenerateTestSigner()
	val1Pub := val1Signer.PublicKey().(*bls.PublicKey)
	val1 := validator.NewValidator(val1Pub, 0)
	val1.AddToStake(10000000000)
	tSandbox.UpdateValidator(val1)

	val2Signer := bls.GenerateTestSigner()
	val2Pub := val2Signer.PublicKey().(*bls.PublicKey)
	val2 := validator.NewValidator(val2Pub, 0)
	val2.AddToStake(10000000000)
	val2.UpdateUnbondingHeight(1)
	tSandbox.UpdateValidator(val2)

	val3Signer := bls.GenerateTestSigner()
	val3Pub := val3Signer.PublicKey().(*bls.PublicKey)
	val3 := validator.NewValidator(val3Pub, 0)
	val3.AddToStake(10000000000)
	tSandbox.UpdateValidator(val3)

	sendTx := tx.NewSendTx(hash1000000.Stamp(), tSandbox.AccSeq(acc1.Address())+1, acc1.Address(), crypto.GenerateTestAddress(), 1000, 1000, "send-tx")
	acc1Signer.SignMsg(sendTx)

	pub, _ := bls.GenerateTestKeyPair()
	bondTx := tx.NewBondTx(hash1000000.Stamp(), tSandbox.AccSeq(acc1.Address())+2, acc1.Address(), pub, 1000, 1000, "bond-tx")
	acc1Signer.SignMsg(bondTx)

	unbondTx := tx.NewUnbondTx(hash1000000.Stamp(), tSandbox.ValSeq(val1.Address())+1, val1.Address(), "unbond-tx")
	val1Signer.SignMsg(unbondTx)

	withdrawTx := tx.NewWithdrawTx(hash1000000.Stamp(), tSandbox.ValSeq(val2.Address())+1, val2.Address(), crypto.GenerateTestAddress(), 1000, 1000, "withdraw-tx")
	val2Signer.SignMsg(withdrawTx)

	tSandbox.AcceptSortition = true
	sortitionTx := tx.NewSortitionTx(hash1000000.Stamp(), tSandbox.ValSeq(val3.Address())+1, val3.Address(), sortition.GenerateRandomProof())
	val3Signer.SignMsg(sortitionTx)

	assert.NoError(t, tPool.AppendTx(sendTx))
	assert.NoError(t, tPool.AppendTx(withdrawTx))
	assert.NoError(t, tPool.AppendTx(unbondTx))
	assert.NoError(t, tPool.AppendTx(bondTx))
	assert.NoError(t, tPool.AppendTx(sortitionTx))

	trxs := tPool.PrepareBlockTransactions()
	assert.Len(t, trxs, 5)
	assert.Equal(t, trxs[0].ID(), sortitionTx.ID())
	assert.Equal(t, trxs[1].ID(), bondTx.ID())
	assert.Equal(t, trxs[2].ID(), unbondTx.ID())
	assert.Equal(t, trxs[3].ID(), withdrawTx.ID())
	assert.Equal(t, trxs[4].ID(), sendTx.ID())
}

func TestAppendAndBroadcast(t *testing.T) {
	setup(t)

	assert.NoError(t, tPool.AppendTxAndBroadcast(tTestTx))
	shouldPublishTransaction(t, tTestTx.ID())

	invTrx, _ := tx.GenerateTestBondTx()
	assert.Error(t, tPool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	setup(t)

	hash88 := hash.GenerateTestHash()
	hash89 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	proposer1 := crypto.GenerateTestAddress()
	proposer2 := crypto.GenerateTestAddress()
	trx1 := tx.NewMintbaseTx(hash88.Stamp(), 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewMintbaseTx(hash88.Stamp(), 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewMintbaseTx(hash88.Stamp(), 89, proposer2, 25000000, "subsidy-tx-2")

	assert.Error(t, tPool.AppendTx(trx1), "Expired subsidy transaction")
	assert.NoError(t, tPool.AppendTx(trx2))
	assert.NoError(t, tPool.AppendTx(trx3))

	tSandbox.AppendNewBlock(89, hash89)

	tPool.SetNewSandboxAndRecheck(sandbox.MockingSandbox())
	assert.Zero(t, tPool.Size())
}

func TestSortingSortitions(t *testing.T) {
	setup(t)

	hash1000000 := hash.GenerateTestHash()
	hash1000001 := hash.GenerateTestHash()
	hash1000002 := hash.GenerateTestHash()
	hash1000003 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(1000000, hash1000000)
	tSandbox.AppendNewBlock(1000001, hash1000001)
	tSandbox.AppendNewBlock(1000002, hash1000002)
	tSandbox.AppendNewBlock(1000003, hash1000003)

	val1Signer := bls.GenerateTestSigner()
	val1Pub := val1Signer.PublicKey().(*bls.PublicKey)
	val1 := validator.NewValidator(val1Pub, 0)
	tSandbox.UpdateValidator(val1)

	tSandbox.AcceptSortition = true
	sortitionTx1 := tx.NewSortitionTx(hash1000000.Stamp(), tSandbox.ValSeq(val1.Address())+1, val1.Address(), sortition.GenerateRandomProof())
	val1Signer.SignMsg(sortitionTx1)

	sortitionTx2 := tx.NewSortitionTx(hash1000001.Stamp(), tSandbox.ValSeq(val1.Address())+1, val1.Address(), sortition.GenerateRandomProof())
	val1Signer.SignMsg(sortitionTx2)

	sortitionTx3 := tx.NewSortitionTx(hash1000002.Stamp(), tSandbox.ValSeq(val1.Address())+1, val1.Address(), sortition.GenerateRandomProof())
	val1Signer.SignMsg(sortitionTx3)

	sortitionTx4 := tx.NewSortitionTx(hash1000003.Stamp(), tSandbox.ValSeq(val1.Address())+1, val1.Address(), sortition.GenerateRandomProof())
	val1Signer.SignMsg(sortitionTx4)

	assert.NoError(t, tPool.AppendTx(sortitionTx2))
	assert.NoError(t, tPool.AppendTx(sortitionTx4))
	assert.NoError(t, tPool.AppendTx(sortitionTx1))
	assert.NoError(t, tPool.AppendTx(sortitionTx3))

	trxs := tPool.PrepareBlockTransactions()
	assert.Equal(t, trxs[0].ID(), sortitionTx4.ID())
	assert.Equal(t, trxs[1].ID(), sortitionTx3.ID())
	assert.Equal(t, trxs[2].ID(), sortitionTx2.ID())
	assert.Equal(t, trxs[3].ID(), sortitionTx1.ID())

}
