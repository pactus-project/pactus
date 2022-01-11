package sandbox

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/validator"
)

type Sandbox interface {
	Account(crypto.Address) *account.Account
	MakeNewAccount(crypto.Address) *account.Account
	UpdateAccount(*account.Account)

	Validator(crypto.Address) *validator.Validator
	MakeNewValidator(*bls.PublicKey) *validator.Validator
	UpdateValidator(*validator.Validator)
	IsInCommittee(crypto.Address) bool

	VerifySortition(hash.Hash, sortition.Proof, *validator.Validator) bool
	EnterCommittee(hash.Hash, crypto.Address) error

	FindBlockInfoByStamp(stamp hash.Stamp) (int, hash.Hash)
	CommitteeSize() int
	UnbondInterval() int
	CurrentHeight() int
	BlockHeight(hash.Hash) int
	TransactionToLiveInterval() int
	MaxMemoLength() int
	FeeFraction() float64
	MinFee() int64
}
