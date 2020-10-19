package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type TxResPayload struct {
	Tx tx.Tx `cbor:"1,keyasint"`
}

func NewTxResMessage(tx tx.Tx) Message {
	return Message{
		Type: PayloadTypeTxRes,
		Payload: &TxResPayload{
			Tx: tx,
		},
	}
}
func (p *TxResPayload) SanityCheck() error {
	if err := p.Tx.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid transaction")
	}
	return nil
}

func (p *TxResPayload) Type() PayloadType {
	return PayloadTypeTxRes
}

func (p *TxResPayload) Fingerprint() string {
	return fmt.Sprintf(" %v", p.Tx.Hash().Fingerprint())
}
