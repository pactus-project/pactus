package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type BlocksReqPayload struct {
	From int `cbor:"1,keyasint"`
	To   int `cbor:"1,keyasint"`
}

func NewBlocksReqMessage(from, to int) Message {
	return Message{
		Type: PayloadTypeBlocksReq,
		Payload: &BlocksReqPayload{
			From: from,
			To:   to,
		},
	}

}
func (p *BlocksReqPayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.To < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *BlocksReqPayload) Type() PayloadType {
	return PayloadTypeBlocksReq
}

func (p *BlocksReqPayload) Fingerprint() string {
	return fmt.Sprintf("{%v-%v}", p.From, p.To)
}
