package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
)

type TransactionsMessage struct {
	Transactions []*tx.Tx `cbor:"1,keyasint"`
}

func NewTransactionsMessage(trxs []*tx.Tx) *TransactionsMessage {
	return &TransactionsMessage{
		Transactions: trxs,
	}
}

func (m *TransactionsMessage) SanityCheck() error {
	if len(m.Transactions) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "no transaction")
	}
	for _, tx := range m.Transactions {
		if err := tx.SanityCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (m *TransactionsMessage) Type() Type {
	return MessageTypeTransactions
}

func (m *TransactionsMessage) Fingerprint() string {
	var s string
	for _, tx := range m.Transactions {
		s += fmt.Sprintf("%v ", tx.ID().Fingerprint())
	}
	return fmt.Sprintf("{%v: ⌘ [%v]}", len(m.Transactions), s)
}
