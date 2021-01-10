package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type VoteSetPayload struct {
	Height int           `cbor:"1,keyasint"`
	Round  int           `cbor:"2,keyasint"`
	Hashes []crypto.Hash `cbor:"3,keyasint"`
}

func (p *VoteSetPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *VoteSetPayload) Type() PayloadType {
	return PayloadTypeVoteSet
}

func (p *VoteSetPayload) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", p.Height, p.Round)
}
