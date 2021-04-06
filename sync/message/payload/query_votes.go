package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type QueryVotesPayload struct {
	Height int `cbor:"1,keyasint"`
	Round  int `cbor:"2,keyasint"`
}

func NewQueryVotesPAyload(h, r int) Payload {
	return &QueryVotesPayload{
		Height: h,
		Round:  r,
	}
}

func (p *QueryVotesPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Round")
	}
	return nil
}

func (p *QueryVotesPayload) Type() PayloadType {
	return PayloadTypeQueryVotes
}

func (p *QueryVotesPayload) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", p.Height, p.Round)
}
