package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
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

	b1, c1 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	applyBlockAndCommitForAllStates(t, b1, c1)

	b2, _ := tState2.ProposeBlock(0)

	invBlockHash := crypto.GenerateTestHash()
	round := 0
	valSig1 := tValSigner1.Sign(vote.CommitSignBytes(b2.Hash(), round))
	valSig2 := tValSigner2.Sign(vote.CommitSignBytes(b2.Hash(), round))
	valSig3 := tValSigner3.Sign(vote.CommitSignBytes(b2.Hash(), round))
	valSig4 := tValSigner4.Sign(vote.CommitSignBytes(b2.Hash(), round))
	invSig1 := tValSigner1.Sign(vote.CommitSignBytes(invBlockHash, round))
	invSig2 := tValSigner2.Sign(vote.CommitSignBytes(invBlockHash, round))
	invSig3 := tValSigner3.Sign(vote.CommitSignBytes(invBlockHash, round))

	validSig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3})
	invalidSig := crypto.Aggregate([]*crypto.Signature{invSig1, invSig2, invSig3})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		invSig := tValSigner1.Sign([]byte("abc"))
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, *invSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Invalid signer, should return error", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: invAddr, Status: 0},
		}, validSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Unexpected signature", func(t *testing.T) {
		validSig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Invalid signature status", func(t *testing.T) {
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, validSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Invalid block hash", func(t *testing.T) {

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, invalidSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Invalid round", func(t *testing.T) {

		c := block.NewCommit(1, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.Error(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Valid signature, should return no error", func(t *testing.T) {

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.NoError(t, tState1.ApplyBlock(2, *b2, *c))
	})

	t.Run("Update last commit- Invalid signer", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})
		invAddr, _, _ := crypto.GenerateTestKeyPair()

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: invAddr, Status: 1},
		}, sig)

		assert.Error(t, tState1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig4})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 0},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, tState1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- Valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, tState1.UpdateLastCommit(c))
	})

}

func TestUpdateBlockTime(t *testing.T) {
	setup(t)

	// Maipulate last block time
	tState1.lastBlockTime = util.Now().Add(-6 * time.Second)
	b, _ := tState1.ProposeBlock(0)
	fmt.Println(b.Header().Time())
	assert.True(t, b.Header().Time().After(tState1.lastBlockTime))
	assert.Zero(t, b.Header().Time().Second()%10)
}

func TestBlockValidation(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	applyBlockAndCommitForAllStates(t, b1, c1)
	assert.False(t, tState1.lastBlockHash.EqualsTo(crypto.UndefHash))

	//
	// Version   			(SanityCheck)
	// UnixTime				(?)
	// LastBlockHash		(OK)
	// StateHash			(OK)
	// TxIDsHash			(?)
	// LastReceiptsHash		(OK)
	// LastCommitHash		(OK)
	// CommittersHash		(OK)
	// ProposerAddress		(OK) -> Check in ApplyBlock
	//
	invAdd, _, _ := crypto.GenerateTestKeyPair()
	invHash := crypto.GenerateTestHash()
	invCommit := block.GenerateTestCommit(tState1.lastBlockHash)
	trx := tState1.createSubsidyTx(0)
	ids := block.NewTxIDs()
	ids.Append(trx.ID())

	b := block.MakeBlock(util.Now(), ids, invHash, tState1.validatorSet.CommittersHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.proposer)
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, invHash, tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.proposer)
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommittersHash(), invHash, tState1.lastReceiptsHash, tState1.lastCommit, tState1.proposer)
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommittersHash(), tState1.stateHash(), invHash, tState1.lastCommit, tState1.proposer)
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommittersHash(), tState1.stateHash(), tState1.lastReceiptsHash, invCommit, tState1.proposer)
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommittersHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, invAdd)
	assert.NoError(t, tState1.validateBlock(b))
	c := makeCommitAndSign(t, b.Hash(), 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.ApplyBlock(2, b, c))

	b = block.MakeBlock(util.Now(), ids, tState1.lastBlockHash, tState1.validatorSet.CommittersHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCommit, tState1.proposer)
	assert.NoError(t, tState1.validateBlock(b))
}
