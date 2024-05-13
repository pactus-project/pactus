package state

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/errors"
)

func (st *state) doValidateBlock(blk *block.Block, round int16) error {
	if blk.Header().Version() != st.params.BlockVersion {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid version")
	}

	if blk.Header().StateRoot() != st.stateRoot() {
		return errors.Errorf(errors.ErrInvalidBlock,
			"state root is not same as we expected, expected %v, got %v", st.stateRoot(), blk.Header().StateRoot())
	}

	// Verify proposer
	proposer := st.committee.Proposer(round)
	if proposer.Address() != blk.Header().ProposerAddress() {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid proposer, expected %s, got %s", proposer.Address(), blk.Header().ProposerAddress())
	}
	// Validate sortition seed
	seed := blk.Header().SortitionSeed()
	if !seed.Verify(proposer.PublicKey(), st.lastInfo.SortitionSeed()) {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid sortition seed")
	}

	return st.validatePrevCertificate(blk.PrevCertificate(), blk.Header().PrevBlockHash())
}

// validatePrevCertificate validates certificate for the previous block.
func (st *state) validatePrevCertificate(cert *certificate.Certificate, blockHash hash.Hash) error {
	if cert == nil {
		if !st.lastInfo.BlockHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"only genesis block has no certificate")
		}
	} else {
		if cert.Round() != st.lastInfo.Certificate().Round() {
			// TODO: we should panic here.
			// It is impossible, unless we have a fork on the latest block
			return InvalidCertificateError{
				Cert: cert,
			}
		}

		signBytes := certificate.BlockCertificateSignBytes(blockHash, cert.Height(), cert.Round())
		err := cert.Validate(st.lastInfo.BlockHeight(), st.lastInfo.Validators(), signBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateCertificate validates certificate for the current height.
func (st *state) validateCertificate(cert *certificate.Certificate, blockHash hash.Hash) error {
	signBytes := certificate.BlockCertificateSignBytes(blockHash, cert.Height(), cert.Round())
	err := cert.Validate(st.lastInfo.BlockHeight()+1, st.committee.Validators(), signBytes)
	if err != nil {
		return err
	}

	return nil
}
