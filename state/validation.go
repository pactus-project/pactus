package state

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
)

func (st *state) validateBlock(blk *block.Block, round int16) error {
	if blk.Header().Version() > protocol.ProtocolVersionLatest ||
		blk.Header().Version() < st.params.BlockVersion {
		return ErrInvalidBlockVersion
	}

	if blk.Header().StateRoot() != st.stateRoot() {
		return InvalidStateRootHashError{
			Expected: st.stateRoot(),
			Got:      blk.Header().StateRoot(),
		}
	}

	// Verify proposer
	proposer := st.committee.Proposer(round)
	if proposer.Address() != blk.Header().ProposerAddress() {
		return InvalidProposerError{
			Expected: proposer.Address(),
			Got:      blk.Header().ProposerAddress(),
		}
	}

	// Validate sortition seed
	seed := blk.Header().SortitionSeed()
	if !seed.Verify(proposer.PublicKey(), st.lastInfo.SortitionSeed()) {
		return ErrInvalidSortitionSeed
	}

	return st.validatePrevCertificate(blk.PrevCertificate(), blk.Header().PrevBlockHash())
}

// validatePrevCertificate validates certificate for the previous block.
func (st *state) validatePrevCertificate(cert *certificate.Certificate, blockHash hash.Hash) error {
	if cert == nil {
		if !st.lastInfo.BlockHash().IsUndef() {
			return ErrInvalidCertificate
		}
	} else {
		err := cert.ValidatePrecommit(st.lastInfo.Validators(), blockHash)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateCurCertificate validates certificate for the current height.
func (st *state) validateCurCertificate(cert *certificate.Certificate, blockHash hash.Hash) error {
	err := cert.ValidatePrecommit(st.committee.Validators(), blockHash)
	if err != nil {
		return err
	}

	return nil
}
