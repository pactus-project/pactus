package message

import (
	"fmt"

	"gitlab.com/zarb-chain/zarb-go/crypto"
)

type VoteSetPayload struct {
	Votes []crypto.Hash `cbor:"1,keyasint"`
}

func NewVoteSetMessage(height int, votes []crypto.Hash) *Message {
	return &Message{
		Type:   PayloadTypeVoteSet,
		Height: height,
		Payload: &VoteSetPayload{
			Votes: votes,
		},
	}
}

func (p *VoteSetPayload) SanityCheck() error {
	return nil
}

func (p *VoteSetPayload) Type() PayloadType {
	return PayloadTypeVoteSet
}

func (p *VoteSetPayload) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(p.Votes))
}
