package sandbox

import (
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
)

type Sandbox interface {
	Account(crypto.Address) *account.Account
	MakeNewAccount(crypto.Address) *account.Account
	UpdateAccount(crypto.Address, *account.Account)

	Validator(crypto.Address) *validator.Validator
	MakeNewValidator(*bls.PublicKey) *validator.Validator
	UpdateValidator(*validator.Validator)
	JoinedToCommittee(crypto.Address)
	IsJoinedCommittee(crypto.Address) bool

	VerifyProof(hash.Stamp, sortition.Proof, *validator.Validator) bool
	Committee() committee.Reader
	FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool)
	FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool)

	Params() param.Params
	CurrentHeight() uint32

	IterateAccounts(consumer func(crypto.Address, *account.Account, bool))
	IterateValidators(consumer func(*validator.Validator, bool, bool))
}
