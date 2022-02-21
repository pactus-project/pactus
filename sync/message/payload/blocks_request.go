package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type BlocksRequestPayload struct {
	SessionID int `cbor:"1,keyasint"`
	From      int `cbor:"2,keyasint"`
	To        int `cbor:"3,keyasint"`
}

func NewBlocksRequestPayload(sid int, from, to int) *BlocksRequestPayload {
	return &BlocksRequestPayload{
		SessionID: sid,
		From:      from,
		To:        to,
	}
}

func (p *BlocksRequestPayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.From > p.To {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid range")
	}
	return nil
}

func (p *BlocksRequestPayload) Type() Type {
	return PayloadTypeBlocksRequest
}

func (p *BlocksRequestPayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %v:%v}", p.SessionID, p.From, p.To)
}
