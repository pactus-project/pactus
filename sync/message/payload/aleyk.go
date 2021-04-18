package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type AleykPayload struct {
	ResponseCode    ResponseCode     `cbor:"1,keyasint"`
	ResponseMessage string           `cbor:"2,keyasint"`
	NodeVersion     version.Version  `cbor:"3,keyasint"`
	Moniker         string           `cbor:"4,keyasint"`
	PublicKey       crypto.PublicKey `cbor:"5,keyasint"`
	Height          int              `cbor:"6,keyasint"`
	Flags           int              `cbor:"7,keyasint"`
}

func NewAleykPayload(code ResponseCode, msg string, moniker string,
	pub crypto.PublicKey, height int, flags int) Payload {
	return &AleykPayload{
		ResponseCode:    code,
		ResponseMessage: msg,
		NodeVersion:     version.NodeVersion,
		Moniker:         moniker,
		PublicKey:       pub,
		Height:          height,
		Flags:           flags,
	}
}

func (p *AleykPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return err
	}
	return nil
}

func (p *AleykPayload) Type() PayloadType {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
