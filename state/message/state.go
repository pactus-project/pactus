package message

import (
	"fmt"

	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/errors"
)

type StateInfoPayload struct {
	Height int         `cbor:"1,keyasint,omitempty"`
	Hash   crypto.Hash `cbor:"2,keyasint"`
}

func NewStateInfoMessage(height int, hash crypto.Hash) *Message {
	return &Message{
		Type: PayloadTypeStateInfo,
		Payload: &StateInfoPayload{
			Height: height,
			Hash:   hash,
		},
	}

}
func (p *StateInfoPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.Hash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid hash: %v", err)
	}
	return nil
}

func (p *StateInfoPayload) Type() PayloadType {
	return PayloadTypeStateInfo
}

func (p *StateInfoPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.Hash.Fingerprint())
}
