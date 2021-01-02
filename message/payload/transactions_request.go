package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type TransactionsRequestPayload struct {
	IDs []crypto.Hash `cbor:"1,keyasint"`
}

func (p *TransactionsRequestPayload) SanityCheck() error {
	if len(p.IDs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Empty list")
	}
	return nil
}

func (p *TransactionsRequestPayload) Type() PayloadType {
	return PayloadTypeTransactionsRequest
}

func (p *TransactionsRequestPayload) Fingerprint() string {
	var s string
	for _, h := range p.IDs {
		s += fmt.Sprintf("%v ", h.Fingerprint())
	}
	return fmt.Sprintf("%v", s)
}
