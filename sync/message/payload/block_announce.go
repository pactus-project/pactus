package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlockAnnouncePayload struct {
	Height      int                `cbor:"1,keyasint"`
	Block       *block.Block       `cbor:"2,keyasint"`
	Certificate *block.Certificate `cbor:"3,keyasint"`
}

func NewBlockAnnouncePayload(h int, b *block.Block, c *block.Certificate) *BlockAnnouncePayload {
	return &BlockAnnouncePayload{
		Height:      h,
		Block:       b,
		Certificate: c,
	}
}

func (p *BlockAnnouncePayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if err := p.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid block: %v", err)
	}
	if err := p.Certificate.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid certificate: %v", err)
	}

	return nil
}

func (p *BlockAnnouncePayload) Type() Type {
	return PayloadTypeBlockAnnounce
}

func (p *BlockAnnouncePayload) Fingerprint() string {
	return fmt.Sprintf("{⌘ %d %v}", p.Height, p.Block.Hash().Fingerprint())
}
