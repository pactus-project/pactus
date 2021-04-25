package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
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

	if !block.Header().LastBlockHash().EqualsTo(st.lastInfo.BlockHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"last block hash is not same as we expected. Expected %v, got %v", st.lastInfo.BlockHash(), block.Header().LastBlockHash())
	}

	if !block.Header().StateHash().EqualsTo(st.stateHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"state hash is not same as we expected. Expected %v, got %v", st.stateHash(), block.Header().StateHash())
	}

	if err := st.validateCertificateForPreviousHeight(block.LastCertificate()); err != nil {
		return err
	}

	return nil
}

func (st *state) validateCertificate(cert *block.Certificate) error {
	if err := cert.SanityCheck(); err != nil {
		return err
	}

	pubs := make([]crypto.PublicKey, 0, len(cert.Committers()))
	totalStake := int64(0)
	signersStake := int64(0)

	for _, num := range cert.Committers() {
		val, _ := st.store.ValidatorByNumber(num)
		if val == nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"certificate has invalid committer: %x", num)
		}
		if !util.HasItem(cert.Absences(), num) {
			pubs = append(pubs, val.PublicKey())
			signersStake += val.Power()
		}
		totalStake += val.Power()
	}

	// Check if signers have 2/3+ of total stake
	if signersStake <= totalStake*2/3 {
		return errors.Errorf(errors.ErrInvalidBlock, "No quorom. Has %v, should be more than %v", signersStake, totalStake*2/3)
	}

	// Check signature
	signBytes := cert.SignBytes()
	if !crypto.VerifyAggregated(cert.Signature(), pubs, signBytes) {
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
		if err := st.validateCertificate(cert); err != nil {
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

// validateCertificateForCurrentHeight validates certificate for the current height
func (st *state) validateCertificateForCurrentHeight(cert *block.Certificate, blockHash crypto.Hash) error {
	if err := st.validateCertificate(cert); err != nil {
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
