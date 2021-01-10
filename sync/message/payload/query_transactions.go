package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type QueryTransactionsPayload struct {
	IDs []crypto.Hash `cbor:"1,keyasint"`
}

func (p *QueryTransactionsPayload) SanityCheck() error {
	if len(p.IDs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Empty list")
	}
	return nil
}

func (p *QueryTransactionsPayload) Type() PayloadType {
	return PayloadTypeQueryTransactions
}

func (p *QueryTransactionsPayload) Fingerprint() string {
	var s string
	for _, h := range p.IDs {
		s += fmt.Sprintf("%v ", h.Fingerprint())
	}
	return fmt.Sprintf("%v", s)
}
