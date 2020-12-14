package state

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/logger"
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

	tree := simpleMerkle.NewTreeFromHashes(hashes)
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
	tree := simpleMerkle.NewTreeFromHashes(hashes)
	return tree.Root()
}

func (st *state) stateHash() crypto.Hash {
	accRootHash := st.accountsMerkleRootHash()
	valRootHash := st.validatorsMerkleRootHash()

	rootHash := simpleMerkle.HashMerkleBranches(&accRootHash, &valRootHash)
	if rootHash == nil {
		logger.Panic("State hash can't be nil")
	}

	return *rootHash
}
