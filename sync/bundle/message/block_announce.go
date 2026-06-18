package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
)

type BlockAnnounceMessage struct {
	Block       *block.Block             `cbor:"1,keyasint"`
	Certificate *certificate.Certificate `cbor:"2,keyasint"`
	Proof       *certificate.Certificate `cbor:"3,keyasint"`
}

func NewBlockAnnounceMessage(blk *block.Block,
	cert *certificate.Certificate, proof *certificate.Certificate,
) *BlockAnnounceMessage {
	return &BlockAnnounceMessage{
		Block:       blk,
		Certificate: cert,
		Proof:       proof,
	}
}

func (m *BlockAnnounceMessage) BasicCheck() error {
	if m.Block == nil {
		return BasicCheckError{Reason: "block is not set"}
	}
	if m.Certificate == nil {
		return BasicCheckError{Reason: "certificate is not set"}
	}

	if err := m.Block.BasicCheck(); err != nil {
		return err
	}

	if m.Proof != nil {
		if err := m.Proof.BasicCheck(); err != nil {
			return err
		}
	}

	return m.Certificate.BasicCheck()
}

func (m *BlockAnnounceMessage) Height() types.Height {
	if m.Certificate == nil {
		return 0
	}

	return m.Certificate.Height()
}

func (*BlockAnnounceMessage) Type() Type {
	return TypeBlockAnnounce
}

func (*BlockAnnounceMessage) TopicID() network.TopicID {
	return network.TopicIDBlock
}

func (*BlockAnnounceMessage) ShouldBroadcast() bool {
	return true
}

func (m *BlockAnnounceMessage) ConsensusHeight() types.Height {
	return m.Height()
}

// LogString returns a concise string representation intended for use in logs.
func (m *BlockAnnounceMessage) LogString() string {
	if m.Block != nil {
		return fmt.Sprintf("{⌘ %d nil}", m.Height())
	}

	return fmt.Sprintf("{⌘ %d %v}", m.Height(), m.Block.Hash().LogString())
}
