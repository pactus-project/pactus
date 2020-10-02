package message

import (
	"gitlab.com/zarb-chain/zarb-go/block"
	"gitlab.com/zarb-chain/zarb-go/errors"
)

type BlockPayload struct {
	Block block.Block `cbor:"1,keyasint"`
}

func NewBlockMessage(height int, block block.Block) *Message {
	return &Message{
		Type:   PayloadTypeBlock,
		Height: height,
		Payload: &BlockPayload{
			Block: block,
		},
	}

}
func (p *BlockPayload) SanityCheck() error {
	if err := p.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
	}
	return nil
}

func (p *BlockPayload) Type() PayloadType {
	return PayloadTypeBlock
}

func (p *BlockPayload) Fingerprint() string {
	return ""
}
