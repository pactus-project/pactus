package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type QueryTransactionsMessage struct {
	IDs []tx.ID `cbor:"1,keyasint"`
}

func NewQueryTransactionsMessage(ids []tx.ID) *QueryTransactionsMessage {
	return &QueryTransactionsMessage{
		IDs: ids,
	}
}

func (m *QueryTransactionsMessage) SanityCheck() error {
	if len(m.IDs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "empty list")
	}
	return nil
}

func (m *QueryTransactionsMessage) Type() Type {
	return MessageTypeQueryTransactions
}

func (m *QueryTransactionsMessage) Fingerprint() string {
	var s string
	for _, h := range m.IDs {
		s += fmt.Sprintf("%v ", h.Fingerprint())
	}
	return fmt.Sprintf("{%v}", s)
}
