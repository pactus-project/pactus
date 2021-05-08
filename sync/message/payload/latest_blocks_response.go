package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

const LatestBlocksResponseCodeOK = 0
const LatestBlocksResponseCodeNoMoreBlock = 1

type LatestBlocksResponsePayload struct {
	ResponseCode    ResponseCode       `cbor:"1,keyasint"`
	SessionID       int                `cbor:"2,keyasint"`
	Target          peer.ID            `cbor:"3,keyasint"`
	From            int                `cbor:"4,keyasint"`
	Blocks          []*block.Block     `cbor:"5,keyasint"`
	Transactions    []*tx.Tx           `cbor:"6,keyasint"`
	LastCertificate *block.Certificate `cbor:"7,keyasint"`
}

func NewLatestBlocksResponsePayload(code ResponseCode, sid int, target peer.ID, from int,
	blocks []*block.Block, trxs []*tx.Tx, cert *block.Certificate) Payload {
	return &LatestBlocksResponsePayload{
		ResponseCode:    code,
		SessionID:       sid,
		Target:          target,
		From:            from,
		Blocks:          blocks,
		Transactions:    trxs,
		LastCertificate: cert,
	}
}
func (p *LatestBlocksResponsePayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid target peer id: %v", err)
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

func (p *LatestBlocksResponsePayload) Type() PayloadType {
	return PayloadTypeLatestBlocksResponse
}

func (p *LatestBlocksResponsePayload) To() int {
	if len(p.Blocks) == 0 {
		return p.From
	}
	return p.From + len(p.Blocks) - 1
}

func (p *LatestBlocksResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", p.SessionID, p.ResponseCode, p.From, p.To())
}

func (p *LatestBlocksResponsePayload) IsRequestNotProcessed() bool {
	if p.ResponseCode == ResponseCodeBusy ||
		p.ResponseCode == ResponseCodeRejected {
		return true
	}

	return false
}
