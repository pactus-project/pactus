package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type QueryProposalPayload struct {
	Height int `cbor:"1,keyasint"`
	Round  int `cbor:"2,keyasint"`
}

func (p *QueryProposalPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid round")
	}

	return nil
}

func (p *QueryProposalPayload) Type() PayloadType {
	return PayloadTypeQueryProposal
}

func (p *QueryProposalPayload) Fingerprint() string {
	return fmt.Sprintf("%v/%v", p.Height, p.Round)
}
