package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type TxsPayload struct {
	Txs []*tx.Tx `cbor:"2,keyasint"`
}

func (p *TxsPayload) SanityCheck() error {
	if len(p.Txs) == 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "No transaction")
	}
	for _, tx := range p.Txs {
		if err := tx.SanityCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (p *TxsPayload) Type() PayloadType {
	return PayloadTypeTxs
}

func (p *TxsPayload) Fingerprint() string {
	var s string
	for _, tx := range p.Txs {
		s += fmt.Sprintf("%v ", tx.ID().Fingerprint())
	}
	return fmt.Sprintf("{%v: ⌘ [%v]}", len(p.Txs), s)
}
