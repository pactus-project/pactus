package types

import (
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
)

// AddressInfo represents the information about a wallet address.
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

// WalletInfo represents the information about the wallet.
type WalletInfo struct {
	Version    int
	Path       string
	Network    genesis.ChainType
	DefaultFee amount.Amount
	UUID       string
	Encrypted  bool
	CreatedAt  time.Time
}
