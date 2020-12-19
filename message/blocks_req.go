package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type BlocksReqPayload struct {
	From          int         `cbor:"1,keyasint"`
	To            int         `cbor:"2,keyasint"`
	LastBlockHash crypto.Hash `cbor:"3,keyasint"`
}

func NewBlocksReqMessage(from, to int, lastBlockHash crypto.Hash) *Message {
	return &Message{
		Type: PayloadTypeBlocksReq,
		Payload: &BlocksReqPayload{
			From:          from,
			To:            to,
			LastBlockHash: lastBlockHash,
		},
	}

}
func (p *BlocksReqPayload) SanityCheck() error {
	if p.From <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.To <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.To < p.From {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *BlocksReqPayload) Type() PayloadType {
	return PayloadTypeBlocksReq
}

func (p *BlocksReqPayload) Fingerprint() string {
	return fmt.Sprintf("{%v-%v ⌘ %v}", p.From, p.To, p.LastBlockHash.Fingerprint())
}
