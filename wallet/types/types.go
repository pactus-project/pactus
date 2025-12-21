package types

import (
	"time"

	"github.com/pactus-project/pactus/types/amount"
)

// AddressInfo represents an address entry in the wallet.
type AddressInfo struct {
	Address   string    `json:"address"`
	PublicKey string    `json:"public_key"`
	Label     string    `json:"label"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HistoryInfo struct {
	TxID        string
	Time        *time.Time
	PayloadType string
	Desc        string
	Amount      amount.Amount
}

// WalletInfo contains wallet metadata.
type WalletInfo struct {
	Version    int
	Network    string
	DefaultFee amount.Amount
	UUID       string
	Encrypted  bool
	CreatedAt  time.Time
}
