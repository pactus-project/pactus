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

	signer1 := bls.GenerateTestSigner()
	addr1 := signer1.Address()
	acc1 := account.NewAccount(addr1, 1)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)

	rcvAddr := crypto.GenerateTestAddress()
	hash1 := hash.GenerateTestHash()
	hash2 := hash.GenerateTestHash()
	hash3 := hash.GenerateTestHash()
	hash8635 := hash.GenerateTestHash()
	hash8640 := hash.GenerateTestHash()
	hash8641 := hash.GenerateTestHash()
	hash8642 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(0, hash.UndefHash)
	tSandbox.AppendNewBlock(1, hash1)
	tSandbox.AppendNewBlock(2, hash2)
	tSandbox.AppendNewBlock(3, hash3)
	tSandbox.AppendNewBlock(8635, hash8635)
	tSandbox.AppendNewBlock(8640, hash8640)
	tSandbox.AppendNewBlock(8641, hash8641)
	tSandbox.AppendNewBlock(8642, hash8642)

	t.Run("Invalid transaction, Should returns error", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Zero(t, tExec.AccumulatedFee())
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Expired stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash1.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "expired-stamp")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Good stamp", func(t *testing.T) {
		trx := tx.NewSendTx(hash3.Stamp(), 1, addr1, rcvAddr, 1000, 1000, "ok")
		signer1.SignMsg(trx)
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase invalid stamp, Should returns error", func(t *testing.T) {
		trx := tx.NewMintbaseTx(hash8641.Stamp(), 1, rcvAddr, 1000, "expired-stamp")
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Mintbase stamp is ok", func(t *testing.T) {
		trx := tx.NewMintbaseTx(hash8642.Stamp(), 1, rcvAddr, 1000, "ok")
		assert.NoError(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Big memo, Should returns error", func(t *testing.T) {
		bigMemo := strings.Repeat("a", 1025)
		trx := tx.NewSendTx(hash8641.Stamp(), 2, addr1, rcvAddr, 1000, 1000, bigMemo)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash2.Stamp(), 2, addr1, rcvAddr, 1000, 1, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash2.Stamp(), 2, addr1, rcvAddr, 1000, 1001, "invalid fee")
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Invalid fee, Should returns error", func(t *testing.T) {
		trx := tx.NewSendTx(hash2.Stamp(), 2, crypto.TreasuryAddress, rcvAddr, 1000, 1001, "invalid fee")
		assert.Error(t, tExec.Execute(trx, tSandbox))
		assert.Error(t, tExec.checkFee(trx, tSandbox))
	})

	t.Run("Sortition tx - Invalid stamp, Should returns error", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(hash8635.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})

	t.Run("Execution failed", func(t *testing.T) {
		proof := sortition.GenerateRandomProof()
		trx := tx.NewSortitionTx(hash8642.Stamp(), 1, addr1, proof)
		signer1.SignMsg(trx)
		assert.Error(t, tExec.Execute(trx, tSandbox))
	})
}

func TestChecker(t *testing.T) {
	tChecker := NewChecker()
	tSandbox := sandbox.MockingSandbox()

	hash1000 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(1000, hash1000)

	t.Run("Accept bond transaction for future blocks", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		acc, signer := account.GenerateTestAccount(1)
		tSandbox.Accounts[acc.Address()] = acc

		tSandbox.InCommittee = true
		trx := tx.NewBondTx(hash1000.Stamp(), acc.Sequence()+1, acc.Address(), pub, 1000, 1000, "")
		signer.SignMsg(trx)
		assert.NoError(t, tChecker.Execute(trx, tSandbox))
	})
}
