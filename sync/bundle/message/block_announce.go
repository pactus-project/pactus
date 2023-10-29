package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
)

type BlockAnnounceMessage struct {
	Block       *block.Block             `cbor:"1,keyasint"`
	Certificate *certificate.Certificate `cbor:"2,keyasint"`
}

func NewBlockAnnounceMessage(blk *block.Block, cert *certificate.Certificate) *BlockAnnounceMessage {
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

func (m *BlockAnnounceMessage) Type() Type {
	return TypeBlockAnnounce
}

func (m *BlockAnnounceMessage) String() string {
	return fmt.Sprintf("{⌘ %d %v}",
		m.Certificate.Height(),
		m.Block.Hash().ShortString())
}
