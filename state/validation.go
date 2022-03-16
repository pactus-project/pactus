package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

func (st *state) validateBlock(block *block.Block) error {
	if err := block.SanityCheck(); err != nil {
		return err
	}

	if block.Header().Version() != st.params.BlockVersion {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid version")
	}

	if !block.Header().StateRoot().EqualsTo(st.stateRoot()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"state root is not same as we expected. Expected %v, got %v", st.stateRoot(), block.Header().StateRoot())
	}

	return st.validateCertificateForPreviousHeight(block.PrevCertificate())
}

func (st *state) checkCertificate(cert *block.Certificate) error {
	if err := cert.SanityCheck(); err != nil {
		return err
	}

	pubs := make([]*bls.PublicKey, 0, len(cert.Committers()))
	committeePower := int64(0)
	signedPower := int64(0)

	for _, num := range cert.Committers() {
		val, _ := st.store.ValidatorByNumber(num)
		if val == nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"certificate has invalid committer: %x", num)
		}
		if !util.Contains(cert.Absentees(), num) {
			pubs = append(pubs, val.PublicKey())
			signedPower += val.Power()
		}
		committeePower += val.Power()
	}

	// Check if signers have 2/3+ of total power
	if signedPower <= committeePower*2/3 {
		return errors.Errorf(errors.ErrInvalidBlock, "No quorom. Has %v, should be more than %v", signedPower, committeePower*2/3)
	}

	// Check signature
	signBytes := cert.SignBytes()
	if !bls.VerifyAggregated(cert.Signature(), pubs, signBytes) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"certificate has invalid signature: %v", cert.Signature())
	}

	return nil
}

// validateCertificateForPreviousHeight validates certificate for the previous height
func (st *state) validateCertificateForPreviousHeight(cert *block.Certificate) error {
	if cert == nil {
		if !st.lastInfo.BlockHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"only genesis block has no certificate")
		}
	} else {
		if err := st.checkCertificate(cert); err != nil {
			return err
		}

		if !cert.BlockHash().EqualsTo(st.lastInfo.BlockHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"certificate has invalid block hash. Expected %v, got %v", st.lastInfo.BlockHash(), cert.BlockHash())
		}

		if cert.Round() != st.lastInfo.Certificate().Round() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"certificate has invalid round. Expected %v, got %v", st.lastInfo.Certificate().Round(), cert.Round())
		}

		if !util.Equal(cert.Committers(), st.lastInfo.Certificate().Committers()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"certificate has invalid committers")
		}
	}

	return nil
}

// validateCertificate validates certificate for the current height
func (st *state) validateCertificate(cert *block.Certificate, blockHash hash.Hash) error {
	if err := st.checkCertificate(cert); err != nil {
		return err
	}

	if !cert.BlockHash().EqualsTo(blockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"certificate has invalid block hash. Expected %v, got %v", st.lastInfo.BlockHash(), cert.BlockHash())
	}

	if !util.Equal(st.committee.Committers(), cert.Committers()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid committers")
	}

	return nil
}
