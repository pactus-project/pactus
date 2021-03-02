package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestTransactionLost(t *testing.T) {
	setup(t)

	b1, _ := tState1.ProposeBlock(0)
	assert.NoError(t, tState2.ValidateBlock(*b1))

	b2, _ := tState1.ProposeBlock(0)
	tCommonTxPool.Txs = make([]*tx.Tx, 0)
	assert.Error(t, tState2.ValidateBlock(*b2))
}

func TestCommitValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	val5, siners5 := validator.GenerateTestValidator(4)
	tState1.store.UpdateValidator(val5)
	tState2.store.UpdateValidator(val5)

	b2, _ := tState2.ProposeBlock(0)

	invBlockHash := crypto.GenerateTestHash()
	round := 0
	valSig1 := tValSigner1.SignData(block.CommitSignBytes(b2.Hash(), round))
	valSig2 := tValSigner2.SignData(block.CommitSignBytes(b2.Hash(), round))
	valSig3 := tValSigner3.SignData(block.CommitSignBytes(b2.Hash(), round))
	valSig4 := tValSigner4.SignData(block.CommitSignBytes(b2.Hash(), round))
	invSig1 := tValSigner1.SignData(block.CommitSignBytes(invBlockHash, round))
	invSig2 := tValSigner2.SignData(block.CommitSignBytes(invBlockHash, round))
	invSig3 := tValSigner3.SignData(block.CommitSignBytes(invBlockHash, round))
	invSig5 := siners5.SignData(block.CommitSignBytes(b2.Hash(), round))

	validSig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3})
	invalidSig := crypto.Aggregate([]crypto.Signature{invSig1, invSig2, invSig3})

	t.Run("Invalid blockhahs, should return error", func(t *testing.T) {
		c := block.NewCommit(invBlockHash, 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, validSig)

		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		invSig := tValSigner1.SignData([]byte("abc"))
		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, invSig)

		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid signer, should return error", func(t *testing.T) {
		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 4, Status: 0},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))

		c2 := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 4, Status: 1},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c2))

		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, invSig5})
		c3 := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 0},
			{Number: 4, Status: 1},
		}, sig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c3))
	})

	t.Run("Unexpected signature", func(t *testing.T) {
		sig1 := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, invSig3, valSig4})
		c1 := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, sig1)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c1))

		sig2 := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, invSig5})
		c2 := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 4, Status: 0},
		}, sig2)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c2)) // committee hash is invalid
	})

	t.Run("duplicated or missed number, should return error", func(t *testing.T) {
		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 2, Status: 1},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))

		c = block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("unexpected block hash", func(t *testing.T) {
		c := block.NewCommit(invBlockHash, 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, invalidSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))

	})

	t.Run("Invalid round", func(t *testing.T) {
		c := block.NewCommit(b2.Hash(), 1, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Doesn't have 2/3 majority, should return no error", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2})

		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 0},
			{Number: 3, Status: 0},
		}, sig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid committer, should return no error", func(t *testing.T) {
		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 5, Status: 1},
		}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Valid signature, should return no error", func(t *testing.T) {
		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 0},
		}, validSig)
		assert.NoError(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Update last commit- Invalid signer", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, invSig5})

		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 4, Status: 1},
		}, sig)
		assert.Error(t, tState1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- valid signature, should return no error", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig4})

		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 0},
			{Number: 3, Status: 1},
		}, sig)
		assert.NoError(t, tState1.UpdateLastCommit(c))
		// Commit didn't change
		assert.NotEqual(t, tState1.lastCommit.Hash(), c.Hash())
	})

	t.Run("Update last commit- Valid signature, should return no error", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCommit(b2.Hash(), 0, []block.Committer{
			{Number: 0, Status: 1},
			{Number: 1, Status: 1},
			{Number: 2, Status: 1},
			{Number: 3, Status: 1},
		}, sig)
		assert.NoError(t, tState1.UpdateLastCommit(c))
		// Commit updated
		assert.Equal(t, tState1.lastCommit.Hash(), c.Hash())
	})

}

func TestUpdateBlockTime(t *testing.T) {
	setup(t)

	fmt.Println(tGenTime)
	// Maipulate last block time
	tState1.lastBlockTime = util.Now().Add(-6 * time.Second)
	b, _ := tState1.ProposeBlock(0)
	fmt.Println(b.Header().Time())
	assert.True(t, b.Header().Time().After(tState1.lastBlockTime))
	assert.Zero(t, b.Header().Time().Second()%10)

	tState1.lastBlockTime = util.Now().Add(-16 * time.Second)
	b, _ = tState1.ProposeBlock(0)
	fmt.Println(b.Header().Time())
	assert.True(t, b.Header().Time().After(tState1.lastBlockTime))
	assert.Zero(t, b.Header().Time().Second()%10)
}

func TestBlockValidation(t *testing.T) {
	setup(t)

	moveToNextHeightForAllStates(t)

	assert.False(t, tState1.lastBlockHash.EqualsTo(crypto.UndefHash))

	//
	// Version   			(SanityCheck)
	// UnixTime				(?)
	// LastBlockHash		(OK)
	// StateHash			(OK)
	// TxIDsHash			(?)
	// LastReceiptsHash		(OK)
	// LastCommitHash		(OK)
	// CommitteeHash		(OK)
	// SortitionSeed		(OK)
	// ProposerAddress		(OK)
	//
	invAddr, _, _ := crypto.GenerateTestKeyPair()
	invHash := crypto.GenerateTestHash()
	invCommit := block.GenerateTestCommit(tState1.lastBlockHash)
	trx := tState2.createSubsidyTx(0)
	tCommonTxPool.AppendTx(trx)
	ids := block.NewTxIDs()
	ids.Append(trx.ID())

	b := block.MakeBlock(1, util.Now(), ids, invHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, invHash, tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), invHash, tState1.lastReceiptsHash, tState1.lastCommit, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), invHash, tState1.lastCommit, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, invCommit, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.lastSortitionSeed, invAddr)
	assert.NoError(t, tState1.validateBlock(b))
	c := makeCommitAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, sortition.GenerateRandomSeed(), tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCommitAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.lastSortitionSeed.Generate(tState2.signer), tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCommitAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, tState1.CommitBlock(2, b, c))
}
