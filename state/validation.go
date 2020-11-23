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
			"lastBlockHash. Expected %v, got %v", st.lastBlockHash, block.Header().LastBlockHash())
	}

	if !block.Header().LastReceiptsHash().EqualsTo(st.lastReceiptsHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"lastReceiptsHash. Expected %v, got %v", st.lastReceiptsHash, block.Header().LastReceiptsHash())
	}

	if !block.Header().NextCommitersHash().EqualsTo(st.NextCommitersHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"NextCommitersHash. Expected %v, got %v", st.NextCommitersHash, block.Header().NextCommitersHash())
	}

	if !block.Header().StateHash().EqualsTo(st.stateHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"StateHash. Expected %v, got %v", st.stateHash(), block.Header().StateHash())
	}

	if !block.Header().LastBlockHash().IsUndef() {
		commit := block.Header().LastCommit()
		if !commit.CommitersHash().EqualsTo(st.lastCommit.CommitersHash()) {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Commiters are not same as we expected. Expected %v, got %v", st.lastCommit.CommitersHash(), commit.CommitersHash())
		}

		if commit.Round() != st.lastCommit.Round() {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Commiters round is not same as we expected. Expected %v, got %v", st.lastCommit.Round(), commit.Round())
		}

		if err := st.validateCommit(st.lastBlockHash, *commit); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"Commiters is not valid. %v", err)
		}
	}

	return nil
}

func (st *state) validateCommit(blockHash crypto.Hash, commit block.Commit) error {
	if err := commit.SanityCheck(); err != nil {
		return err
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
