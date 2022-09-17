package message

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/version"
)

const (
	FlagHelloAck    = 0x0001
	FlagNodeNetwork = 0x0100
)

type HelloMessage struct {
	PeerID      peer.ID        `cbor:"1,keyasint"`
	Agent       string         `cbor:"2,keyasint"`
	Moniker     string         `cbor:"3,keyasint"`
	PublicKey   *bls.PublicKey `cbor:"4,keyasint"`
	Signature   *bls.Signature `cbor:"5,keyasint"`
	Height      uint32         `cbor:"6,keyasint"`
	Flags       int            `cbor:"7,keyasint"`
	GenesisHash hash.Hash      `cbor:"8,keyasint"`
}

func NewHelloMessage(pid peer.ID, moniker string,
	height uint32, flags int, genesisHash hash.Hash) *HelloMessage {
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
	if m.Signature == nil {
		return errors.Error(errors.ErrInvalidSignature)
	}
	if m.PublicKey == nil {
		return errors.Error(errors.ErrInvalidPublicKey)
	}
	if err := m.PublicKey.Verify(m.SignBytes(), m.Signature); err != nil {
		return err
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
	ack := ""
	if util.IsFlagSet(m.Flags, FlagHelloAck) {
		ack = " ack"
	}
	return fmt.Sprintf("{%s %v%s}", m.Moniker, m.Height, ack)
}
