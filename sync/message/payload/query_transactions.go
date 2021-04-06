package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type QueryTransactionsPayload struct {
	IDs []tx.ID `cbor:"1,keyasint"`
}

func NewQueryTransactionsPayload(ids []tx.ID) Payload {
	return &QueryTransactionsPayload{
		IDs: ids,
	}
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
	return fmt.Sprintf("{%v}", s)
}
