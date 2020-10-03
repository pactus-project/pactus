package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type RequestPayload struct {
	Hash crypto.Hash `cbor:"1,keyasint"`
}

func NewRequestMessage(hash crypto.Hash) *Message {
	return &Message{
		Type: PayloadTypeRequest,
		Payload: &RequestPayload{
			Hash: hash,
		},
	}
}

func (p *RequestPayload) SanityCheck() error {
	if err := p.Hash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid hash")
	}
	return nil
}

func (p *RequestPayload) Type() PayloadType {
	return PayloadTypeRequest
}

func (p *RequestPayload) Fingerprint() string {
	return fmt.Sprintf(" %v", p.Hash.Fingerprint())
}
