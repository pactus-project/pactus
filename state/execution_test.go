package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func TestProposeBlock(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.NoError(t, tState2.CommitBlock(1, b1, c1))

	subsidy := tState1.params.BlockReward
	invSubsidyTx := tx.NewMintbaseTx(tState1.lastInfo.BlockHash(), 1, tValSigner2.Address(), subsidy, "")
	invSendTx, _ := tx.GenerateTestSendTx()
	invBondTx, _ := tx.GenerateTestBondTx()
	invSortitionTx, _ := tx.GenerateTestSortitionTx()

	_, pub, _ := crypto.GenerateTestKeyPair()
	trx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "")
	tValSigner1.SignMsg(trx1)

	trx2 := tx.NewBondTx(b1.Hash(), 2, tValSigner1.Address(), pub, 1000, 1000, "")
	tValSigner1.SignMsg(trx2)

	assert.NoError(t, tState1.txPool.AppendTx(invSendTx))
	assert.NoError(t, tState1.txPool.AppendTx(invBondTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSortitionTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(trx1))
	assert.NoError(t, tState1.txPool.AppendTx(trx2))

	b2, c2 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.Equal(t, b2.Header().LastBlockHash(), b1.Hash())
	assert.Equal(t, b2.TxIDs().IDs()[1:], []crypto.Hash{trx1.ID(), trx2.ID()})
	assert.NoError(t, tState1.CommitBlock(2, b2, c2))

	assert.Equal(t, tState1.sortition.TotalStake(), int64(1000))
}

func TestExecuteBlock(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))

	invSubsidyTx := tState1.createSubsidyTx(1001)
	validSubsidyTx := tState1.createSubsidyTx(1000)
	invSendTx, _ := tx.GenerateTestSendTx()

	validTx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "")
	tValSigner1.SignMsg(validTx1)

	assert.NoError(t, tState1.txPool.AppendTx(invSendTx))
	assert.NoError(t, tState1.txPool.AppendTx(validSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, tState1.txPool.AppendTx(validTx1))

	t.Run("Subsidy tx is invalid", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(invSubsidyTx.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.Error(t, err)
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validSubsidyTx.ID())
		txIDs.Append(invSendTx.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.Error(t, err)
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validTx1.ID())
		txIDs.Append(validSubsidyTx.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.Error(t, err)
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validTx1.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.Error(t, err)
	})

	t.Run("Two subsidy transactions", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validSubsidyTx.ID())
		txIDs.Append(validSubsidyTx.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validSubsidyTx.ID())
		txIDs.Append(validTx1.ID())
		invBlock := block.MakeBlock(1, util.Now(), txIDs, tState1.lastInfo.BlockHash(), tState1.stateHash(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState1.signer.Address())
		sb := tState1.makeSandbox()
		_, err := tState1.executeBlock(invBlock, sb)
		assert.NoError(t, err)

		// Check if fee is claimed
		treasury := sb.Account(crypto.TreasuryAddress)
		subsidy := tState1.params.BlockReward
		assert.Equal(t, treasury.Balance(), 21*1e14-(2*subsidy)) // Two blocks has committed yet
	})
}
