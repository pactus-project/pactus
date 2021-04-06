package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
)

type LatestBlocksRequestPayload struct {
	SessionID int     `cbor:"1,keyasint"`
	Target    peer.ID `cbor:"2,keyasint"`
	From      int     `cbor:"3,keyasint"`
	To        int     `cbor:"4,keyasint"`
}

func NewLatestBlocksRequestPayload(sid int, target peer.ID, from, to int) Payload {
	return &LatestBlocksRequestPayload{
		SessionID: sid,
		Target:    target,
		From:      from,
		To:        to,
	}
}

func (p *LatestBlocksRequestPayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if p.From > p.To {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid range")
	}
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid target peer id: %v", err)
	}
	return nil
}

func (p *LatestBlocksRequestPayload) Type() PayloadType {
	return PayloadTypeLatestBlocksRequest
}

func (p *LatestBlocksRequestPayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %v}", p.SessionID, p.From)
}
