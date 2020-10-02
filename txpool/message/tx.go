package message

import (
	"gitlab.com/zarb-chain/zarb-go/errors"
	"gitlab.com/zarb-chain/zarb-go/tx"
)

type TxPayload struct {
	TxData []byte `cbor:"1,keyasint"`
}

func NewTxMessage(tx *tx.Tx) *Message {
	txData, _ := tx.Encode()
	return &Message{
		Type: PayloadTypeTx,
		Payload: &TxPayload{
			TxData: txData,
		},
	}
}
func (p *TxPayload) SanityCheck() error {
	if len(p.TxData) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "empty transaction")
	}
	return nil
}

func (p *TxPayload) Type() PayloadType {
	return PayloadTypeTx
}

func (p *TxPayload) Fingerprint() string {
	return ""
}
