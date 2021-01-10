package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type DownloadRequestPayload struct {
	SessionID int     `cbor:"1,keyasint"`
	Initiator peer.ID `cbor:"2,keyasint"`
	Target    peer.ID `cbor:"3,keyasint"`
	From      int     `cbor:"4,keyasint"`
	To        int     `cbor:"5,keyasint"`
}

func (p *DownloadRequestPayload) SanityCheck() error {
	if err := p.Initiator.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid initiator peer is: %v", err)
	}
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid target peer is: %v", err)
	}
	if p.From <= 0 {
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
	return fmt.Sprintf("{%v %v:%v}", util.FingerprintPeerID(p.Initiator), p.From, p.To)
}
