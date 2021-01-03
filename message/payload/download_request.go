package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type DownloadRequestPayload struct {
	PeerID peer.ID `cbor:"1,keyasint"`
	From   int     `cbor:"2,keyasint"`
	To     int     `cbor:"3,keyasint"`
}

func (p *DownloadRequestPayload) SanityCheck() error {
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
	return fmt.Sprintf("{%v ☍ %v}", p.From, util.FingerprintPeerID(p.PeerID))
}
