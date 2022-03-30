package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlockAnnounceMessage struct {
	Height      int32              `cbor:"1,keyasint"`
	Block       *block.Block       `cbor:"2,keyasint"`
	Certificate *block.Certificate `cbor:"3,keyasint"`
}

func NewBlockAnnounceMessage(h int32, b *block.Block, c *block.Certificate) *BlockAnnounceMessage {
	return &BlockAnnounceMessage{
		Height:      h,
		Block:       b,
		Certificate: c,
	}
}

func (m *BlockAnnounceMessage) SanityCheck() error {
	if m.Height < 0 {
		return errors.Error(errors.ErrInvalidHeight)
	}
	if err := m.Block.SanityCheck(); err != nil {
		return err
	}
	if err := m.Certificate.SanityCheck(); err != nil {
		return err
	}

	return nil
}

func (m *BlockAnnounceMessage) Type() Type {
	return MessageTypeBlockAnnounce
}

func (m *BlockAnnounceMessage) Fingerprint() string {
	return fmt.Sprintf("{⌘ %d %v}", m.Height, m.Block.Hash().Fingerprint())
}
