package execution

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecution(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	tSandbox := sandbox.MockingSandbox()
	tExec := NewExecution(tSandbox)

	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(21*1e14 - 10000000000)
	tSandbox.UpdateAccount(acc0)

	signer1 := crypto.GenerateTestSigner()
	addr1 := signer1.Address()
	acc1 := account.NewAccount(addr1, 1)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	rcvAddr, _, _ := crypto.GenerateTestKeyPair()
	stamp1 := crypto.GenerateTestHash()
	stamp2 := crypto.GenerateTestHash()
	stamp3 := crypto.GenerateTestHash()
	stamp8635 := crypto.GenerateTestHash()
	stamp8640 := crypto.GenerateTestHash()
	stamp8641 := crypto.GenerateTestHash()
	stamp8642 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(0, crypto.UndefHash)
	tSandbox.AppendStampAndUpdateHeight(1, stamp1)
	tSandbox.AppendStampAndUpdateHeight(2, stamp2)
	tSandbox.AppendStampAndUpdateHeight(3, stamp3)
	tSandbox.AppendStampAndUpdateHeight(8635, stamp8635)
	tSandbox.AppendStampAndUpdateHeight(8640, stamp8640)
	tSandbox.AppendStampAndUpdateHeight(8641, stamp8641)
	tSandbox.AppendStampAndUpdateHeight(8642, stamp8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, tExec.Execute(trx))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(crypto.UndefHash, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp1, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Good stamp", func(t *testing.T) {
		trx := tx.NewSendTx(stamp3, 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, tExec.Execute(trx))
	})

	t.Run("Mintbase invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewMintbaseTx(stamp8641, 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Mintbase stamp is ok", func(t *testing.T) {
		trx := tx.NewMintbaseTx(stamp8642, 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx))
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 1025)
		trx := tx.NewSendTx(stamp8641, 2, addr1, rcvAddr, 1000, 1000, bigMemo)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, crypto.TreasuryAddress, rcvAddr, 1000, 1001, "invalid fee")
		assert.Error(t, tExec.Execute(trx))
		assert.Error(t, tExec.checkFee(trx))
	})

	t.Run("Sortition tx - Invalid stamp, Should returns error", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(stamp8635, 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(stamp8642, 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Claim accumulated fee", func(t *testing.T) {
		tExec.ClaimAccumulatedFee()
		total := int64(0)
		for _, acc := range tSandbox.Accounts {
			total += acc.Balance()
		}
		assert.Equal(t, total, int64(21*1e14))
	})

	assert.Equal(t, tExec.AccumulatedFee(), int64(1000))
	tExec.ResetFee()
	assert.Equal(t, tExec.AccumulatedFee(), int64(0))
}
