package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/tx"
)

type TxsPayload struct {
	Txs []tx.Tx `cbor:"2,keyasint"`
}

func NewTxsMessage(txs []tx.Tx) *Message {
	return &Message{
		Type: PayloadTypeTxs,
		Payload: &TxsPayload{
			Txs: txs,
		},
	}
}
func (p *TxsPayload) SanityCheck() error {
	for _, tx := range p.Txs {
		if err := tx.SanityCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (p *TxsPayload) Type() PayloadType {
	return PayloadTypeTxs
}

func (p *TxsPayload) Fingerprint() string {
	var s string
	for _, tx := range p.Txs {
		s += fmt.Sprintf("%v ", tx.ID().Fingerprint())
	}
	return fmt.Sprintf("{%v: ⌘ [%v]}", len(p.Txs), s)
}
