package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/errors"
)

type QueryProposalMessage struct {
	Height  uint32         `cbor:"1,keyasint"`
	Round   int16          `cbor:"3,keyasint"`
	Querier crypto.Address `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(height uint32, round int16, querier crypto.Address) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height:  height,
		Round:   round,
		Querier: querier,
	}
}

func (m *QueryProposalMessage) BasicCheck() error {
	if m.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}

	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return TypeQueryProposal
}

func (m *QueryProposalMessage) String() string {
	return fmt.Sprintf("{%v %s}", m.Height, m.Querier.ShortString())
}
