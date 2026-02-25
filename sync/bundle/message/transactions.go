package message

import (
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types/tx"
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
		return BasicCheckError{Reason: "no transaction"}
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

func (*TransactionsMessage) TopicID() network.TopicID {
	return network.TopicIDTransaction
}

func (*TransactionsMessage) ShouldBroadcast() bool {
	return true
}

func (*TransactionsMessage) ConsensusHeight() uint32 {
	return 0
}

// LogString returns a concise string representation intended for use in logs.
func (m *TransactionsMessage) LogString() string {
	var builder strings.Builder

	for _, trx := range m.Transactions {
		fmt.Fprintf(&builder, "%v ", trx.ID().LogString())
	}
	fmt.Fprintf(&builder, "{%v: ⌘ [%v]}", len(m.Transactions), builder.String())

	return builder.String()
}
