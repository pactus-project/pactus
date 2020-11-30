package txpool

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func (pool *txPool) validateTx(trx *tx.Tx) error {
	curHeight := pool.sandbox.CurrentHeight()
	height := pool.sandbox.RecentBlockHeight(trx.Stamp())
	interval := pool.sandbox.TransactionToLiveInterval()

	if height == -1 || curHeight-height > interval {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid stamp")
	}
	if len(trx.Memo()) > pool.sandbox.MaxMemoLenght() {
		return errors.Errorf(errors.ErrInvalidTx, "Memo length exceeded")
	}
	if !trx.IsMintbaseTx() {
		fee := int64(float64(trx.Payload().Value()) * pool.sandbox.FeeFraction())
		fee = util.Max64(fee, pool.sandbox.MinFee())
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: %v, got: %v", fee, trx.Fee())
		}
	}

	// TODO:
	// validate transaction
	// Sequence number and amount

	return nil
}
