package message

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/version"
)

type HelloMessage struct {
	PeerID          peer.ID          `cbor:"1,keyasint"`
	Agent           string           `cbor:"2,keyasint"`
	Moniker         string           `cbor:"3,keyasint"`
	PublicKeys      []*bls.PublicKey `cbor:"4,keyasint"`
	Signature       *bls.Signature   `cbor:"5,keyasint"`
	Height          uint32           `cbor:"6,keyasint"`
	Services        service.Services `cbor:"7,keyasint"`
	GenesisHash     hash.Hash        `cbor:"8,keyasint"`
	BlockHash       hash.Hash        `cbor:"9,keyasint"`
	MyTimeUnixMilli int64            `cbor:"10,keyasint"`
}

func NewHelloMessage(pid peer.ID, moniker string,
	height uint32, services service.Services, blockHash, genesisHash hash.Hash,
) *HelloMessage {
	return &HelloMessage{
		PeerID:          pid,
		Agent:           version.Agent(),
		Moniker:         moniker,
		GenesisHash:     genesisHash,
		BlockHash:       blockHash,
		Height:          height,
		Services:        services,
		MyTimeUnixMilli: time.Now().UnixMilli(),
	}
}

func (m *HelloMessage) BasicCheck() error {
	if m.Signature == nil {
		return errors.Error(errors.ErrInvalidSignature)
	}
	if len(m.PublicKeys) == 0 {
		return errors.Error(errors.ErrInvalidPublicKey)
	}
	aggPublicKey := bls.PublicKeyAggregate(m.PublicKeys...)

	return aggPublicKey.Verify(m.SignBytes(), m.Signature)
}

func (m *HelloMessage) MyTime() time.Time {
	return time.UnixMilli(m.MyTimeUnixMilli)
}

func (m *HelloMessage) SignBytes() []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s", m.Type(), m.Agent, m.PeerID, m.GenesisHash.String()))
}

func (m *HelloMessage) Type() Type {
	return TypeHello
}

func (m *HelloMessage) String() string {
	return fmt.Sprintf("{%s %d %s}", m.Moniker, m.Height, m.Services)
}

func (m *HelloMessage) Sign(valKeys []*bls.ValidatorKey) {
	signatures := make([]*bls.Signature, len(valKeys))
	publicKeys := make([]*bls.PublicKey, len(valKeys))
	signBytes := m.SignBytes()
	for i, key := range valKeys {
		signatures[i] = key.Sign(signBytes)
		publicKeys[i] = key.PublicKey()
	}
	aggSignature := bls.SignatureAggregate(signatures...)
	m.Signature = aggSignature
	m.PublicKeys = publicKeys
}
