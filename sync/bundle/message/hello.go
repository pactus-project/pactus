package message

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
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
	services service.Services, height uint32, blockHash, genesisHash hash.Hash,
) *HelloMessage {
	return &HelloMessage{
		PeerID:          pid,
		Agent:           version.NodeAgent.String(),
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
		return BasicCheckError{"no signature"}
	}
	if len(m.PublicKeys) == 0 {
		return BasicCheckError{"no public key"}
	}
	aggPub, err := bls.PublicKeyAggregate(m.PublicKeys...)
	if err != nil {
		return BasicCheckError{err.Error()}
	}

	return aggPub.Verify(m.SignBytes(), m.Signature)
}

func (m *HelloMessage) MyTime() time.Time {
	return time.UnixMilli(m.MyTimeUnixMilli)
}

func (m *HelloMessage) SignBytes() []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s", m.Type(), m.Agent, m.PeerID, m.GenesisHash.String()))
}

func (*HelloMessage) Type() Type {
	return TypeHello
}

func (*HelloMessage) TopicID() network.TopicID {
	return network.TopicIDUnspecified
}

func (*HelloMessage) ShouldBroadcast() bool {
	return false
}

func (*HelloMessage) ConsensusHeight() uint32 {
	return 0
}

// LogString returns a concise string representation intended for use in logs.
func (m *HelloMessage) LogString() string {
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
	aggSig, _ := bls.SignatureAggregate(signatures...)
	m.Signature = aggSig
	m.PublicKeys = publicKeys
}
