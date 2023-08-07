package message

import (
	"fmt"

	"github.com/pactus-project/pactus/util/errors"
)

type QueryProposalMessage struct {
	Height uint32 `cbor:"1,keyasint"`
	Round  int16  `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(h uint32, r int16) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height: h,
		Round:  r,
	}
}

func (m *QueryProposalMessage) SanityCheck() error {
	if m.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}

	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return TypeQueryProposal
}

func (m *QueryProposalMessage) String() string {
	return fmt.Sprintf("{%v/%v}", m.Height, m.Round)
}
