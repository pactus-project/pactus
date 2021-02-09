package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

type AleykPayload struct {
	ResponseCode    ResponseCode     `cbor:"1,keyasint"`
	ResponseMessage string           `cbor:"2,keyasint,omitempty"`
	NodeVersion     version.Version  `cbor:"3,keyasint"`
	Moniker         string           `cbor:"4,keyasint"`
	PublicKey       crypto.PublicKey `cbor:"5,keyasint"`
	PeerID          peer.ID          `cbor:"6,keyasint"`
	Height          int              `cbor:"7,keyasint"`
	Flags           int              `cbor:"8,keyasint"`
}

func (p *AleykPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return err
	}
	if err := p.PeerID.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *AleykPayload) Type() PayloadType {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%v %v}", util.FingerprintPeerID(p.PeerID), p.Height)
}
