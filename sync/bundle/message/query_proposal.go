package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type QueryProposalMessage struct {
	Height int `cbor:"1,keyasint"`
	Round  int `cbor:"2,keyasint"`
}

func NewQueryProposalMessage(h, r int) *QueryProposalMessage {
	return &QueryProposalMessage{
		Height: h,
		Round:  r,
	}
}

func (m *QueryProposalMessage) SanityCheck() error {
	if m.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if m.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid round")
	}

	return nil
}

func (m *QueryProposalMessage) Type() Type {
	return MessageTypeQueryProposal
}

func (m *QueryProposalMessage) Fingerprint() string {
	return fmt.Sprintf("%v/%v", m.Height, m.Round)
}
