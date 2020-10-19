package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type BlocksResPayload struct {
	From   int           `cbor:"1,keyasint"`
	Blocks []block.Block `cbor:"2,keyasint"`
	Txs    []tx.Tx       `cbor:"3,keyasint"`

	// TODO: do we need it
	//LastCommit block.Commit  `cbor:"3,keyasint"`
}

func NewBlocksMessage(height int, blocks []block.Block, txs []tx.Tx) Message {
	return Message{
		Type: PayloadTypeBlocksRes,
		Payload: &BlocksResPayload{
			Blocks: blocks,
			Txs:    txs,
		},
	}

}
func (p *BlocksResPayload) SanityCheck() error {
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

func (p *BlocksResPayload) Type() PayloadType {
	return PayloadTypeBlocksRes
}

func (p *BlocksResPayload) Fingerprint() string {
	return fmt.Sprintf("{%v-%v}", p.From, p.From+len(p.Blocks))
}
