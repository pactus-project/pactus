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

	TotalPower() int64
	IsInCommittee(crypto.Address) bool
	CommitteeAge() int
	CommitteePower() int64
	JoinedPower() int64
	CommitteeHasFreeSeats() bool

	BlockHeightByStamp(stamp hash.Stamp) int
	BlockSeedByStamp(stamp hash.Stamp) sortition.VerifiableSeed

	UnbondInterval() int
	CommitteeSize() int
	BondInterval() int
	CurrentHeight() int
	TransactionToLiveInterval() int
	MaxMemoLength() int
	FeeFraction() float64
	MinFee() int64

	IterateAccounts(consumer func(*AccountStatus))
	IterateValidators(consumer func(*ValidatorStatus))
}
