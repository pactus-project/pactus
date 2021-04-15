package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type TransactionsPayload struct {
	Transactions []*tx.Tx `cbor:"1,keyasint"`
}

func NewTransactionsPayload(trxs []*tx.Tx) Payload {
	return &TransactionsPayload{
		Transactions: trxs,
	}
}

func (p *TransactionsPayload) SanityCheck() error {
	if len(p.Transactions) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "no transaction")
	}
	for _, tx := range p.Transactions {
		if err := tx.SanityCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (p *TransactionsPayload) Type() PayloadType {
	return PayloadTypeTransactions
}

func (p *TransactionsPayload) Fingerprint() string {
	var s string
	for _, tx := range p.Transactions {
		s += fmt.Sprintf("%v ", tx.ID().Fingerprint())
	}
	return fmt.Sprintf("{%v: ⌘ [%v]}", len(p.Transactions), s)
}
