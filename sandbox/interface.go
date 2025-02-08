package sandbox

import (
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
)

type Sandbox interface {
	Account(crypto.Address) *account.Account
	MakeNewAccount(crypto.Address) *account.Account
	UpdateAccount(crypto.Address, *account.Account)

	CommitTransaction(trx *tx.Tx)
	RecentTransaction(txID tx.ID) bool
	IsBanned(trx *tx.Tx) bool

	Validator(crypto.Address) *validator.Validator
	MakeNewValidator(*bls.PublicKey) *validator.Validator
	UpdateValidator(*validator.Validator)
	JoinedToCommittee(crypto.Address)
	IsJoinedCommittee(crypto.Address) bool
	UpdatePowerDelta(delta int64)
	PowerDelta() int64
	AccumulatedFee() amount.Amount

	VerifyProof(uint32, sortition.Proof, *validator.Validator) bool
	Committee() committee.Reader

	Params() *param.Params
	CurrentHeight() uint32

	IterateAccounts(consumer func(crypto.Address, *account.Account, bool))
	IterateValidators(consumer func(*validator.Validator, bool, bool))
}
