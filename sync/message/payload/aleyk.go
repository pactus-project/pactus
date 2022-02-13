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
	PeerID          peer.ID        `cbor:"1,keyasint"`
	Agent           string         `cbor:"2,keyasint"`
	Moniker         string         `cbor:"3,keyasint"`
	PublicKey       *bls.PublicKey `cbor:"4,keyasint"`
	Signature       *bls.Signature `cbor:"5,keyasint"`
	Height          int            `cbor:"6,keyasint"`
	Flags           int            `cbor:"7,keyasint"`
	ResponseTarget  peer.ID        `cbor:"8,keyasint"`
	ResponseCode    ResponseCode   `cbor:"9,keyasint"`
	ResponseMessage string         `cbor:"10,keyasint"`
}

func NewAleykPayload(pid peer.ID, moniker string,
	height int, flags int, target peer.ID, code ResponseCode, msg string) *AleykPayload {
	return &AleykPayload{
		PeerID:          pid,
		ResponseTarget:  target,
		ResponseCode:    code,
		ResponseMessage: msg,
		Agent:           version.Agent(),
		Moniker:         moniker,
		Height:          height,
		Flags:           flags,
	}
}

func (p *AleykPayload) SignBytes() []byte {
	return []byte(fmt.Sprintf("%d:%s:%s", p.Type(), p.Agent, p.PeerID))
}

func (p *AleykPayload) SetSignature(sig crypto.Signature) {
	p.Signature = sig.(*bls.Signature)
}

func (p *AleykPayload) SetPublicKey(pub crypto.PublicKey) {
	p.PublicKey = pub.(*bls.PublicKey)
}

func (p *AleykPayload) SanityCheck() error {
	if err := p.ResponseTarget.Validate(); err != nil {
		return err
	}
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if !p.PublicKey.Verify(p.SignBytes(), p.Signature) {
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
