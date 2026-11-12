package state

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
)

func (st *state) validateBlock(blk *block.Block, round types.Round) error {
	if blk.Header().Version() > protocol.ProtocolVersionLatest ||
		blk.Header().Version() < st.params.BlockVersion {
		return InvalidBlockVersionError{
			Version: blk.Header().Version(),
		}
	}

	// After the first halving (block 8,000,000), blocks must use version 4 or higher.
	// This guards against old nodes that didn't upgrade in time.
	if blk.Header().Version() < protocol.ProtocolVersion4 &&
		blk.Height() > 8_000_000 {
		return InvalidBlockVersionError{
			Version: blk.Header().Version(),
		}
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
	if blockHash != st.lastInfo.BlockHash() {
		// The certificate below can still be a validly signed certificate,
		// just for an older block: ValidatePrecommit only checks that
		// blockHash was precommitted by the given validators, not that
		// blockHash is the chain's actual tip.
		return ForkDetectedError{
			Expected: st.lastInfo.BlockHash(),
			Got:      blockHash,
		}
	}

	if cert == nil {
		if !st.lastInfo.BlockHash().IsUndef() {
			return ErrInvalidCertificate
		}
	} else if !st.isLastCertificate(cert) {
		err := cert.ValidatePrecommit(st.lastInfo.Validators(), blockHash)
		if err != nil {
			return err
		}
	}

	return nil
}

// isLastCertificate reports whether the certificate is the one that the state holds
// for the last block. That certificate is verified already: when the last block was
// committed, and each time a vote is added to it. While the node synchronizes, the
// previous certificate of a block is the certificate that committed the block before,
// so verifying it again doubles the signature verifications of each block.
//
// The certificate of a proposer can hold other signers. It is verified as before.
func (st *state) isLastCertificate(cert *certificate.Certificate) bool {
	lastCert := st.lastInfo.Certificate()

	return lastCert != nil && cert.Hash() == lastCert.Hash()
}

// validateCurCertificate validates certificate for the current height.
func (st *state) validateCurCertificate(cert *certificate.Certificate, blockHash hash.Hash) error {
	// The aggregated signature doesn't cover the list of absentees: it must be checked on its own.
	if err := cert.BasicCheck(); err != nil {
		return err
	}

	err := cert.ValidatePrecommit(st.committee.Validators(), blockHash)
	if err != nil {
		return err
	}

	return nil
}
