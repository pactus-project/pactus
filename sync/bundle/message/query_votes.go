package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/errors"
)

type QueryVotesMessage struct {
	Height  uint32         `cbor:"1,keyasint"`
	Round   int16          `cbor:"2,keyasint"`
	Querier crypto.Address `cbor:"3,keyasint"`
}

func NewQueryVotesMessage(height uint32, round int16, querier crypto.Address) *QueryVotesMessage {
	return &QueryVotesMessage{
		Height:  height,
		Round:   round,
		Querier: querier,
	}
}

func (m *QueryVotesMessage) BasicCheck() error {
	if m.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}
	return nil
}

func (m *QueryVotesMessage) Type() Type {
	return TypeQueryVotes
}

func (m *QueryVotesMessage) String() string {
	return fmt.Sprintf("{%d/%d %s}", m.Height, m.Round, m.Querier.ShortString())
}
