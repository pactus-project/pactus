package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type DownloadRequestPayload struct {
	SessionID int     `cbor:"1,keyasint"`
	Target    peer.ID `cbor:"2,keyasint"`
	From      int     `cbor:"3,keyasint"`
	To        int     `cbor:"4,keyasint"`
}

func NewDownloadRequestPayload(sid int, target peer.ID, from, to int) Payload {
	return &DownloadRequestPayload{
		SessionID: sid,
		Target:    target,
		From:      from,
		To:        to,
	}
}

func (p *DownloadRequestPayload) SanityCheck() error {
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid target peer id: %v", err)
	}
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if p.From > p.To {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid range")
	}
	return nil
}

func (p *DownloadRequestPayload) Type() PayloadType {
	return PayloadTypeDownloadRequest
}

func (p *DownloadRequestPayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %v %v:%v}", p.SessionID, util.FingerprintPeerID(p.Target), p.From, p.To)
}
