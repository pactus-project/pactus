package state

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/validator"
)

func (state *State) accountsMerkleRootHash() *crypto.Hash {
	accs := make([]*account.Account, 0)
	state.store.IterateAccounts(func(acc *account.Account) (stop bool) {
		accs = append(accs, acc)
		return false
	})

	hashes := make([]crypto.Hash, len(accs))
	for i, a := range accs {
		hashes[i] = a.Hash()
	}

	tree := simpleMerkle.NewTreeFromHashes(hashes)
	return tree.Root()
}

func (state *State) validatorsMerkleRootHash() *crypto.Hash {
	vals := make([]*validator.Validator, 0)
	state.store.IterateValidators(func(val *validator.Validator) (stop bool) {
		vals = append(vals, val)
		return false
	})

	hashes := make([]crypto.Hash, len(vals))
	for i, v := range vals {
		hashes[i] = v.Hash()
	}

	tree := simpleMerkle.NewTreeFromHashes(hashes)
	return tree.Root()
}
