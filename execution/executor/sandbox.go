package executor

import (
	"gitlab.com/zarb-chain/zarb-go/account"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/validator"
)

type Sandbox interface {
	HasAccount(addr crypto.Address) bool
	Account(addr crypto.Address) *account.Account
	UpdateAccount(acc *account.Account)

	HasValidator(addr crypto.Address) bool
	Validator(addr crypto.Address) *validator.Validator
	UpdateValidator(val *validator.Validator)
}
