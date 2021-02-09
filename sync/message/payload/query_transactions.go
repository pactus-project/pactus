package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type QueryTransactionsPayload struct {
	Querier peer.ID       `cbor:"1,keyasint"`
	IDs     []crypto.Hash `cbor:"2,keyasint"`
}

func (p *QueryTransactionsPayload) SanityCheck() error {
	if len(p.IDs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Empty list")
	}
	if err := p.Querier.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid querier peer id: %v", err)
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
