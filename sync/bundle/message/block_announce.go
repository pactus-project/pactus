package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
)

type BlockAnnounceMessage struct {
	Height      uint32                   `cbor:"1,keyasint"`
	Block       *block.Block             `cbor:"2,keyasint"`
	Certificate *certificate.Certificate `cbor:"3,keyasint"`
}

func NewBlockAnnounceMessage(h uint32, b *block.Block, c *certificate.Certificate) *BlockAnnounceMessage {
	return &BlockAnnounceMessage{
		Height:      h,
		Block:       b,
		Certificate: c,
	}
}

func (m *BlockAnnounceMessage) BasicCheck() error {
	if err := m.Block.BasicCheck(); err != nil {
		return err
	}
	return m.Certificate.BasicCheck()
}

func (m *BlockAnnounceMessage) Type() Type {
	return TypeBlockAnnounce
}

func (m *BlockAnnounceMessage) String() string {
	return fmt.Sprintf("{⌘ %d %v}", m.Height, m.Block.Hash().ShortString())
}
