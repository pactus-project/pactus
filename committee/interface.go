package committee

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Committers() []int
	Contains(addr crypto.Address) bool
	Proposer(round int) *validator.Validator
	IsProposer(addr crypto.Address, round int) bool
	Size() int
	TotalPower() int64
}

type Committee interface {
	Reader

	Update(lastRound int, joined []*validator.Validator)
}
