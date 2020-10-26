package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlockPayload struct {
	Height int          `cbor:"1,keyasint"`
	Block  block.Block  `cbor:"2,keyasint"`
	Commit block.Commit `cbor:"3,keyasint"`
}

func NewBlockMessage(height int, block block.Block, commit block.Commit) Message {
	return Message{
		Type: PayloadTypeBlock,
		Payload: &BlockPayload{
			Height: height,
			Block:  block,
			Commit: commit,
		},
	}

}
func (p *BlockPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
	}
	if err := p.Commit.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
	}
	return nil
}

func (p *BlockPayload) Type() PayloadType {
	return PayloadTypeBlock
}

func (p *BlockPayload) Fingerprint() string {
	return fmt.Sprintf("{🗃 %v}", p.Block.Fingerprint())
}
