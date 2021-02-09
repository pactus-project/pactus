package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
)

type QueryProposalPayload struct {
	Querier peer.ID `cbor:"1,keyasint"`
	Height  int     `cbor:"2,keyasint"`
	Round   int     `cbor:"3,keyasint"`
}

func (p *QueryProposalPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid round")
	}
	if err := p.Querier.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid querier peer id: %v", err)
	}

	return nil
}

func (p *QueryProposalPayload) Type() PayloadType {
	return PayloadTypeQueryProposal
}

func (p *QueryProposalPayload) Fingerprint() string {
	return fmt.Sprintf("%v/%v", p.Height, p.Round)
}
