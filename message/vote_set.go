package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type VoteSetPayload struct {
	Height int
	Hashes []crypto.Hash `cbor:"1,keyasint"`
}

func NewVoteSetMessage(height int, Hashes []crypto.Hash) Message {
	return Message{
		Type: PayloadTypeVoteSet,
		Payload: &VoteSetPayload{
			Height: height,
			Hashes: Hashes,
		},
	}
}

func (p *VoteSetPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *VoteSetPayload) Type() PayloadType {
	return PayloadTypeVoteSet
}

func (p *VoteSetPayload) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(p.Hashes))
}
