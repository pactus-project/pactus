package committee

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Committers() []int32
	Contains(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	Size() int
	TotalPower() int64
}

type Committee interface {
	Reader

	Update(lastround int16, joined []*validator.Validator)
}
