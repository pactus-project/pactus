package tx

import "github.com/zarbchain/zarb-go/errors"

type CommittedTx struct {
	Tx      *Tx      `cbor:"1,keyasint"`
	Receipt *Receipt `cbor:"2,keyasint"`
}

func (ctrx *CommittedTx) SanityCheck() error {
	if err := ctrx.Tx.SanityCheck(); err != nil {
		return err
	}
	if err := ctrx.Receipt.SanityCheck(); err != nil {
		return err
	}
	if !ctrx.Receipt.data.TxID.EqualsTo(ctrx.Tx.ID()) {
		return errors.Errorf(errors.ErrInvalidReceipt, "Mismatched transaction hash")
	}
	return nil
}
