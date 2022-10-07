package wallet

import (
	"encoding/hex"
	"sort"
	"time"

	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
)

type HistoryInfo struct {
	TxID        string
	Time        time.Time
	PayloadType string
	Desc        string
	Amount      int64
}

type transaction struct {
	BlockHash   string `json:"block"`
	BlockTime   uint32 `json:"time"`
	PayloadType string `json:"type"`
	Data        string `json:"data"`
}

type activity struct {
	TxID   string `json:"id"`
	Desc   string `json:"desc"`
	Amount int64  `json:"amount"`
}

type history struct {
	Activities   map[string][]activity  `json:"activities"`
	Transactions map[string]transaction `json:"transactions"`
}

func (h *history) hasTransaction(id string) bool {
	_, ok := h.Transactions[id]
	return ok
}

func (h *history) addActivity(addr string, amount int64, trx *pactus.TransactionResponse) {
	if h.Activities == nil {
		h.Activities = map[string][]activity{}
		h.Transactions = map[string]transaction{}
	}
	if len(h.Activities[addr]) == 0 {
		h.Activities[addr] = make([]activity, 0, 1)
	}
	act := activity{
		TxID:   hex.EncodeToString(trx.Transaction.Id),
		Amount: amount,
	}
	h.Activities[addr] = append(h.Activities[addr], act)
	sort.Slice(h.Activities[addr], func(i, j int) bool {
		return h.Transactions[h.Activities[addr][i].TxID].BlockTime <
			h.Transactions[h.Activities[addr][j].TxID].BlockTime
	})

	h.Transactions[act.TxID] = transaction{
		BlockHash:   hex.EncodeToString(trx.BlockHash),
		BlockTime:   trx.BlockTime,
		PayloadType: payload.Type(trx.Transaction.Type).String(),
		Data:        hex.EncodeToString(trx.Transaction.Data),
	}
}

func (h *history) getAddrHistory(addr string) []HistoryInfo {
	addrActs := h.Activities[addr]
	history := make([]HistoryInfo, len(addrActs))
	for i, act := range addrActs {
		t := h.Transactions[act.TxID]
		history[i].Amount = act.Amount
		history[i].TxID = act.TxID
		history[i].Time = time.Unix(int64(t.BlockTime), 0)
		history[i].PayloadType = t.PayloadType
	}

	return history
}
