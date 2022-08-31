package sandbox

import (
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/param"
	"github.com/zarbchain/zarb-go/types/validator"
)

type Sandbox interface {
	Account(crypto.Address) *account.Account
	MakeNewAccount(crypto.Address) *account.Account
	UpdateAccount(*account.Account)

	Validator(crypto.Address) *validator.Validator
	MakeNewValidator(*bls.PublicKey) *validator.Validator
	UpdateValidator(*validator.Validator)

	VerifyProof(hash.Stamp, sortition.Proof, *validator.Validator) bool
	Committee() committee.Reader
	FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool)
	FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool)

	Params() param.Params
	CurrentHeight() uint32

	IterateAccounts(consumer func(*AccountStatus))
	IterateValidators(consumer func(*ValidatorStatus))
}
