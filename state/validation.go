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

	if !block.Header().CommitersHash().EqualsTo(st.validatorSet.CommitersHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Commiters hash is not same as we expected. Expected %v, got %v", st.validatorSet.CommitersHash(), block.Header().CommitersHash())
	}

	if !block.Header().StateHash().EqualsTo(st.stateHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"State hash is not same as we expected. Expected %v, got %v", st.stateHash(), block.Header().StateHash())
	}

	if err := st.validateLastCommit(block.Header().LastCommit(), st.lastBlockHash); err != nil {
		return err
	}

	return nil
}

func (st *state) validateLastCommit(commit *block.Commit, blockHash crypto.Hash) error {
	if commit == nil {
		if !blockHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Only genesis block has no commit")
		}
	} else {
		if err := commit.SanityCheck(); err != nil {
			return err
		}

		if !commit.CommitersHash().EqualsTo(st.lastCommit.CommitersHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last commiters are not same as we expected. Expected %v, got %v", st.lastCommit.CommitersHash(), commit.CommitersHash())
		}

		if commit.Round() != st.lastCommit.Round() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last commit round is not same as we expected. Expected %v, got %v", st.lastCommit.Round(), commit.Round())
		}

		if !blockHash.EqualsTo(st.lastBlockHash) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Last commit block hash is not same as we expected. Expected %v, got %v", st.lastBlockHash, blockHash)

		}

		signBytes := vote.CommitSignBytes(blockHash, commit.Round())
		pubs := make([]crypto.PublicKey, 0)
		for _, c := range commit.Commiters() {
			if c.Signed {
				val, _ := st.store.Validator(c.Address)
				if val == nil {
					return errors.Errorf(errors.ErrInvalidBlock,
						"invalid commiter: %x", c.Address)
				}
				pubs = append(pubs, val.PublicKey())
			}
		}

		if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"invalid commit signature: %x", commit.Signature())
		}
	}

	return nil
}

func (st *state) validateCommit(commit block.Commit, blockHash crypto.Hash) error {
	if err := commit.SanityCheck(); err != nil {
		return err
	}

	if !commit.CommitersHash().EqualsTo(st.validatorSet.CommitersHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"Last commiters are not same as we expected. Expected %v, got %v", st.validatorSet.CommitersHash(), commit.CommitersHash())
	}

	signBytes := vote.CommitSignBytes(blockHash, commit.Round())
	pubs := make([]crypto.PublicKey, 0)
	for _, c := range commit.Commiters() {
		if c.Signed {
			val := st.validatorSet.Validator(c.Address)
			if val == nil {
				return errors.Errorf(errors.ErrInvalidBlock,
					"invalid commiter: %x", c.Address)
			}
			pubs = append(pubs, val.PublicKey())
		}
	}

	if !crypto.VerifyAggregated(commit.Signature(), pubs, signBytes) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"invalid commit signature: %x", commit.Signature())
	}

	return nil
}
