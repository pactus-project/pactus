package message

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
	FlagNodeNetwork  = 0x0001
	FlagNeedResponse = 0x1000
)

type HelloMessage struct {
	PeerID      peer.ID        `cbor:"1,keyasint"`
	Agent       string         `cbor:"2,keyasint"`
	Moniker     string         `cbor:"3,keyasint"`
	PublicKey   *bls.PublicKey `cbor:"4,keyasint"`
	Signature   *bls.Signature `cbor:"5,keyasint"`
	Height      int            `cbor:"6,keyasint"`
	Flags       int            `cbor:"7,keyasint"`
	GenesisHash hash.Hash      `cbor:"8,keyasint"`
}

func NewHelloMessage(pid peer.ID, moniker string,
	height int, flags int, genesisHash hash.Hash) *HelloMessage {
	return &HelloMessage{
		PeerID:      pid,
		Agent:       version.Agent(),
		Moniker:     moniker,
		GenesisHash: genesisHash,
		Height:      height,
		Flags:       flags,
	}
}

func (m *HelloMessage) SanityCheck() error {
	if m.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if !m.PublicKey.Verify(m.SignBytes(), m.Signature) {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid signature")
	}
	return nil
}

func (m *HelloMessage) SignBytes() []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", m.Type(), m.Agent, m.PeerID))
}

func (m *HelloMessage) SetSignature(sig crypto.Signature) {
	m.Signature = sig.(*bls.Signature)
}

func (m *HelloMessage) SetPublicKey(pub crypto.PublicKey) {
	m.PublicKey = pub.(*bls.PublicKey)
}

func (m *HelloMessage) Type() Type {
	return MessageTypeHello
}

func (m *HelloMessage) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", m.Moniker, m.Height)
}
