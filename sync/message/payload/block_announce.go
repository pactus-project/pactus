package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlockAnnouncePayload struct {
	Height int           `cbor:"1,keyasint"`
	Block  *block.Block  `cbor:"2,keyasint"`
	Commit *block.Commit `cbor:"3,keyasint"`
}

func (p *BlockAnnouncePayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if err := p.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
	}
	if err := p.Commit.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid commit: %v", err)
	}

	return nil
}

func (p *BlockAnnouncePayload) Type() PayloadType {
	return PayloadTypeBlockAnnounce
}

func (p *BlockAnnouncePayload) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v}", p.Block.Hash().Fingerprint())
}
