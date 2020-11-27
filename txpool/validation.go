package txpool

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func (pool *txPool) validateTx(trx *tx.Tx) error {
	if !pool.stamps.Has(trx.Stamp()) {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid stamp")
	}
	if len(trx.Memo()) > pool.maxMemoLenght {
		return errors.Errorf(errors.ErrInvalidTx, "Memo length exceeded")
	}
	if !trx.IsMintbaseTx() {
		fee := int64(float64(trx.Payload().Value()) * pool.feeFraction)
		fee = util.Max64(fee, pool.minFee)
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: %v, got: %v", fee, trx.Fee())
		}
	}

	// TODO:
	// validate transaction
	// Sequence number and amount

	return nil
}
