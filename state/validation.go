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

	if !block.Header().LastBlockHash().EqualsTo(st.lastBlockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last block hash is not same as we expected. Expected %v, got %v", st.lastBlockHash, block.Header().LastBlockHash())
	}

	if !block.Header().LastReceiptsHash().EqualsTo(st.lastReceiptsHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"last receipts hash is not same as we expected. Expected %v, got %v", st.lastReceiptsHash, block.Header().LastReceiptsHash())
	}

	if !block.Header().CommittersHash().EqualsTo(st.validatorSet.CommittersHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Committers hash is not same as we expected. Expected %v, got %v", st.validatorSet.CommittersHash(), block.Header().CommittersHash())
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

	if !commit.HasTwoThirdThreshold() {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Commit has not two third threshold.")
	}

	pubs := make([]crypto.PublicKey, len(commit.Signed()))
	for i, num := range commit.Signed() {
		val, _ := st.store.ValidatorByNumber(num)
		if val == nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"invalid committer: %x", num)
		}
		pubs[i] = val.PublicKey()
	}

	signBytes := commit.SignBytes()
	if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid commit signature: %v", commit.Signature())
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

		// TODO: add tests for this case
		if !commit.CommittersHash().EqualsTo(st.lastCommit.CommittersHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last committers are not same as we expected. Expected %v, got %v", st.lastCommit.CommittersHash(), commit.CommittersHash())
		}
	}

	return nil
}

// validateCommitForCurrentHeight validates commit for the current height
func (st *state) validateCommitForCurrentHeight(commit block.Commit, blockHash crypto.Hash) error {
	if err := st.validateCommit(&commit); err != nil {
		return err
	}

	if !commit.CommittersHash().EqualsTo(st.validatorSet.CommittersHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last committers are not same as we expected. Expected %v, got %v", st.validatorSet.CommittersHash(), commit.CommittersHash())
	}

	if !commit.BlockHash().EqualsTo(blockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Commit has invalid block hash. Expected %v, got %v", st.lastBlockHash, commit.BlockHash())
	}

	for _, num := range commit.Signed() {
		// Check if validator is eligible to commit the block
		val, _ := st.store.ValidatorByNumber(num)
		if val == nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"invalid committer: %x", num)
		}
		if !st.validatorSet.Contains(val.Address()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"validator is eligible to commit the block: %x", val.Address())
		}
	}

	return nil
}
