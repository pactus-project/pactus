package wallet2

import (
	"encoding/hex"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/wallet2/db"
)

type history struct {
	db db.DB
}

func newHistory(database db.DB) *history {
	return &history{
		db: database,
	}
}

func (h *history) hasTransaction(id string) bool {
	t, err := h.db.GetTransactionByTxID(id)
	if err != nil || t == nil {
		return false
	}

	return true
}

func (h *history) addTransaction(
	addr string,
	amount int64,
	txID hash.Hash,
	blockHeight uint32,
	blockTime uint32,
	payloadType string,
	data []byte,
	status db.Status,
) error {
	var desc string
	if status == db.Confirmed {
		desc = "Confirmed..."
	} else if status == db.Pending {
		desc = "Pending..."
	}

	t := &db.Transaction{
		TxID:        txID.String(),
		Address:     addr,
		BlockHeight: blockHeight,
		BlockTime:   blockTime,
		PayloadType: payloadType,
		Data:        hex.EncodeToString(data),
		Description: desc,
		Amount:      amount,
		Status:      int(status),
	}

	_, err := h.db.InsertIntoTransaction(t)

	return err
}

func (h *history) getAddrHistory(addr string) ([]db.Transaction, error) {
	return h.db.GetAllTransactions(db.WithTransactionAddr(addr))
}
