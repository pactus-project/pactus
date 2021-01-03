package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

const DownloadResponseCodeOK = 0
const DownloadResponseCodeNoMoreBlock = 1
const DownloadResponseCodeRejected = 2
const DownloadResponseCodeBusy = 3

type DownloadResponsePayload struct {
	Status       int            `cbor:"1,keyasint"`
	From         int            `cbor:"2,keyasint"`
	Blocks       []*block.Block `cbor:"3,keyasint, omitempty"`
	Transactions []*tx.Tx       `cbor:"4,keyasint, omitempty"`
}

func (p *DownloadResponsePayload) SanityCheck() error {
	if p.From <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
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

func (p *DownloadResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{%v %v %v}", p.Status, p.From, len(p.Blocks))
}
