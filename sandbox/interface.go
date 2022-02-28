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
	HasAnyValidatorJoinedCommittee() bool
	JoinCommittee(crypto.Address)

	FindBlockInfoByStamp(stamp hash.Stamp) (int, hash.Hash)
	CommitteeHasFreeSeats() bool
	CommitteeStake() int64
	CommitteeSize() int
	UnbondInterval() int
	BondInterval() int
	CurrentHeight() int
	BlockHeight(hash.Hash) int
	TransactionToLiveInterval() int
	MaxMemoLength() int
	FeeFraction() float64
	MinFee() int64
}
