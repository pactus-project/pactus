package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

func (st *State) validateBlock(block *block.Block) error {
	if block == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "Block is empty")
	}

	if err := block.SanityCheck(); err != nil {
		return err
	}

	if !block.Header().LastBlockHash().EqualsTo(st.lastBlockHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"lastBlockHash.  Expected %v, got %v", st.lastBlockHash, block.Header().LastBlockHash())
	}

	if !block.Header().LastReceiptsHash().EqualsTo(st.lastReceiptsHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"lastReceiptsHash.  Expected %v, got %v", st.lastReceiptsHash, block.Header().LastReceiptsHash())
	}

	if !block.Header().NextValidatorsHash().EqualsTo(st.nextValidatorsHash) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"NextValidatorsHash.  Expected %v, got %v", st.nextValidatorsHash, block.Header().NextValidatorsHash())
	}

	if !block.Header().StateHash().EqualsTo(st.stateHash()) {
		return errors.Errorf(errors.ErrInvalidBlock,
			"StateHash.  Expected %v, got %v", st.stateHash(), block.Header().StateHash())
	}

	if block.Header().LastBlockHash().IsUndef() {
		if block.LastCommit() != nil {
			return errors.Errorf(errors.ErrInvalidBlock,
				"block at height 1 can't have Commit signatures")
		}
	} else {
		// TODO
		// Verify commit signers
	}

	// TODO: Validate block Time

	// TODO: validate proposer is correct

	return nil
}
