package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type AleykPayload struct {
	Agent           string         `cbor:"1,keyasint"`
	Moniker         string         `cbor:"2,keyasint"`
	PublicKey       *bls.PublicKey `cbor:"3,keyasint"`
	Signature       *bls.Signature `cbor:"4,keyasint"`
	Height          int            `cbor:"5,keyasint"`
	Flags           int            `cbor:"6,keyasint"`
	ResponseTarget  peer.ID        `cbor:"7,keyasint"`
	ResponseCode    ResponseCode   `cbor:"8,keyasint"`
	ResponseMessage string         `cbor:"9,keyasint"`
}

func NewAleykPayload(moniker string, pub crypto.PublicKey, sig crypto.Signature,
	height int, flags int, target peer.ID, code ResponseCode, msg string) Payload {
	return &AleykPayload{
		ResponseTarget:  target,
		ResponseCode:    code,
		ResponseMessage: msg,
		Agent:           version.Agent(),
		Moniker:         moniker,
		PublicKey:       pub.(*bls.PublicKey),
		Signature:       sig.(*bls.Signature),
		Height:          height,
		Flags:           flags,
	}
}

func (p *AleykPayload) SanityCheck() error {
	if err := p.ResponseTarget.Validate(); err != nil {
		return err
	}
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if !p.PublicKey.Verify(p.PublicKey.RawBytes(), p.Signature) {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid signature")
	}
	return nil
}

func (p *AleykPayload) Type() Type {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
