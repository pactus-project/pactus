package payload

import (
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/errors"
)

type VotePayload struct {
	Vote *vote.Vote `cbor:"1,keyasint"`
}

func NewVotePayload(v *vote.Vote) Payload {
	return &VotePayload{
		Vote: v,
	}
}

func (p *VotePayload) SanityCheck() error {
	if err := p.Vote.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}
	return nil
}

func (p *VotePayload) Type() PayloadType {
	return PayloadTypeVote
}

func (p *VotePayload) Fingerprint() string {
	return p.Vote.Fingerprint()
}
