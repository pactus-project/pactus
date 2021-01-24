package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/errors"
)

type QueryVotesPayload struct {
	Querier peer.ID `cbor:"1,keyasint"`
	Height  int     `cbor:"2,keyasint"`
	Round   int     `cbor:"3,keyasint"`
}

func (p *QueryVotesPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Round")
	}
	if err := p.Querier.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid querier peer id: %v", err)
	}
	return nil
}

func (p *QueryVotesPayload) Type() PayloadType {
	return PayloadTypeQueryVotes
}

func (p *QueryVotesPayload) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", p.Height, p.Round)
}
