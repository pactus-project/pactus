package executor

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type Sandbox interface {
	HasAccount(addr crypto.Address) bool
	Account(addr crypto.Address) *account.Account
	UpdateAccount(acc *account.Account)

	HasValidator(addr crypto.Address) bool
	Validator(addr crypto.Address) *validator.Validator
	UpdateValidator(val *validator.Validator)
}
