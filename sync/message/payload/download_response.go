package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type DownloadResponsePayload struct {
	ResponseCode ResponseCode   `cbor:"1,keyasint"`
	SessionID    int            `cbor:"2,keyasint"`
	Initiator    peer.ID        `cbor:"3,keyasint"`
	Target       peer.ID        `cbor:"4,keyasint"`
	From         int            `cbor:"5,keyasint"`
	Blocks       []*block.Block `cbor:"6,keyasint"`
	Transactions []*tx.Tx       `cbor:"7,keyasint"`
}

func (p *DownloadResponsePayload) SanityCheck() error {
	for _, b := range p.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
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
	return p.From + len(p.Blocks) - 1
}

func (p *DownloadResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{%v %v-%v}", p.ResponseCode, p.From, p.To())
}
