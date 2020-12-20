package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/vote"
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

	if err := st.validateLastCommit(block.LastCommit()); err != nil {
		return err
	}

	return nil
}

func (st *state) validateLastCommit(commit *block.Commit) error {
	if commit == nil {
		if !st.lastBlockHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Only genesis block has no commit")
		}
	} else {
		if err := commit.SanityCheck(); err != nil {
			return err
		}

		// TODO: add tests for this case
		// Make sure the committers are the cprrect one
		if !commit.CommittersHash().EqualsTo(st.lastCommit.CommittersHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last committers are not same as we expected. Expected %v, got %v", st.lastCommit.CommittersHash(), commit.CommittersHash())
		}

		if commit.Round() != st.lastCommit.Round() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last commit round is not same as we expected. Expected %v, got %v", st.lastCommit.Round(), commit.Round())
		}

		signBytes := vote.CommitSignBytes(st.lastBlockHash, commit.Round())
		pubs := make([]crypto.PublicKey, 0)
		for _, c := range commit.Committers() {
			if c.HasSigned() {
				val, _ := st.store.Validator(c.Address)
				if val == nil {
					return errors.Errorf(errors.ErrInvalidBlock,
						"invalid committer: %x", c.Address)
				}
				pubs = append(pubs, val.PublicKey())
			}
		}

		if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"invalid commit signature: %v", commit.Signature())
		}
	}

	return nil
}

// validateCommitForCurrentHeight validates commit for the current height
func (st *state) validateCommitForCurrentHeight(commit block.Commit, blockHash crypto.Hash) error {
	if err := commit.SanityCheck(); err != nil {
		return err
	}

	if !commit.CommittersHash().EqualsTo(st.validatorSet.CommittersHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last committers are not same as we expected. Expected %v, got %v", st.validatorSet.CommittersHash(), commit.CommittersHash())
	}

	signBytes := vote.CommitSignBytes(blockHash, commit.Round())
	pubs := make([]crypto.PublicKey, 0)
	for _, c := range commit.Committers() {
		if c.HasSigned() {
			// Since this block belongs to current height, we get validator info from validator set
			val := st.validatorSet.Validator(c.Address)
			if val == nil {
				return errors.Errorf(errors.ErrInvalidBlock,
					"invalid committer: %x", c.Address)
			}
			pubs = append(pubs, val.PublicKey())
		}
	}

	if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid commit signature: %v", commit.Signature())
	}

	return nil
}
