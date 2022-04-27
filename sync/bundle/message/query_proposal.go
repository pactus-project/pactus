package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/util/errors"
)

type QueryProposalMessage struct {
	Height int32 `cbor:"1,keyasint"`
	Round  int16 `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(h int32, r int16) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height: h,
		Round:  r,
	}
}

func (m *QueryProposalMessage) SanityCheck() error {
	if m.Height < 0 {
		return errors.Error(errors.ErrInvalidHeight)
	}
	if m.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}

	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return MessageTypeQueryProposal
}

func (m *QueryProposalMessage) Fingerprint() string {
	return fmt.Sprintf("%v/%v", m.Height, m.Round)
}
