package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlocksPayload struct {
	From       int            `cbor:"1,keyasint"`
	Blocks     []*block.Block `cbor:"2,keyasint"`
	LastCommit *block.Commit  `cbor:"3,keyasint, omitempty"`
}

func (p *BlocksPayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid Height")
	}
	if len(p.Blocks) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "No block")
	}
	for _, b := range p.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
		}
	}
	if p.LastCommit != nil {
		if err := p.LastCommit.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid commit: %v", err)
		}
	}
	return nil
}

func (p *BlocksPayload) Type() PayloadType {
	return PayloadTypeBlocks
}

func (p *BlocksPayload) To() int {
	return p.From + len(p.Blocks) - 1
}

func (p *BlocksPayload) Fingerprint() string {
	var s string
	for _, b := range p.Blocks {
		s += fmt.Sprintf("%v ", b.Hash().Fingerprint())
	}
	return fmt.Sprintf("{%v-%v: ⌘ [%v]}", p.From, p.From+len(p.Blocks)-1, s)
}
