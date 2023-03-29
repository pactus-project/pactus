package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/block"
)

type BlockAnnounceMessage struct {
	Height      uint32             `cbor:"1,keyasint"`
	Block       *block.Block       `cbor:"2,keyasint"`
	Certificate *block.Certificate `cbor:"3,keyasint"`
}

func NewBlockAnnounceMessage(h uint32, b *block.Block, c *block.Certificate) *BlockAnnounceMessage {
	return &BlockAnnounceMessage{
		Height:      h,
		Block:       b,
		Certificate: c,
	}
}

func (m *BlockAnnounceMessage) SanityCheck() error {
	if err := m.Block.SanityCheck(); err != nil {
		return err
	}
	return m.Certificate.SanityCheck()
}

func (m *BlockAnnounceMessage) Type() Type {
	return MessageTypeBlockAnnounce
}

func (m *BlockAnnounceMessage) Fingerprint() string {
	return fmt.Sprintf("{⌘ %d %v}", m.Height, m.Block.Hash().Fingerprint())
}
