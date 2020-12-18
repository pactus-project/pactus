package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type TxsReqPayload struct {
	IDs []crypto.Hash `cbor:"1,keyasint"`
}

func NewTxsReqMessage(ids []crypto.Hash) *Message {
	return &Message{
		Type: PayloadTypeTxsReq,
		Payload: &TxsReqPayload{
			IDs: ids,
		},
	}
}

func (p *TxsReqPayload) SanityCheck() error {
	if len(p.IDs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Empty list")
	}
	return nil
}

func (p *TxsReqPayload) Type() PayloadType {
	return PayloadTypeTxsReq
}

func (p *TxsReqPayload) Fingerprint() string {
	var s string
	for _, h := range p.IDs {
		s += fmt.Sprintf("%v ", h.Fingerprint())
	}
	return fmt.Sprintf("%v", s)
}
