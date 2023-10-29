package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
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

	val5, valKey5 := td.GenerateTestValidator(4)
	td.state1.store.UpdateValidator(val5)
	td.state2.store.UpdateValidator(val5)

	nextBlock, _ := td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress())
	nextBlockHash := nextBlock.Hash()
	height := uint32(6)
	round := int16(0)

	t.Run("Invalid signature, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, round)
		aggSig := valKey5.Sign(signBytes)
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Invalid round, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, round+1)
		sig1 := td.valKey1.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		sig4 := td.valKey4.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Invalid block hash, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		invBlockHash := td.RandHash()
		signBytes := certificate.BlockCertificateSignBytes(invBlockHash, height, 0)
		sig1 := td.valKey1.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		sig4 := td.valKey4.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		committers = append(committers, 666)
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valKey1.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		sig4 := td.valKey4.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Invalid committers, should return error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		committers[0] = val5.Number()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := valKey5.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		sig4 := td.valKey4.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Doesn't have 2/3 majority", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valKey1.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2, 3}, aggSig)

		assert.Error(t, td.state1.CommitBlock(nextBlock, cert))
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		committers := td.state2.committee.Committers()
		signBytes := certificate.BlockCertificateSignBytes(nextBlockHash, height, 0)
		sig1 := td.valKey1.Sign(signBytes)
		sig2 := td.valKey2.Sign(signBytes)
		sig4 := td.valKey4.Sign(signBytes)
		aggSig := aggregate([]crypto.Signature{sig1, sig2, sig4})
		cert := certificate.NewCertificate(height, 0, committers, []int32{2}, aggSig)

		assert.NoError(t, td.state1.CommitBlock(nextBlock, cert))
	})
}

func TestBlockValidation(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	assert.False(t, td.state1.lastInfo.BlockHash().IsUndef())

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
	prevCert := td.state1.lastInfo.Certificate()
	proposerAddr := td.state2.valKeys[0].Address()
	trx := td.state2.createSubsidyTx(td.RandAccAddress(), 0)
	txs := block.NewTxs()
	txs.Append(trx)

	t.Run("Invalid version", func(t *testing.T) {
		b := block.MakeBlock(2, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			prevCert, td.state1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, td.state1.validateBlock(b), "Invalid Version")
	})

	t.Run("Invalid StateRoot", func(t *testing.T) {
		invHash := td.RandHash()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), invHash,
			prevCert, td.state1.lastInfo.SortitionSeed(), proposerAddr)

		assert.Error(t, td.state1.validateBlock(b), "Invalid StateRoot")
	})

	t.Run("Invalid PrevCertificate round", func(t *testing.T) {
		invCert := certificate.NewCertificate(prevCert.Height(), prevCert.Round()+1,
			prevCert.Committers(), prevCert.Absentees(), prevCert.Signature())
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			invCert, td.state1.lastInfo.SortitionSeed(), proposerAddr)

		err := td.state1.validateBlock(b)
		assert.ErrorIs(t, InvalidCertificateError{Cert: invCert}, err)
	})

	t.Run("Invalid PrevCertificate signature", func(t *testing.T) {
		invCert := certificate.NewCertificate(prevCert.Height(), prevCert.Round(),
			prevCert.Committers(), prevCert.Absentees(), td.RandBLSSignature())
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			invCert, td.state1.lastInfo.SortitionSeed(), proposerAddr)

		err := td.state1.validateBlock(b)
		assert.ErrorIs(t, crypto.ErrInvalidSignature, err)
	})

	t.Run("Invalid ProposerAddress", func(t *testing.T) {
		invAddr := td.RandAccAddress()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			prevCert, td.state1.lastInfo.SortitionSeed(), invAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.Error(t, td.state1.CommitBlock(b, c), "Invalid ProposerAddress")
	})

	t.Run("Invalid SortitionSeed", func(t *testing.T) {
		invSeed := td.RandSeed()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			prevCert, invSeed, proposerAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.Error(t, td.state1.CommitBlock(b, c), "Invalid SortitionSeed")
	})

	t.Run("Ok", func(t *testing.T) {
		seed := td.state1.lastInfo.SortitionSeed()
		b := block.MakeBlock(1, util.Now(), txs, td.state1.lastInfo.BlockHash(), td.state1.stateRoot(),
			prevCert, seed.GenerateNext(td.state2.valKeys[0].PrivateKey()), proposerAddr)
		c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)

		assert.NoError(t, td.state1.validateBlock(b))
		assert.NoError(t, td.state1.CommitBlock(b, c), "Looks Good")
	})
}
