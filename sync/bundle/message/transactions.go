package message

import (
	"fmt"
	"strings"

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

func (m *TransactionsMessage) BasicCheck() error {
	if len(m.Transactions) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "no transaction")
	}
	for _, trx := range m.Transactions {
		if err := trx.BasicCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (*TransactionsMessage) Type() Type {
	return TypeTransaction
}

func (m *TransactionsMessage) String() string {
	var builder strings.Builder

	for _, trx := range m.Transactions {
		builder.WriteString(fmt.Sprintf("%v ", trx.ID().ShortString()))
	}
	builder.WriteString(fmt.Sprintf("{%v: ⌘ [%v]}", len(m.Transactions), builder.String()))

	return builder.String()
}
