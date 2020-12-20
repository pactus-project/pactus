package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestTransactionLost(t *testing.T) {
	st1 := setupStatewithFourValidators(t, tValSigner1)
	st2 := setupStatewithFourValidators(t, tValSigner2)

	b1 := st1.ProposeBlock()
	assert.NoError(t, st2.ValidateBlock(b1))

	b2 := st1.ProposeBlock()
	tCommonTxPool.Txs = make([]*tx.Tx, 0)
	assert.Error(t, st2.ValidateBlock(b2))
}

func TestCommitValidation(t *testing.T) {
	st1 := setupStatewithFourValidators(t, tValSigner1)
	st2 := setupStatewithFourValidators(t, tValSigner2)

	b1 := st1.ProposeBlock()
	{
		v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b1.Hash(), tValSigner1.Address())
		v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b1.Hash(), tValSigner2.Address())
		v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, b1.Hash(), tValSigner3.Address())
		tValSigner1.SignMsg(v1)
		tValSigner2.SignMsg(v2)
		tValSigner3.SignMsg(v3)

		validSig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v3.Signature()})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.NoError(t, st1.ApplyBlock(1, b1, *c))
		assert.NoError(t, st2.ApplyBlock(1, b1, *c))
	}

	b2 := st2.ProposeBlock()

	invBlockHash := crypto.GenerateTestHash()
	v1 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, b2.Hash(), tValSigner1.Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, b2.Hash(), tValSigner2.Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, b2.Hash(), tValSigner3.Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, b2.Hash(), tValSigner4.Address())
	invVote1 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, invBlockHash, tValSigner1.Address())
	invVote2 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, invBlockHash, tValSigner1.Address())
	invVote3 := vote.NewVote(vote.VoteTypePrecommit, 2, 0, invBlockHash, tValSigner1.Address())
	tValSigner1.SignMsg(v1)
	tValSigner2.SignMsg(v2)
	tValSigner3.SignMsg(v3)
	tValSigner4.SignMsg(v4)
	tValSigner1.SignMsg(invVote1)
	tValSigner2.SignMsg(invVote2)
	tValSigner3.SignMsg(invVote3)

	validSig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v3.Signature()})
	invalidSig := crypto.Aggregate([]*crypto.Signature{invVote1.Signature(), invVote2.Signature(), invVote3.Signature()})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		invSig := tValSigner1.Sign([]byte("abc"))
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, *invSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Invalid signer, should return error", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: invAddr, Status: 0},
		}, validSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Unexpected signature", func(t *testing.T) {
		validSig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v3.Signature(), v4.Signature()})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Invalid signature status", func(t *testing.T) {
		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, validSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Invalid block hash", func(t *testing.T) {

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, invalidSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Invalid round", func(t *testing.T) {

		c := block.NewCommit(1, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.Error(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Valid signature, should return no error", func(t *testing.T) {

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 0},
		}, validSig)

		assert.NoError(t, st1.ApplyBlock(2, b2, *c))
	})

	t.Run("Update last commit- Invalid signer", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v3.Signature(), v4.Signature()})
		invAddr, _, _ := crypto.GenerateTestKeyPair()

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: invAddr, Status: 1},
		}, sig)

		assert.Error(t, st1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- Invalid status", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v4.Signature()})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 0},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.Error(t, st1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v4.Signature()})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 0},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, st1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- Valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{v1.Signature(), v2.Signature(), v3.Signature(), v4.Signature()})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, st1.UpdateLastCommit(c))
	})

}
