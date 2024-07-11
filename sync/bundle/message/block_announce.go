package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
)

type BlockAnnounceMessage struct {
	Block       *block.Block                  `cbor:"1,keyasint"`
	Certificate *certificate.BlockCertificate `cbor:"2,keyasint"`
}

func NewBlockAnnounceMessage(blk *block.Block, cert *certificate.BlockCertificate) *BlockAnnounceMessage {
	return &BlockAnnounceMessage{
		Block:       blk,
		Certificate: cert,
	}
}

func (m *BlockAnnounceMessage) BasicCheck() error {
	if err := m.Block.BasicCheck(); err != nil {
		return err
	}

	return m.Certificate.BasicCheck()
}

func (m *BlockAnnounceMessage) Height() uint32 {
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

func (m *BlockAnnounceMessage) String() string {
	return fmt.Sprintf("{⌘ %d %v}",
		m.Certificate.Height(),
		m.Block.Hash().ShortString())
}
