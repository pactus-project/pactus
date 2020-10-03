package message

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlocksPayload struct {
	From   int           `cbor:"1,keyasint"`
	Blocks []block.Block `cbor:"2,keyasint"`
	// TODO: do we need it
	//LastCommit block.Commit  `cbor:"3,keyasint"`
}

func NewBlocksMessage(height int, blocks []block.Block) *Message {
	return &Message{
		Type: PayloadTypeBlocks,
		Payload: &BlocksPayload{
			Blocks: blocks,
		},
	}

}
func (p *BlocksPayload) SanityCheck() error {
	for _, b := range p.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
		}
	}
	return nil
}

func (p *BlocksPayload) Type() PayloadType {
	return PayloadTypeBlocks
}

func (p *BlocksPayload) Fingerprint() string {
	return ""
}
