package committee

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Committers() []int32
	Contains(addr crypto.Address) bool
	Proposer(round int32) *validator.Validator
	IsProposer(addr crypto.Address, round int32) bool
	Size() int
	TotalPower() int64
}

type Committee interface {
	Reader

	Update(lastRound int32, joined []*validator.Validator)
}
