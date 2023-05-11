package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
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

	nextBlock, _ := tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 0)
	nextBlockHash := nextBlock.Hash()

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		aggSig := signer5.SignData(signBytes).(*bls.Signature)
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Invalid round, should return error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 1)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2})
		cert := block.NewCertificate(0, committers, []int32{2, 3}, aggSig)

		assert.Error(t, tState1.CommitBlock(2, nextBlock, cert))
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		committers := tState2.committee.Committers()
		signBytes := block.CertificateSignBytes(nextBlockHash, 0)
		sig1 := tValSigner1.SignData(signBytes)
		sig2 := tValSigner2.SignData(signBytes)
		sig4 := tValSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := block.NewCertificate(0, committers, []int32{2}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{}, aggSig)

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
		cert := block.NewCertificate(1, committers, []int32{}, aggSig)

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
		cert := block.NewCertificate(0, committers, []int32{}, aggSig)

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
	// PrevCertificate   	(OK)
	// SortitionSeed		(OK)
	// ProposerAddress		(OK)
	//
	proposerAddr := tState2.signers[0].Address()
	trx := tState2.createSubsidyTx(crypto.GenerateTestAddress(), 0)
	txs := block.NewTxs()
	txs.Append(trx)

	t.Run("Invalid version", func(t *testing.T) {
		b := block.MakeBlock(2, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(),
			tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, tState1.validateBlock(b), "Invalid Version")
	})

	t.Run("Invalid StateRoot", func(t *testing.T) {
		invHash := hash.GenerateTestHash()
		b := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), invHash,
			tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, tState1.validateBlock(b), "Invalid StateRoot")
	})

	t.Run("Invalid PrevCertificate", func(t *testing.T) {
		invCert := block.GenerateTestCertificate(tState1.lastInfo.BlockHash())
		b := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(),
			invCert, tState1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, tState1.validateBlock(b), "Invalid PrevCertificate")
	})

	t.Run("Invalid ProposerAddress", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		b := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(),
			tState1.lastInfo.Certificate(), tState1.lastInfo.SortitionSeed(), invAddr)
		c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

		assert.NoError(t, tState1.validateBlock(b))
		assert.Error(t, tState1.CommitBlock(2, b, c), "Invalid ProposerAddress")
	})

	t.Run("Invalid SortitionSeed", func(t *testing.T) {
		invSeed := sortition.GenerateRandomSeed()
		b := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(),
			tState1.lastInfo.Certificate(), invSeed, proposerAddr)
		c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

		assert.NoError(t, tState1.validateBlock(b))
		assert.Error(t, tState1.CommitBlock(2, b, c), "Invalid SortitionSeed")
	})

	t.Run("Ok", func(t *testing.T) {
		seed := tState1.lastInfo.SortitionSeed()
		b := block.MakeBlock(1, util.Now(), txs, tState1.lastInfo.BlockHash(), tState1.stateRoot(),
			tState1.lastInfo.Certificate(), seed.GenerateNext(tState2.signers[0]), proposerAddr)
		c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

		assert.NoError(t, tState1.validateBlock(b))
		assert.NoError(t, tState1.CommitBlock(2, b, c), "Looks Good")
	})
}
