package execution

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecution(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	tSandbox := sandbox.MockingSandbox()
	tExec := NewExecution(tSandbox)

	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(2100000000000000 - 10000000000)
	tSandbox.UpdateAccount(acc0)

	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	acc1 := account.NewAccount(addr1, 1)
	acc1.AddToBalance(2100000000000000 - 10000000000)
	tSandbox.UpdateAccount(acc1)

	rcvAddr, _, _ := crypto.GenerateTestKeyPair()
	stamp1 := crypto.GenerateTestHash()
	stamp2 := crypto.GenerateTestHash()
	stamp3 := crypto.GenerateTestHash()
	stamp4 := crypto.GenerateTestHash()
	stamp5 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(0, crypto.UndefHash)
	tSandbox.AppendStampAndUpdateHeight(1, stamp1)
	tSandbox.AppendStampAndUpdateHeight(2, stamp2)
	tSandbox.AppendStampAndUpdateHeight(3, stamp3)
	tSandbox.AppendStampAndUpdateHeight(4, stamp4)
	tSandbox.AppendStampAndUpdateHeight(5, stamp5)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		trx.SetPublicKey(nil)
		assert.Error(t, tExec.Execute(trx))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(crypto.UndefHash, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp", &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp1, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp", &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Good stamp", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 1, addr1, rcvAddr, 1000, 1000, "ok", &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.NoError(t, tExec.Execute(trx))
	})

	t.Run("Subsidy invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSubsidyTx(stamp4, 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Subsidy stamp is ok", func(t *testing.T) {
		trx := tx.NewSubsidyTx(stamp5, 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx))
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 1025)
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1000, bigMemo, &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1, "invalid fee", &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1001, "invalid fee", &pub1, nil)
		trx.SetSignature(priv1.Sign(trx.SignBytes()))
		assert.Error(t, tExec.Execute(trx))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, crypto.TreasuryAddress, rcvAddr, 1000, 1001, "invalid fee", nil, nil)
		assert.Error(t, tExec.Execute(trx))
		assert.Error(t, tExec.checkFee(trx))
	})

	assert.Equal(t, tExec.AccumulatedFee(), int64(1000))
	tExec.ResetFee()
	assert.Equal(t, tExec.AccumulatedFee(), int64(0))
}
