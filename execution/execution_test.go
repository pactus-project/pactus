package execution

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecution(t *testing.T) {
	logger.InitLogger(logger.TestConfig())

	tSandbox := sandbox.MockingSandbox()
	tExec := NewExecution()

	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(21*1e14 - 10000000000)
	tSandbox.UpdateAccount(acc0)

	signer1 := crypto.GenerateTestSigner()
	addr1 := signer1.Address()
	acc1 := account.NewAccount(addr1, 1)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	rcvAddr, _, _ := bls.GenerateTestKeyPair()
	stamp1 := hash.GenerateTestHash()
	stamp2 := hash.GenerateTestHash()
	stamp3 := hash.GenerateTestHash()
	stamp8635 := hash.GenerateTestHash()
	stamp8640 := hash.GenerateTestHash()
	stamp8641 := hash.GenerateTestHash()
	stamp8642 := hash.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(0, hash.UndefHash)
	tSandbox.AppendStampAndUpdateHeight(1, stamp1)
	tSandbox.AppendStampAndUpdateHeight(2, stamp2)
	tSandbox.AppendStampAndUpdateHeight(3, stamp3)
	tSandbox.AppendStampAndUpdateHeight(8635, stamp8635)
	tSandbox.AppendStampAndUpdateHeight(8640, stamp8640)
	tSandbox.AppendStampAndUpdateHeight(8641, stamp8641)
	tSandbox.AppendStampAndUpdateHeight(8642, stamp8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash.UndefHash, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp1, 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Good stamp", func(t *testing.T) {
		trx := tx.NewSendTx(stamp3, 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewMintbaseTx(stamp8641, 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase stamp is ok", func(t *testing.T) {
		trx := tx.NewMintbaseTx(stamp8642, 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 1025)
		trx := tx.NewSendTx(stamp8641, 2, addr1, rcvAddr, 1000, 1000, bigMemo)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(stamp2, 2, crypto.TreasuryAddress, rcvAddr, 1000, 1001, "invalid fee")
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Error(t, tExec.checkFee(trx, tSandbox))
	})

	t.Run("Sortition tx - Invalid stamp, Should returns error", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(stamp8635, 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(stamp8642, 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})
}

func TestChecker(t *testing.T) {
	tChecker := NewChecker()
	tSandbox := sandbox.MockingSandbox()

	stamp1000 := hash.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(1000, stamp1000)

	t.Run("Accept bond transaction for future blocks", func(t *testing.T) {
		acc, signer := account.GenerateTestAccount(1)
		tSandbox.Accounts[acc.Address()] = acc

		tSandbox.InCommittee = true
		trx := tx.NewBondTx(stamp1000, acc.Sequence()+1, signer.Address(), signer.PublicKey(), 1000, 1000, "")
		signer.SignMsg(trx)
		assert.NoError(t, tChecker.Execute(trx, tSandbox))
	})
}
