package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/errors"
)

type LatestBlocksRequestPayload struct {
	RequestID int     `cbor:"1,keyasint"`
	Initiator peer.ID `cbor:"2,keyasint"`
	Target    peer.ID `cbor:"3,keyasint"`
	From      int     `cbor:"4,keyasint"`
}

func (p *LatestBlocksRequestPayload) SanityCheck() error {
	if err := p.Initiator.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid initiator peer is: %v", err)
	}
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid target peer is: %v", err)
	}
	if p.From <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *LatestBlocksRequestPayload) Type() PayloadType {
	return PayloadTypeLatestBlocksRequest
}

func (p *LatestBlocksRequestPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.From)
}
