package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type TxsReqPayload struct {
	Hashes []crypto.Hash `cbor:"1,keyasint"`
}

func NewTxsReqMessage(hashes []crypto.Hash) Message {
	return Message{
		Type: PayloadTypeTxsReq,
		Payload: &TxsReqPayload{
			Hashes: hashes,
		},
	}
}

func (p *TxsReqPayload) SanityCheck() error {
	if len(p.Hashes) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Empty list")
	}
	return nil
}

func (p *TxsReqPayload) Type() PayloadType {
	return PayloadTypeTxsReq
}

func (p *TxsReqPayload) Fingerprint() string {
	return fmt.Sprintf("%v", len(p.Hashes))
}
