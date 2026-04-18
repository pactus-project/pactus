package committee

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/validator"
)

type Reader interface {
	Validators() []*validator.Validator
	Committers() []int32
	Contains(addr crypto.Address) bool
	Proposer(round types.Round) *validator.Validator
	IsProposer(addr crypto.Address, round types.Round) bool
	ProtocolVersions() map[protocol.Version]float64
	SupportProtocolVersion(version protocol.Version) bool
	Size() int
	TotalPower() int64
	String() string
}

type Committee interface {
	Reader

	Update(lastround types.Round, joined []*validator.Validator)
}
