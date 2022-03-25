package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type QueryVotesMessage struct {
	Height int32 `cbor:"1,keyasint"`
	Round  int16 `cbor:"2,keyasint"`
}

func NewQueryVotesMessage(h int32, r int16) *QueryVotesMessage {
	return &QueryVotesMessage{
		Height: h,
		Round:  r,
	}
}

func (m *QueryVotesMessage) SanityCheck() error {
	if m.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if m.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Round")
	}
	return nil
}

func (m *QueryVotesMessage) Type() Type {
	return MessageTypeQueryVotes
}

func (m *QueryVotesMessage) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", m.Height, m.Round)
}
