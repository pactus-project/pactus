package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type AleykPayload struct {
	ResponseTarget  peer.ID          `cbor:"1,keyasint"`
	ResponseCode    ResponseCode     `cbor:"2,keyasint"`
	ResponseMessage string           `cbor:"3,keyasint"`
	NodeVersion     string           `cbor:"4,keyasint"`
	Moniker         string           `cbor:"5,keyasint"`
	PublicKey       crypto.PublicKey `cbor:"6,keyasint"`
	Height          int              `cbor:"7,keyasint"`
	Flags           int              `cbor:"8,keyasint"`
}

func NewAleykPayload(target peer.ID, code ResponseCode, msg string, moniker string,
	pub crypto.PublicKey, height int, flags int) Payload {
	return &AleykPayload{
		ResponseTarget:  target,
		ResponseCode:    code,
		ResponseMessage: msg,
		NodeVersion:     version.Version(),
		Moniker:         moniker,
		PublicKey:       pub,
		Height:          height,
		Flags:           flags,
	}
}

func (p *AleykPayload) SanityCheck() error {
	if err := p.ResponseTarget.Validate(); err != nil {
		return err
	}
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return err
	}
	return nil
}

func (p *AleykPayload) Type() Type {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
