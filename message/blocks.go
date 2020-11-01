package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlocksPayload struct {
	From       int           `cbor:"1,keyasint"`
	Blocks     []block.Block `cbor:"2,keyasint"`
	LastCommit *block.Commit `cbor:"3,keyasint, omitempty"`
}

func NewBlocksMessage(from int, blocks []block.Block, lastCommit *block.Commit) Message {
	return Message{
		Type: PayloadTypeBlocks,
		Payload: &BlocksPayload{
			From:       from,
			Blocks:     blocks,
			LastCommit: lastCommit,
		},
	}

}
func (p *BlocksPayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
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
	return fmt.Sprintf("{%v-%v}", p.From, p.From+len(p.Blocks))
}
