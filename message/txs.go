package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
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
			return errors.Errorf(errors.ErrInvalidMessage, "invalid transaction")
		}
	}

	return nil
}

func (p *TxsPayload) Type() PayloadType {
	return PayloadTypeTxs
}

func (p *TxsPayload) Fingerprint() string {
	return fmt.Sprintf("%v", len(p.Txs))
}
