package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type DownloadResponsePayload struct {
	ResponseCode ResponseCode   `cbor:"1,keyasint"`
	SessionID    int            `cbor:"2,keyasint"`
	Target       peer.ID        `cbor:"3,keyasint"`
	From         int            `cbor:"4,keyasint"`
	Blocks       []*block.Block `cbor:"5,keyasint"`
	Transactions []*tx.Tx       `cbor:"6,keyasint"`
}

func NewDownloadResponsePayload(code ResponseCode, sid int, target peer.ID, from int,
	blocks []*block.Block, trxs []*tx.Tx) Payload {
	return &DownloadResponsePayload{
		ResponseCode: code,
		SessionID:    sid,
		Target:       target,
		From:         from,
		Blocks:       blocks,
		Transactions: trxs,
	}
}

func (p *DownloadResponsePayload) SanityCheck() error {
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
	for _, trx := range p.Transactions {
		if err := trx.SanityCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (p *DownloadResponsePayload) Type() PayloadType {
	return PayloadTypeDownloadResponse
}

func (p *DownloadResponsePayload) To() int {
	if len(p.Blocks) == 0 {
		return p.From
	}
	return p.From + len(p.Blocks) - 1
}

func (p *DownloadResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", p.SessionID, p.ResponseCode, p.From, p.To())
}
