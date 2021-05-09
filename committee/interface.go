package committee

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Contains(addr crypto.Address) bool
	Proposer(round int) *validator.Validator
	IsProposer(addr crypto.Address, round int) bool
	Size() int
	TotalPower() int64
	TotalStake() int64
}
