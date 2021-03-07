package payload

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/proposal"
)

type ProposalPayload struct {
	Proposal *proposal.Proposal `cbor:"1,keyasint"`
}

func (p *ProposalPayload) SanityCheck() error {
	if p.Proposal == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "No proposal")
	}
	if err := p.Proposal.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return nil
}

func (p *ProposalPayload) Type() PayloadType {
	return PayloadTypeProposal
}

func (p *ProposalPayload) Fingerprint() string {
	return p.Proposal.Fingerprint()
}
