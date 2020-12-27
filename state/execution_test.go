package state

import (
	"testing"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/util"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestProposeBlock(t *testing.T) {
	st1 := setupStatewithOneValidator(t)

	b1, c1 := proposeAndSignBlock(t, st1)
	assert.NoError(t, st1.ApplyBlock(1, b1, c1))

	subsidy := calcBlockSubsidy(st1.LastBlockHeight(), st1.params.SubsidyReductionInterval)
	invSubsidyTx := tx.NewSubsidyTx(st1.LastBlockHash(), 1, tValSigner2.Address(), subsidy, "")
	invSendTx, _ := tx.GenerateTestSendTx()
	invBondTx, _ := tx.GenerateTestBondTx()
	invSortitionTx, _ := tx.GenerateTestSortitionTx()

	pub := tValSigner1.PublicKey()
	trx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "", &pub, nil)
	tValSigner1.SignMsg(trx1)

	trx2 := tx.NewBondTx(b1.Hash(), 2, tValSigner1.Address(), pub, 1, "", &pub, nil)
	tValSigner1.SignMsg(trx2)

	assert.NoError(t, st1.txPool.AppendTx(invSendTx))
	assert.NoError(t, st1.txPool.AppendTx(invBondTx))
	assert.NoError(t, st1.txPool.AppendTx(invSortitionTx))
	assert.NoError(t, st1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, st1.txPool.AppendTx(trx1))
	assert.NoError(t, st1.txPool.AppendTx(trx2))

	b2, c2 := proposeAndSignBlock(t, st1)
	assert.Equal(t, b2.Header().LastBlockHash(), b1.Hash())
	assert.Equal(t, b2.TxIDs().IDs()[1:], []crypto.Hash{trx1.ID(), trx2.ID()})
	assert.NoError(t, st1.ApplyBlock(2, b2, c2))
}

func TestExecuteBlock(t *testing.T) {
	st := setupStatewithOneValidator(t)

	b1, c1 := proposeAndSignBlock(t, st)
	assert.NoError(t, st.ApplyBlock(1, b1, c1))

	invSubsidyTx := st.createSubsidyTx(1001)
	validSubsidyTx := st.createSubsidyTx(1000)
	invSendTx, _ := tx.GenerateTestSendTx()

	pub := tValSigner1.PublicKey()
	validTx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "", &pub, nil)
	tValSigner1.SignMsg(validTx1)

	assert.NoError(t, st.txPool.AppendTx(invSendTx))
	assert.NoError(t, st.txPool.AppendTx(validSubsidyTx))
	assert.NoError(t, st.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, st.txPool.AppendTx(validTx1))

	t.Run("Subsidy tx is invalid", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(invSubsidyTx.ID())
		invBlock := block.MakeBlock(util.Now(), txIDs, st.lastBlockHash, st.validatorSet.CommittersHash(), st.stateHash(), st.lastReceiptsHash, st.lastCommit, st.proposer)
		_, err := st.executeBlock(invBlock)
		assert.Error(t, err)
	})

	t.Run("Has invalid tx", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validSubsidyTx.ID())
		txIDs.Append(invSendTx.ID())
		invBlock := block.MakeBlock(util.Now(), txIDs, st.lastBlockHash, st.validatorSet.CommittersHash(), st.stateHash(), st.lastReceiptsHash, st.lastCommit, st.proposer)
		_, err := st.executeBlock(invBlock)
		assert.Error(t, err)
	})

	t.Run("Subsidy is not first tx", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validTx1.ID())
		txIDs.Append(validSubsidyTx.ID())
		invBlock := block.MakeBlock(util.Now(), txIDs, st.lastBlockHash, st.validatorSet.CommittersHash(), st.stateHash(), st.lastReceiptsHash, st.lastCommit, st.proposer)
		_, err := st.executeBlock(invBlock)
		assert.Error(t, err)
	})

	t.Run("Has no subsidy", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validTx1.ID())
		invBlock := block.MakeBlock(util.Now(), txIDs, st.lastBlockHash, st.validatorSet.CommittersHash(), st.stateHash(), st.lastReceiptsHash, st.lastCommit, st.proposer)
		_, err := st.executeBlock(invBlock)
		assert.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		txIDs := block.NewTxIDs()
		txIDs.Append(validSubsidyTx.ID())
		txIDs.Append(validTx1.ID())
		invBlock := block.MakeBlock(util.Now(), txIDs, st.lastBlockHash, st.validatorSet.CommittersHash(), st.stateHash(), st.lastReceiptsHash, st.lastCommit, st.proposer)
		_, err := st.executeBlock(invBlock)
		assert.NoError(t, err)
	})
}
