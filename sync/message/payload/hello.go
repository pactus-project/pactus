package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

const (
	FlagNeedResponse         = 0x1
	FlagInitialBlockDownload = 0x2
)

type HelloPayload struct {
	PeerID      peer.ID        `cbor:"1,keyasint"`
	Agent       string         `cbor:"2,keyasint"`
	Moniker     string         `cbor:"3,keyasint"`
	PublicKey   *bls.PublicKey `cbor:"4,keyasint"`
	Signature   *bls.Signature `cbor:"5,keyasint"`
	Height      int            `cbor:"6,keyasint"`
	Flags       int            `cbor:"7,keyasint"`
	GenesisHash hash.Hash      `cbor:"8,keyasint"`
}

func NewHelloPayload(pid peer.ID, moniker string,
	height int, flags int, genesisHash hash.Hash) *HelloPayload {
	return &HelloPayload{
		PeerID:      pid,
		Agent:       version.Agent(),
		Moniker:     moniker,
		GenesisHash: genesisHash,
		Height:      height,
		Flags:       flags,
	}
}

func (p *HelloPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if !p.PublicKey.Verify(p.SignBytes(), p.Signature) {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid signature")
	}
	return nil
}

func (p *HelloPayload) SignBytes() []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", p.Type(), p.Agent, p.PeerID))
}

func (p *HelloPayload) SetSignature(sig crypto.Signature) {
	p.Signature = sig.(*bls.Signature)
}

func (p *HelloPayload) SetPublicKey(pub crypto.PublicKey) {
	p.PublicKey = pub.(*bls.PublicKey)
}

func (p *HelloPayload) Type() Type {
	return PayloadTypeHello
}

func (p *HelloPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
