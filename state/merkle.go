package state

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	simplemerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/validator"
)

func (st *state) accountsMerkleRootHash() crypto.Hash {
	total := st.store.TotalAccounts()

	hashes := make([]crypto.Hash, total)
	st.store.IterateAccounts(func(acc *account.Account) (stop bool) {
		if acc.Number() >= total {
			panic("Account number is out of range")
		}
		if !hashes[acc.Number()].IsUndef() {
			panic("Duplicated account number")
		}
		hashes[acc.Number()] = acc.Hash()

		return false
	})

	tree := simplemerkle.NewTreeFromHashes(hashes)
	return tree.Root()
}

func (st *state) validatorsMerkleRootHash() crypto.Hash {
	total := st.store.TotalValidators()
	hashes := make([]crypto.Hash, total)
	st.store.IterateValidators(func(val *validator.Validator) (stop bool) {
		if val.Number() >= total {
			panic("Validator number is out of range")
		}
		if !hashes[val.Number()].IsUndef() {
			panic("Duplicated validator number")
		}

		hashes[val.Number()] = val.Hash()
		return false
	})
	tree := simplemerkle.NewTreeFromHashes(hashes)
	return tree.Root()
}

func (st *state) stateHash() crypto.Hash {
	accRootHash := st.accountsMerkleRootHash()
	valRootHash := st.validatorsMerkleRootHash()

	rootHash := simplemerkle.HashMerkleBranches(&accRootHash, &valRootHash)

	return *rootHash
}

func (st *state) calculateGenesisStateHashFromGenesisDoc() crypto.Hash {
	accs := st.genDoc.Accounts()
	vals := st.genDoc.Validators()

	accHashes := make([]crypto.Hash, len(accs))
	valHashes := make([]crypto.Hash, len(vals))
	for i, acc := range accs {
		accHashes[i] = acc.Hash()
	}
	for i, val := range vals {
		valHashes[i] = val.Hash()
	}

	accTree := simplemerkle.NewTreeFromHashes(accHashes)
	valTree := simplemerkle.NewTreeFromHashes(valHashes)
	accRootHash := accTree.Root()
	valRootHash := valTree.Root()

	return *simplemerkle.HashMerkleBranches(&accRootHash, &valRootHash)
}
