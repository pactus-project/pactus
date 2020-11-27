package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

type TxPoolReader interface {
	PendingTx(hash crypto.Hash) *tx.Tx
	HasTx(hash crypto.Hash) bool
	Size() int

	Fingerprint() string
}

type TxPool interface {
	TxPoolReader

	UpdateStampsCount(stampsCount int)
	UpdateMaxMemoLenght(maxMemoLenght int)
	UpdateFeeFraction(feeFraction float64)
	UpdateMinFee(minFee int64)
	AppendStamp(height int, stamp crypto.Hash)
	AppendTxs(txs []tx.Tx)
	AppendTx(tx tx.Tx) error
	AppendTxAndBroadcast(trx tx.Tx) error
	RemoveTx(hash crypto.Hash) *tx.Tx
}
