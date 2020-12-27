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
		c := makeCommitAndSign(t, b1.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

		assert.NoError(t, st1.ApplyBlock(1, b1, c))
		assert.NoError(t, st2.ApplyBlock(1, b1, c))
	}

	b2 := st2.ProposeBlock()

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
		validSig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})

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

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})
		invAddr, _, _ := crypto.GenerateTestKeyPair()

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: invAddr, Status: 1},
		}, sig)

		assert.Error(t, st1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig4})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 0},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, st1.UpdateLastCommit(c))
	})

	t.Run("Update last commit- Valid signature, should return no error", func(t *testing.T) {

		sig := crypto.Aggregate([]*crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCommit(0, []block.Committer{
			{Address: tValSigner1.Address(), Status: 1},
			{Address: tValSigner2.Address(), Status: 1},
			{Address: tValSigner3.Address(), Status: 1},
			{Address: tValSigner4.Address(), Status: 1},
		}, sig)

		assert.NoError(t, st1.UpdateLastCommit(c))
	})

}
