package state

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestBlockValidation(t *testing.T) {
	td := setup(t)

	round := td.RandRound()
	t.Run("Invalid version", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invBlockVersion := uint8(2)
		blk := block.MakeBlock(
			invBlockVersion,
			blk0.Header().Time(),
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			blk0.Header().StateRoot(),
			blk0.PrevCertificate(),
			blk0.Header().SortitionSeed(),
			blk0.Header().ProposerAddress())
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))

		// Receiving a block with version 2 and rejects it.
		// It is possible that the same block would be considered valid by other nodes (Soft fork).
		assert.Error(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Invalid time", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invBlockTime := util.RoundNow(td.state.params.BlockIntervalInSecond).Add(30 * time.Second)
		blk := block.MakeBlock(
			blk0.Header().Version(),
			invBlockTime,
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			blk0.Header().StateRoot(),
			blk0.PrevCertificate(),
			blk0.Header().SortitionSeed(),
			blk0.Header().ProposerAddress())
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))
		assert.NoError(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Invalid StateRoot", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invStateRoot := td.RandHash()
		blk := block.MakeBlock(
			blk0.Header().Version(),
			blk0.Header().Time(),
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			invStateRoot,
			blk0.PrevCertificate(),
			blk0.Header().SortitionSeed(),
			blk0.Header().ProposerAddress())
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))
		assert.Error(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Invalid PrevCertificate", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invPrevCert := certificate.NewBlockCertificate(
			blk0.PrevCertificate().Height(),
			blk0.PrevCertificate().Round(),
		)
		invPrevCert.SetSignature(
			blk0.PrevCertificate().Committers(),
			blk0.PrevCertificate().Absentees(),
			td.RandBLSSignature())

		blk := block.MakeBlock(
			blk0.Header().Version(),
			blk0.Header().Time(),
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			blk0.Header().StateRoot(),
			invPrevCert,
			blk0.Header().SortitionSeed(),
			blk0.Header().ProposerAddress())
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))
		assert.Error(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Invalid ProposerAddress", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invProposerAddress := td.RandValAddress()
		blk := block.MakeBlock(
			blk0.Header().Version(),
			blk0.Header().Time(),
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			blk0.Header().StateRoot(),
			blk0.PrevCertificate(),
			blk0.Header().SortitionSeed(),
			invProposerAddress)
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))
		assert.Error(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Invalid SortitionSeed", func(t *testing.T) {
		blk0, _ := td.makeBlockAndCertificate(t, round)
		invSortitionSeed, _ := sortition.VerifiableSeedFromBytes(td.RandBLSSignature().Bytes())
		blk := block.MakeBlock(
			blk0.Header().Version(),
			blk0.Header().Time(),
			blk0.Transactions(),
			blk0.Header().PrevBlockHash(),
			blk0.Header().StateRoot(),
			blk0.PrevCertificate(),
			invSortitionSeed,
			blk0.Header().ProposerAddress())
		cert := td.makeCertificateAndSign(t, blk.Hash(), round)

		assert.Error(t, td.state.ValidateBlock(blk, round))
		assert.Error(t, td.state.CommitBlock(blk, cert))
	})

	t.Run("Ok", func(t *testing.T) {
		blk, cert := td.makeBlockAndCertificate(t, round)

		assert.NoError(t, td.state.ValidateBlock(blk, round))
		assert.NoError(t, td.state.CommitBlock(blk, cert))
	})
}
