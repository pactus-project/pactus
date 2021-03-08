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

func TestCertificateValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	val5, siners5 := validator.GenerateTestValidator(4)
	tState1.store.UpdateValidator(val5)
	tState2.store.UpdateValidator(val5)

	b2, _ := tState2.ProposeBlock(0)

	invBlockHash := crypto.GenerateTestHash()
	round := 0
	valSig1 := tValSigner1.SignData(block.CertificateSignBytes(b2.Hash(), round))
	valSig2 := tValSigner2.SignData(block.CertificateSignBytes(b2.Hash(), round))
	valSig3 := tValSigner3.SignData(block.CertificateSignBytes(b2.Hash(), round))
	valSig4 := tValSigner4.SignData(block.CertificateSignBytes(b2.Hash(), round))
	invSig1 := tValSigner1.SignData(block.CertificateSignBytes(invBlockHash, round))
	invSig2 := tValSigner2.SignData(block.CertificateSignBytes(invBlockHash, round))
	invSig3 := tValSigner3.SignData(block.CertificateSignBytes(invBlockHash, round))
	invSig5 := siners5.SignData(block.CertificateSignBytes(b2.Hash(), round))

	validSig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3})
	invalidSig := crypto.Aggregate([]crypto.Signature{invSig1, invSig2, invSig3})

	t.Run("Invalid blockhahs, should return error", func(t *testing.T) {
		c := block.NewCertificate(invBlockHash, 0, []int{0, 1, 2, 3}, []int{3}, validSig)

		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		invSig := tValSigner1.SignData([]byte("abc"))
		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{3}, invSig)

		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 4}, []int{4}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))

		c2 := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 4}, []int{}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c2))

		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, invSig5})
		c3 := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 5}, []int{2}, sig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c3))
	})

	t.Run("Unexpected signature", func(t *testing.T) {
		sig1 := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, invSig3, valSig4})
		c1 := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{3}, sig1)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c1))

		sig2 := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, invSig5})
		c2 := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 4}, []int{4}, sig2)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c2)) // committee hash is invalid
	})

	t.Run("duplicated or missed number, should return error", func(t *testing.T) {
		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 2}, []int{}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))

		c = block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2}, []int{}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("unexpected block hash", func(t *testing.T) {
		c := block.NewCertificate(invBlockHash, 0, []int{0, 1, 2, 3}, []int{3}, invalidSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("invalid signature", func(t *testing.T) {
		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{3}, invalidSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid round", func(t *testing.T) {
		valSig1 := tValSigner1.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		valSig2 := tValSigner2.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		valSig3 := tValSigner3.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		validSig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3})

		c := block.NewCertificate(b2.Hash(), 1, []int{0, 1, 2, 3}, []int{3}, validSig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2})

		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{2, 3}, sig)
		assert.Error(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Valid signature, should return no error", func(t *testing.T) {
		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{3}, validSig)
		assert.NoError(t, tState1.CommitBlock(2, *b2, *c))
	})

	t.Run("Invalid round", func(t *testing.T) {
		valSig1 := tValSigner1.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		valSig2 := tValSigner2.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		valSig3 := tValSigner3.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		valSig4 := tValSigner4.SignData(block.CertificateSignBytes(b2.Hash(), round+1))
		validSig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCertificate(b2.Hash(), 1, []int{0, 1, 2, 3}, []int{}, validSig)
		assert.Error(t, tState1.UpdateLastCertificate(c))
	})

	t.Run("Update last commit- Invalid committers", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, invSig5})

		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 4}, []int{}, sig)
		assert.Error(t, tState1.UpdateLastCertificate(c))
	})

	t.Run("Update last commit- Invalid signer", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, invSig5})

		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{}, sig)
		assert.Error(t, tState1.UpdateLastCertificate(c))
	})

	t.Run("Update last commit- valid signature, should not update the certificate", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig4})

		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{2}, sig)
		assert.NoError(t, tState1.UpdateLastCertificate(c))
		// Certificate didn't change
		assert.NotEqual(t, tState1.lastCertificate.Hash(), c.Hash())
	})

	t.Run("Update last commit- Valid signature, should update the certificate", func(t *testing.T) {
		sig := crypto.Aggregate([]crypto.Signature{valSig1, valSig2, valSig3, valSig4})

		c := block.NewCertificate(b2.Hash(), 0, []int{0, 1, 2, 3}, []int{}, sig)
		assert.NoError(t, tState1.UpdateLastCertificate(c))
		// Certificate updated
		assert.Equal(t, tState1.lastCertificate.Hash(), c.Hash())
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
	// LastCertificateHash		(OK)
	// CommitteeHash		(OK)
	// SortitionSeed		(OK)
	// ProposerAddress		(OK)
	//
	invAddr, _, _ := crypto.GenerateTestKeyPair()
	invHash := crypto.GenerateTestHash()
	invCert := block.GenerateTestCertificate(tState1.lastBlockHash)
	trx := tState2.createSubsidyTx(0)
	assert.NoError(t, tCommonTxPool.AppendTx(trx))
	ids := block.NewTxIDs()
	ids.Append(trx.ID())

	b := block.MakeBlock(1, util.Now(), ids, invHash, tState1.committee.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCertificate, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, invHash, tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCertificate, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), invHash, tState1.lastReceiptsHash, tState1.lastCertificate, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), tState1.stateHash(), invHash, tState1.lastCertificate, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, invCert, tState1.lastSortitionSeed, tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCertificate, tState1.lastSortitionSeed, invAddr)
	assert.NoError(t, tState1.validateBlock(b))
	c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCertificate, sortition.GenerateRandomSeed(), tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	b = block.MakeBlock(1, util.Now(), ids, tState1.lastBlockHash, tState1.committee.CommitteeHash(), tState1.stateHash(), tState1.lastReceiptsHash, tState1.lastCertificate, tState1.lastSortitionSeed.Generate(tState2.signer), tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, tState1.CommitBlock(2, b, c))
}
