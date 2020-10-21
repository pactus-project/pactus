package message

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type BlockPayload struct {
	Height int         `cbor:"1,keyasint"`
	Block  block.Block `cbor:"2,keyasint"`
	Txs    []tx.Tx     `cbor:"3,keyasint"`
}

func NewBlockMessage(height int, block block.Block, txs []tx.Tx) Message {
	return Message{
		Type: PayloadTypeBlock,
		Payload: &BlockPayload{
			Height: height,
			Block:  block,
			Txs:    txs,
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
	// TODO: compare tx.hash() with tx_hashes insied block
	if p.Block.TxHashes().Count() != len(p.Txs) {
		return errors.Errorf(errors.ErrInvalidMessage, "Not enough transactions")
	}
	return nil
}

func (p *BlockPayload) Type() PayloadType {
	return PayloadTypeBlock
}

func (p *BlockPayload) Fingerprint() string {
	return ""
}
