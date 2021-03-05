package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

func (st *state) validateBlock(block block.Block) error {
	if err := block.SanityCheck(); err != nil {
		return err
	}

	if block.Header().Version() != st.params.BlockVersion {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid version")
	}

	if !block.Header().LastBlockHash().EqualsTo(st.lastBlockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last block hash is not same as we expected. Expected %v, got %v", st.lastBlockHash, block.Header().LastBlockHash())
	}

	if !block.Header().LastReceiptsHash().EqualsTo(st.lastReceiptsHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"last receipts hash is not same as we expected. Expected %v, got %v", st.lastReceiptsHash, block.Header().LastReceiptsHash())
	}

	if !block.Header().CommitteeHash().EqualsTo(st.validatorSet.CommitteeHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Committee hash is not same as we expected. Expected %v, got %v", st.validatorSet.CommitteeHash(), block.Header().CommitteeHash())
	}

	if !block.Header().StateHash().EqualsTo(st.stateHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"State hash is not same as we expected. Expected %v, got %v", st.stateHash(), block.Header().StateHash())
	}

	if err := st.validateCommitForPreviousHeight(block.LastCommit()); err != nil {
		return err
	}

	return nil
}

func (st *state) validateCommit(commit *block.Commit) error {
	if err := commit.SanityCheck(); err != nil {
		return err
	}

	pubs := make([]crypto.PublicKey, 0, len(commit.Committers()))
	totalStake := int64(0)
	signersStake := int64(0)
	for _, c := range commit.Committers() {
		val, _ := st.store.ValidatorByNumber(c.Number)
		if c.HasSigned() {
			if val == nil {
				return errors.Errorf(errors.ErrInvalidBlock,
					"Invalid committer: %x", c.Number)
			}
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
	signBytes := commit.SignBytes()
	if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Invalid commit signature: %v", commit.Signature())
	}

	return nil
}

// validateCommitForPreviousHeight validates commit for the previous height
func (st *state) validateCommitForPreviousHeight(commit *block.Commit) error {
	if commit == nil {
		if !st.lastBlockHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Only genesis block has no commit")
		}
	} else {
		if err := st.validateCommit(commit); err != nil {
			return err
		}

		if !commit.BlockHash().EqualsTo(st.lastBlockHash) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Commit has invalid block hash. Expected %v, got %v", st.lastBlockHash, commit.BlockHash())
		}

		if commit.Round() != st.lastCommit.Round() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last commit round is not same as we expected. Expected %v, got %v", st.lastCommit.Round(), commit.Round())
		}

		if !commit.CommitteeHash().EqualsTo(st.lastCommit.CommitteeHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last committee hash are not same as we expected. Expected %v, got %v", st.lastCommit.CommitteeHash(), commit.CommitteeHash())
		}
	}

	return nil
}

// validateCommitForCurrentHeight validates commit for the current height
func (st *state) validateCommitForCurrentHeight(commit block.Commit, blockHash crypto.Hash) error {
	if err := st.validateCommit(&commit); err != nil {
		return err
	}

	if !commit.BlockHash().EqualsTo(blockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Commit has invalid block hash. Expected %v, got %v", st.lastBlockHash, commit.BlockHash())
	}

	if !commit.CommitteeHash().EqualsTo(st.validatorSet.CommitteeHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last committee hash are not same as we expected. Expected %v, got %v", st.validatorSet.CommitteeHash(), commit.CommitteeHash())
	}

	return nil
}
