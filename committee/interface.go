package committee

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Committers() []int32
	Contains(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	ProtocolVersions() map[protocol.Version]float64
	Size() int
	TotalPower() int64
	String() string
}

type Committee interface {
	Reader

	Update(lastRound int16, joined []*validator.Validator)
}
