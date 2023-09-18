package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func aggregate(sigs []crypto.Signature) *bls.Signature {
	blsSigs := make([]*bls.Signature, len(sigs))
	for i, s := range sigs {
		blsSigs[i] = s.(*bls.Signature)
	}
	return bls.SignatureAggregate(blsSigs...)
}

func TestCertificateValidation(t *testing.T) {
	td := setup(t)
	td.moveToNextHeightForAllStates(t)
	td.moveToNextHeightForAllStates(t)
	td.moveToNextHeightForAllStates(t)
	td.moveToNextHeightForAllStates(t)
	td.moveToNextHeightForAllStates(t)

	val5, signer5 := td.GenerateTestValidator(4)
	td.state1.store.UpdateValidator(val5)
	td.state2.store.UpdateValidator(val5)

	nextBlock, _ := td.state2.ProposeBlock(td.state2.signers[0], td.RandAddress(), 0)
	nextBlockHash := nextBlock.Hash()
	height := uint32(6)
	round := int16(0)

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, round)
		aggSig := signer5.SignData(signBytes).(*bls.Signature)
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Invalid round, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, round+1)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Invalid height, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height+1, round)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height+1, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Invalid block hash, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		invBlockHash := td.RandHash()
		signBytes := certificate.BlockCertificateSignBytes(invBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		committers = append(committers, 666)
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Invalid committers, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		committers[0] = val5.Number()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := signer5.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2, 3}, aggSig)

		assert.Error(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.NoError(t, td.state1.CommitBlock(height, nextBlock, cert))
	})

	t.Run("Update last certificate, Invalid committers", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		committers = append(committers, val5.Number())
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig3 := td.valSigner3.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		sig5 := signer5.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4, sig5})
		cert := certificate.NewCertificate(height, 0, committers, []int32{}, aggSig)

		assert.Error(t, td.state1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Invalid block hash", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		invBlockHash := td.RandHash()
		signBytes := certificate.BlockCertificateSignBytes(invBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig3 := td.valSigner3.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{}, aggSig)

		assert.Error(t, td.state1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Invalid round", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 1)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig3 := td.valSigner3.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := certificate.NewCertificate(height, 1, committers, []int32{}, aggSig)

		assert.Error(t, td.state1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Invalid height", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height+1, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig3 := td.valSigner3.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := certificate.NewCertificate(height+1, 0, committers, []int32{}, aggSig)

		assert.Error(t, td.state1.UpdateLastCertificate(cert))
	})

	t.Run("Update last certificate, Ok", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valSigner1.SignData(signBytes)
		sig2 := td.valSigner2.SignData(signBytes)
		sig3 := td.valSigner3.SignData(signBytes)
		sig4 := td.valSigner4.SignData(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{}, aggSig)

		assert.NoError(t, td.state1.UpdateLastCertificate(cert))
	})
}

func TestBlockValidation(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	assert.False(t, td.state1.lastInfo.BlockHash().EqualsTo(hash.UndefHash))

	//
	// Version   			(OK)
	// UnixTime				(TestValidateBlockTime)
	// PrevBlockHash		(OK)
	// StateRoot			(OK)
	// TxsRoot			    (BasicCheck)
	// PrevCertificate   	(OK)
	// SortitionSeed		(OK)
	// ProposerAddress		(OK)
	//
	proposerAddr := td.state2.signers[0].Address()
	trx := td.state2.createSubsidyTx(td.RandAddress(), 0)
	txs := block.NewTxs()
	txs.Append(trx)

	t.Run("Invalid version", func(t *testing.T) {
		b := block.MakeBlock(2, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			td.state1.lastInfo.Certificate(), td.state1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, td.state1.validateBlock(b), "Invalid Version")
	})

	t.Run("Invalid StateRoot", func(t *testing.T) {
		invHash := td.RandHash()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), invHash,
			td.state1.lastInfo.Certificate(), td.state1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, td.state1.validateBlock(b), "Invalid StateRoot")
	})

	t.Run("Invalid PrevCertificate", func(t *testing.T) {
		invCert := td.GenerateTestCertificate()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			invCert, td.state1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, td.state1.validateBlock(b), "Invalid PrevCertificate")
	})

	t.Run("Invalid ProposerAddress", func(t *testing.T) {
		invAddr := td.RandAddress()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			td.state1.lastInfo.Certificate(), td.state1.lastInfo.SortitionSeed(), invAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valSigner1, td.valSigner2, td.valSigner3, td.valSigner4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.Error(t, td.state1.CommitBlock(2, b, c), "Invalid ProposerAddress")
	})

	t.Run("Invalid SortitionSeed", func(t *testing.T) {
		invSeed := td.RandSeed()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			td.state1.lastInfo.Certificate(), invSeed, proposerAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valSigner1, td.valSigner2, td.valSigner3, td.valSigner4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.Error(t, td.state1.CommitBlock(2, b, c), "Invalid SortitionSeed")
	})

	t.Run("Ok", func(t *testing.T) {
		seed := td.state1.lastInfo.SortitionSeed()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			td.state1.lastInfo.Certificate(), seed.GenerateNext(td.state2.signers[0]), proposerAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valSigner1, td.valSigner2, td.valSigner3, td.valSigner4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.NoError(t, td.state1.CommitBlock(2, b, c), "Looks Good")
	})
}
