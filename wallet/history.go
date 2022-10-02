package wallet

import "sort"

type Transaction struct {
	BlockHash string `json:"block"`
	BlockTime uint32 `json:"time"`
	Data      string `json:"data"`
}

type Activity struct {
	TxID        string `json:"id"`
	Status      string `json:"status"`
	BlockTime   uint32 `json:"time"`
	PayloadType string `json:"type"`
	Amount      int64  `json:"amount"`
}

type history struct {
	Activities   map[string][]Activity  `json:"activities"`
	Transactions map[string]Transaction `json:"transactions"`
}

func (h *history) hasTransaction(id string) bool {
	_, ok := h.Transactions[id]
	return ok
}

func (h *history) addTransaction(addr string, id string,
	activity Activity, transaction Transaction) {
	if h.Activities == nil {
		h.Activities = map[string][]Activity{}
		h.Transactions = map[string]Transaction{}
	}
	if len(h.Activities[addr]) == 0 {
		h.Activities[addr] = make([]Activity, 0, 1)
	}
	h.Activities[addr] = append(h.Activities[addr], activity)
	sort.Slice(h.Activities[addr], func(i, j int) bool {
		return h.Activities[addr][i].BlockTime < h.Activities[addr][j].BlockTime
	})

	h.Transactions[id] = transaction
}
