package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type LatestBlocksPayload struct {
	From         int            `cbor:"1,keyasint"`
	Blocks       []*block.Block `cbor:"2,keyasint"`
	Transactions []*tx.Tx       `cbor:"3,keyasint, omitempty"`
	Commit       *block.Commit  `cbor:"4,keyasint, omitempty"`
}

func (p *LatestBlocksPayload) SanityCheck() error {
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
	if p.Commit != nil {
		if err := p.Commit.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid commit: %v", err)
		}
	}

	for _, trx := range p.Transactions {
		if err := trx.SanityCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (p *LatestBlocksPayload) Type() PayloadType {
	return PayloadTypeLatestBlocks
}

func (p *LatestBlocksPayload) To() int {
	return p.From + len(p.Blocks) - 1
}

func (p *LatestBlocksPayload) Fingerprint() string {
	var s string
	for _, b := range p.Blocks {
		s += fmt.Sprintf("%v ", b.Hash().Fingerprint())
	}
	return fmt.Sprintf("{%v-%v: ⌘ [%v]}", p.From, p.From+len(p.Blocks)-1, s)
}
