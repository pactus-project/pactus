package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func aggregate(sigs []crypto.Signature) *bls.Signature {
	blsSigs := make([]*bls.Signature, len(sigs))
	for i, s := range sigs {
		blsSigs[i] = s.(*bls.Signature)
	}
	return bls.Aggregate(blsSigs)
}

func TestCertificateValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	val5, signer5 := validator.GenerateTestValidator(4)
	tState1.store.UpdateValidator(val5)
	tState2.store.UpdateValidator(val5)

	nextBlock, _ := tState2.ProposeBlock(0)
	nextBlockHash := nextBlock.Hash()

	t.Run("SanityCheck fails, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes).(*bls.Signature)
		sig2 := tValSigner2.SignData(signBytes).(*bls.Signature)
		sig4 := tValSigner4.SignData(signBytes).(*bls.Signature)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(hash.UndefHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		aggSig := signer5.SignData(signBytes).(*bls.Signature)
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid round, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 1)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		committers = append(committers, 666)
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid block hash, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		invBlockHash := hash.GenerateTestHash()
		signBytes := block.CertificateSignBytes(invBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(invBlockHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid committers, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		committers[0] = val5.Number()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := signer5.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2, 3}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{2}, aggSig)

		assert.NoError(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Update last certificate, Invalid committers", func(t *testing.T) {
		committers := tState2.committee.Committers()
		committers = append(committers, val5.Number())
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig3 := tValSigner3.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		sig5 := signer5.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4, sig5})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{}, aggSig)

		assert.Error(t, tState1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Invalid block hash", func(t *testing.T) {
		committers := tState2.committee.Committers()
		invBlockHash := hash.GenerateTestHash()
		signBytes := block.CertificateSignBytes(invBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig3 := tValSigner3.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := block.NewCertificate(invBlockHash, 0, committers, []int{}, aggSig)

		assert.Error(t, tState1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Invalid round", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 1)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig3 := tValSigner3.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := block.NewCertificate(nextBlockHash, 1, committers, []int{}, aggSig)

		assert.Error(t, tState1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Ok", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig3 := tValSigner3.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := block.NewCertificate(nextBlockHash, 0, committers, []int{}, aggSig)

		assert.NoError(t, tState1.UpdateLastCertificate(cert))
	})
}

func TestBlockValidation(t *testing.T) {
	setup(t)

	moveToNextHeightForAllStates(t)

	assert.False(t, tState1.lastInfo.BlockHash().EqualsTo(hash.UndefHash))

	//
	// Version   			(OK)
	// UnixTime				(TestValidateBlockTime)
	// PrevBlockHash		(OK)
	// StateRoot			(OK)
	// TxsRoot			    (SanityCheck)
	// PrevCertificateHash	(OK)
	// SortitionSeed		(OK)
	// ProposerAddress		(OK)
	//
	invAddr := crypto.GenerateTestAddress()
	invHash := hash.GenerateTestHash()
	invCert := block.GenerateTestCertificate(tState1.lastInfo.BlockHash())
	invSeed := sortition.GenerateRandomSeed()
	trx := tState2.createSubsidyTx(0)
	assert.NoError(t, tState2.AddPendingTx(trx))
	txs := block.NewTxs()
	txs.Append(trx)

	b := block.MakeBlock(2, util.Now(), txs, invHash, tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), invHash, tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), invCert, tState1.lastInfo.SortitionSeed(), tState2.signer.Address())
	assert.Error(t, tState1.validateBlock(b))

	b = block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), invAddr)
	assert.NoError(t, tState1.validateBlock(b))
	c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	b = block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), invSeed, tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.Error(t, tState1.CommitBlock(2, b, c))

	seed := tState1.lastInfo.SortitionSeed()
	b = block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(), tState1.lastInfo.Certificate(), seed.Generate(tState2.signer), tState2.signer.Address())
	assert.NoError(t, tState1.validateBlock(b))
	c = makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, tState1.CommitBlock(2, b, c))
}
