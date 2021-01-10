package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type LatestBlocksResponsePayload struct {
	RequestID    int            `cbor:"1,keyasint"`
	From         int            `cbor:"2,keyasint"`
	Blocks       []*block.Block `cbor:"3,keyasint"`
	Transactions []*tx.Tx       `cbor:"4,keyasint, omitempty"`
	LastCommit   *block.Commit  `cbor:"5,keyasint, omitempty"`
}

func (p *LatestBlocksResponsePayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid Height")
	}
	if len(p.Blocks) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "No block")
	}
	if len(p.Transactions) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "No transaction")
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

	for _, trx := range p.Transactions {
		if err := trx.SanityCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (p *LatestBlocksResponsePayload) Type() PayloadType {
	return PayloadTypeLatestBlocksResponse
}

func (p *LatestBlocksResponsePayload) To() int {
	return p.From + len(p.Blocks) - 1
}

func (p *LatestBlocksResponsePayload) Fingerprint() string {
	var s string
	for _, b := range p.Blocks {
		s += fmt.Sprintf("%v ", b.Hash().Fingerprint())
	}
	return fmt.Sprintf("{%v-%v: ⌘ [%v]}", p.From, p.From+len(p.Blocks)-1, s)
}
