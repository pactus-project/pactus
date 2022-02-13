package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

const LatestBlocksResponseCodeOK = 0
const LatestBlocksResponseCodeNoMoreBlock = 1

type BlocksResponsePayload struct {
	ResponseCode    ResponseCode       `cbor:"1,keyasint"`
	SessionID       int                `cbor:"2,keyasint"`
	From            int                `cbor:"3,keyasint"`
	Blocks          []*block.Block     `cbor:"4,keyasint"`
	Transactions    []*tx.Tx           `cbor:"5,keyasint"`
	LastCertificate *block.Certificate `cbor:"6,keyasint"`
}

func NewBlocksResponsePayload(code ResponseCode, sid int, from int,
	blocks []*block.Block, trxs []*tx.Tx, cert *block.Certificate) *BlocksResponsePayload {
	return &BlocksResponsePayload{
		ResponseCode:    code,
		SessionID:       sid,
		From:            from,
		Blocks:          blocks,
		Transactions:    trxs,
		LastCertificate: cert,
	}
}
func (p *BlocksResponsePayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	for _, b := range p.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "invalid block: %v", err)
		}
	}
	if p.LastCertificate != nil {
		if err := p.LastCertificate.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "invalid certificate: %v", err)
		}
	}
	for _, trx := range p.Transactions {
		if err := trx.SanityCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (p *BlocksResponsePayload) Type() Type {
	return PayloadTypeBlocksResponse
}

func (p *BlocksResponsePayload) To() int {
	if len(p.Blocks) == 0 {
		return p.From
	}
	return p.From + len(p.Blocks) - 1
}

func (p *BlocksResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", p.SessionID, p.ResponseCode, p.From, p.To())
}

func (p *BlocksResponsePayload) IsRequestRejected() bool {
	if p.ResponseCode == ResponseCodeBusy ||
		p.ResponseCode == ResponseCodeRejected {
		return true
	}

	return false
}
