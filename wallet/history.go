package wallet

import (
	"encoding/hex"
	"sort"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type HistoryInfo struct {
	TxID        string
	Time        *time.Time
	PayloadType string
	Desc        string
	Amount      amount.Amount
}

type transaction struct {
	BlockHeight uint32 `json:"height"`
	BlockTime   uint32 `json:"time"`
	PayloadType string `json:"type"`
	Data        string `json:"data"`
}

type activity struct {
	TxID   string        `json:"id"`
	Desc   string        `json:"desc"`
	Amount amount.Amount `json:"amount"`
}

type pending struct {
	TxID   string        `json:"id"`
	Amount amount.Amount `json:"amount"`
	Data   string        `json:"data"`
}

type history struct {
	Transactions map[string]transaction `json:"transactions"`
	Activities   map[string][]activity  `json:"activities"`
	Pendings     map[string][]pending   `json:"pendings"`
}

func (h *history) hasTransaction(id string) bool {
	_, ok := h.Transactions[id]

	return ok
}

func (h *history) addActivity(addr string, amt amount.Amount, trx *pactus.GetTransactionResponse) {
	if h.Activities == nil {
		h.Activities = map[string][]activity{}
		h.Transactions = map[string]transaction{}
	}
	if len(h.Activities[addr]) == 0 {
		h.Activities[addr] = make([]activity, 0, 1)
	}
	act := activity{
		TxID:   hex.EncodeToString(trx.Transaction.Id),
		Amount: amt,
	}
	h.Activities[addr] = append(h.Activities[addr], act)
	sort.Slice(h.Activities[addr], func(i, j int) bool {
		return h.Transactions[h.Activities[addr][i].TxID].BlockTime <
			h.Transactions[h.Activities[addr][j].TxID].BlockTime
	})

	h.Transactions[act.TxID] = transaction{
		BlockHeight: trx.BlockHeight,
		BlockTime:   trx.BlockTime,
		PayloadType: payload.Type(trx.Transaction.PayloadType).String(),
		Data:        hex.EncodeToString(trx.Transaction.Data),
	}
}

func (h *history) addPending(addr string, amt amount.Amount, txID hash.Hash, data []byte) {
	if h.Pendings == nil {
		h.Pendings = map[string][]pending{}
	}
	if len(h.Pendings[addr]) == 0 {
		h.Pendings[addr] = make([]pending, 0, 1)
	}
	pnd := pending{
		TxID:   txID.String(),
		Amount: amt,
		Data:   hex.EncodeToString(data),
	}
	h.Pendings[addr] = append(h.Pendings[addr], pnd)
}

func (h *history) getAddrHistory(addr string) []HistoryInfo {
	addrActs := h.Activities[addr]
	addrPnds := h.Pendings[addr]
	history := make([]HistoryInfo, 0, len(addrActs)+len(addrPnds))
	for _, pnd := range addrPnds {
		history = append(history, HistoryInfo{
			Amount: pnd.Amount,
			TxID:   pnd.TxID,
			Desc:   "Pending...",
			Time:   nil,
		})
	}

	for _, act := range addrActs {
		trx := h.Transactions[act.TxID]
		tme := time.Unix(int64(trx.BlockTime), 0)
		history = append(history, HistoryInfo{
			Amount:      act.Amount,
			TxID:        act.TxID,
			Desc:        act.Desc,
			PayloadType: trx.PayloadType,
			Time:        &tme,
		})
	}

	return history
}
