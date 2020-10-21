package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type TxReqPayload struct {
	Hash crypto.Hash `cbor:"1,keyasint"`
}

func NewTxReqMessage(hash crypto.Hash) Message {
	return Message{
		Type: PayloadTypeTxReq,
		Payload: &TxReqPayload{
			Hash: hash,
		},
	}
}

func (p *TxReqPayload) SanityCheck() error {
	if err := p.Hash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid hash")
	}
	return nil
}

func (p *TxReqPayload) Type() PayloadType {
	return PayloadTypeTxReq
}

func (p *TxReqPayload) Fingerprint() string {
	return fmt.Sprintf(" %v", p.Hash.Fingerprint())
}
